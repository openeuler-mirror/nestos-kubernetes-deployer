/*
Copyright 2024 KylinSoft  Co., Ltd.

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

package cloudinit

import (
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/constants"
	"nestos-kubernetes-deployer/pkg/osmanager/bootconfig"
	"path/filepath"
)

const (
	cloudinitControlplane = "controlplane.cfg"
	cloudinitMaster       = "master.cfg"
	cloudinitWorker       = "worker.cfg"
)

type Cloudinit struct {
	ClusterAsset     *asset.ClusterAsset
	BootstrapBaseurl string
}

func NewCloudinit(clusterAsset *asset.ClusterAsset, bootstrapBaseurl string) *Cloudinit {
	return &Cloudinit{
		ClusterAsset:     clusterAsset,
		BootstrapBaseurl: bootstrapBaseurl,
	}
}

var (
	enabledFiles = []string{
		constants.ReleaseImagePivotFile,
		constants.SetKernelParaConf,
		constants.Hosts,
	}

	enabledServices = []string{
		constants.ReleaseImagePivotService,
		constants.SetKernelPara,
	}
)

func (c *Cloudinit) GenerateBootConfig() error {
	if err := c.generateNodeConfig(constants.Controlplane, constants.InitClusterService, constants.InitClusterYaml, cloudinitControlplane); err != nil {
		return err
	}

	if len(c.ClusterAsset.Master) > 1 {
		if err := c.generateNodeConfig(constants.Master, constants.JoinMasterService, "", cloudinitMaster); err != nil {
			return err
		}
	}

	if len(c.ClusterAsset.Worker) > 0 {
		if err := c.generateNodeConfig(constants.Worker, constants.JoinWorkerService, "", cloudinitWorker); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cloudinit) generateNodeConfig(nodeType, service string, yamlPath string, filename string) error {
	tmpl := newTemplate(c.ClusterAsset, &CloudinitConfig{}, append(enabledServices, service), enabledFiles)
	if yamlPath != "" {
		tmpl.enabledFiles = append(tmpl.enabledFiles, yamlPath)
	}
	if err := tmpl.GenerateBootConfig(c.BootstrapBaseurl, nodeType); err != nil {
		return err
	}
	savePath := bootconfig.GetSavePath(c.ClusterAsset.Cluster_ID)
	if err := bootconfig.SaveYAML(tmpl.config, savePath, filename, "#cloud-config\n"); err != nil {
		return err
	}

	switch nodeType {
	case constants.Controlplane:
		c.ClusterAsset.BootConfig.Controlplane = asset.BootFile{
			Path: filepath.Join(savePath, filename),
		}
	case constants.Master:
		c.ClusterAsset.BootConfig.Master = asset.BootFile{
			Path: filepath.Join(savePath, filename),
		}
	case constants.Worker:
		c.ClusterAsset.BootConfig.Worker = asset.BootFile{
			Path: filepath.Join(savePath, filename),
		}
	}

	return nil
}
