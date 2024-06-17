package kickstart

import (
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/constants"
	"nestos-kubernetes-deployer/pkg/osmanager/bootconfig"
	"os"
	"path/filepath"
	"strings"
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
	if err := c.generateNodeConfig(constants.Controlplane, constants.InitClusterService, constants.InitClusterYaml, c.ClusterAsset.Master[0].Hostname+constants.KickstartSuffix); err != nil {
		return err
	}

	n := len(c.ClusterAsset.Master)
	if n > 1 {
		for i := 1; i < n; i++ {
			if err := c.generateNodeConfig(constants.Master, constants.JoinMasterService, "", c.ClusterAsset.Master[i].Hostname+constants.KickstartSuffix); err != nil {
				return err
			}
		}
	}

	if err := c.generateNodeConfig(constants.Worker, constants.JoinWorkerService, "", constants.Worker+constants.KickstartSuffix); err != nil {
		return err
	}

	return nil
}

func (c *Kickstart) generateNodeConfig(nodeType, service string, yamlPath string, filename string) error {
	tmpl := newTemplate(c.ClusterAsset, append(enabledServices, service), enabledFiles)
	if yamlPath != "" {
		tmpl.enabledFiles = append(tmpl.enabledFiles, yamlPath)
	}
	if err := tmpl.GenerateBootConfig(c.BootstrapBaseurl, nodeType, strings.TrimSuffix(filename, constants.KickstartSuffix)); err != nil {
		return err
	}
	savePath := bootconfig.GetSavePath(c.ClusterAsset.ClusterID)
	if err := bootconfig.SaveFile(tmpl.config, savePath, filename); err != nil {
		return err
	}

	ksPath := filepath.Join(savePath, filename)
	ksContent, err := os.ReadFile(ksPath)
	if err != nil {
		return err
	}

	switch nodeType {
	case constants.Controlplane:
		c.ClusterAsset.BootConfig.Controlplane = asset.BootFile{
			Content: ksContent,
			Path:    ksPath,
		}
	case constants.Master:
		c.ClusterAsset.BootConfig.KickstartMaster = append(c.ClusterAsset.BootConfig.KickstartMaster, asset.BootFile{
			Content: ksContent,
			Path:    ksPath,
		})
	case constants.Worker:
		c.ClusterAsset.BootConfig.Worker = asset.BootFile{
			Content: ksContent,
			Path:    ksPath,
		}
	}

	return nil
}
