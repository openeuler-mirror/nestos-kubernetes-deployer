/*
Copyright 2023 KylinSoft  Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"io/ioutil"
	"nestos-kubernetes-deployer/cmd/command"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func NewTemplateCommand() *cobra.Command {
	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Create a default template of nkd config",
		RunE:  createTemplate,
	}

	command.SetupTemplateCmdOpts(templateCmd)

	return templateCmd
}

func createTemplate(cmd *cobra.Command, args []string) error {
	return createDeployConfigTemplate(opts.Opts.ClusterConfigFile, opts.Opts.Platform)
}

func createDeployConfigTemplate(file string, platform string) error {
	conf := &asset.ClusterAsset{
		Cluster_ID: opts.Opts.ClusterID,
		Kubernetes: asset.Kubernetes{
			Kubernetes_Version: opts.Opts.KubeVersion,
			ApiServer_Endpoint: fmt.Sprintf("%s:%s", opts.Opts.Master.IP[0], "6443"),
			Image_Registry:     opts.Opts.ImageRegistry,
			Pause_Image:        opts.Opts.PauseImage,
			Release_Image_URL:  opts.Opts.ReleaseImageUrl,
			Token:              opts.Opts.Token,
			CertificateKey:     opts.Opts.CertificateKey,
			Network: asset.Network{
				Service_Subnet:        opts.Opts.NetWork.ServiceSubnet,
				Pod_Subnet:            opts.Opts.NetWork.PodSubnet,
				CoreDNS_Image_Version: opts.Opts.NetWork.DNS.ImageVersion,
			},
		},
		Housekeeper: asset.Housekeeper{
			DeployHousekeeper:  opts.Opts.DeployHousekeeper,
			OperatorImageUrl:   opts.Opts.OperatorImageUrl,
			ControllerImageUrl: opts.Opts.ControllerImageUrl,
		},
	}
	if opts.Opts.Master.IP != nil {
		for i, v := range opts.Opts.Master.IP {
			nodeAsset := asset.NodeAsset{
				Hostname: opts.Opts.Master.Hostname[i],
				HardwareInfo: asset.HardwareInfo{
					CPU:  opts.Opts.Master.CPU,
					RAM:  opts.Opts.Master.RAM,
					Disk: opts.Opts.Master.Disk,
				},
				UserName: opts.Opts.Master.UserName,
				Password: opts.Opts.Master.Password,
				SSHKey:   opts.Opts.Master.SSHKey,
				IP:       v,
			}
			conf.Master = append(conf.Master, nodeAsset)
		}
	}
	if opts.Opts.Worker.IP != nil {
		for i, v := range opts.Opts.Worker.IP {
			nodeAsset := asset.NodeAsset{
				Hostname: opts.Opts.Worker.Hostname[i],
				HardwareInfo: asset.HardwareInfo{
					CPU:  opts.Opts.Worker.CPU,
					RAM:  opts.Opts.Worker.RAM,
					Disk: opts.Opts.Worker.Disk,
				},
				UserName: opts.Opts.Worker.UserName,
				Password: opts.Opts.Worker.Password,
				SSHKey:   opts.Opts.Worker.SSHKey,
				IP:       v,
			}
			conf.Worker = append(conf.Worker, nodeAsset)
		}
	}

	d, err := yaml.Marshal(conf)
	if err != nil {
		logrus.Errorf("faild to marshal template config: %v", err)
		return err
	}

	if err := ioutil.WriteFile(file, d, utils.DeployConfigFileMode); err != nil {
		logrus.Errorf("faild to write template config file: %v", err)
		return err
	}

	return nil
}
