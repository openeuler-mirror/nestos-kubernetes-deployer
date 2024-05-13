package kickstart

import (
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/constants"
	"nestos-kubernetes-deployer/pkg/osmanager/bootconfig"
	"path/filepath"
)

const (
	kickstartControlplane = "controlplane.cfg"
	kickstartMaster       = "master.cfg"
	kickstartWorker       = "worker.cfg"
)

type Kickstart struct {
	ClusterAsset     *asset.ClusterAsset
	BootstrapBaseurl string
}

func NewKickstart(clusterAsset *asset.ClusterAsset, bootstrapBaseurl string) *Kickstart {
	return &Kickstart{
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

func (c *Kickstart) GenerateBootConfig() error {
	if err := c.generateNodeConfig(constants.Controlplane, constants.InitClusterService, constants.InitClusterYaml, kickstartControlplane); err != nil {
		return err
	}

	if len(c.ClusterAsset.Master) > 1 {
		if err := c.generateNodeConfig(constants.Master, constants.JoinMasterService, "", kickstartMaster); err != nil {
			return err
		}
	}

	if len(c.ClusterAsset.Worker) > 0 {
		if err := c.generateNodeConfig(constants.Worker, constants.JoinWorkerService, "", kickstartWorker); err != nil {
			return err
		}
	}

	return nil
}

func (c *Kickstart) generateNodeConfig(nodeType, service string, yamlPath string, filename string) error {
	tmpl := newTemplate(c.ClusterAsset, append(enabledServices, service), enabledFiles)
	if yamlPath != "" {
		tmpl.enabledFiles = append(tmpl.enabledFiles, yamlPath)
	}
	if err := tmpl.GenerateBootConfig(c.BootstrapBaseurl, nodeType); err != nil {
		return err
	}
	savePath := bootconfig.GetSavePath(c.ClusterAsset.Cluster_ID)
	if err := bootconfig.SaveFile(tmpl.config, savePath, filename); err != nil {
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
