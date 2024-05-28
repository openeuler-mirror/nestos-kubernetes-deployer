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
	"errors"
	"nestos-kubernetes-deployer/cmd/command"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/infraasset"
	"nestos-kubernetes-deployer/pkg/utils"
	"os"
	"runtime"
	"strings"

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
	arch := runtime.GOARCH
	if opts.Opts.Arch != "" {
		arch = opts.Opts.Arch
	}

	conf, err := asset.GetDefaultClusterConfig(arch, opts.Opts.Platform)
	if err != nil {
		return err
	}
	conf.InfraPlatform = getDefaultInfraAsset(strings.ToLower(conf.Platform))

	data, err := yaml.Marshal(conf)
	if err != nil {
		logrus.Errorf("Faild to marshal template config: %v", err)
		return err
	}

	output, _ := cmd.Flags().GetString("output")
	if output == "" {
		output = "."
	}
	if !strings.HasSuffix(output, "/") {
		output = output + "/"
	}

	if err := os.WriteFile(output+"template.yaml", data, utils.DeployConfigFileMode); err != nil {
		logrus.Errorf("Faild to write template config file: %v", err)
		return err
	}

	return nil
}

func getDefaultInfraAsset(platform string) interface{} {
	switch platform {
	case "libvirt":
		return infraasset.LibvirtAsset{
			URI:     "qemu:///system",
			CIDR:    "192.168.132.0/24",
			Gateway: "192.168.132.1",
		}
	case "openstack":
		return infraasset.OpenStackAsset{
			UserName:         "admin",
			TenantName:       "admin",
			AuthURL:          "http://controller:5000/v3",
			Region:           "RegionOne",
			AvailabilityZone: "nova",
		}
	case "pxe":
		return infraasset.PXEAsset{
			HTTPServerPort: "9080",
			HTTPRootDir:    "/var/www/html/",
			TFTPServerPort: "69",
			TFTPRootDir:    "/var/lib/tftpboot/",
		}
	case "ipxe":
		return infraasset.IPXEAsset{
			Port:              "9080",
			OSInstallTreePath: "/var/www/html/",
		}
	default:
		return errors.New("unsupported platform")
	}
}
