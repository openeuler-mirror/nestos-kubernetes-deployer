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
package ignition

import (
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/constants"
	"nestos-kubernetes-deployer/pkg/osmanager/bootconfig"
	"path/filepath"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/sirupsen/logrus"
	"github.com/vincent-petithory/dataurl"
)

var (
	enabledFiles = []string{
		constants.ReleaseImagePivotFile,
		constants.SetKernelParaConf,
		constants.KubeletServiceConf,
		constants.Hosts,
	}

	enabledServices = []string{
		constants.KubeletService,
		constants.ReleaseImagePivotService,
		constants.SetKernelPara,
	}
)

type Ignition struct {
	ClusterAsset     *asset.ClusterAsset
	BootstrapBaseurl string
}

func NewIgnition(clusterAsset *asset.ClusterAsset, bootstrapBaseurl string) *Ignition {
	return &Ignition{
		ClusterAsset:     clusterAsset,
		BootstrapBaseurl: bootstrapBaseurl,
	}
}

func (ig *Ignition) GenerateBootConfig() error {
	if err := ig.generateNodeIgnition(constants.Controlplane, constants.InitClusterService, constants.InitClusterYaml, constants.ControlplaneIgn, constants.ControlplaneMergeIgn); err != nil {
		return err
	}

	if len(ig.ClusterAsset.Master) > 1 {
		if err := ig.generateNodeIgnition(constants.Master, constants.JoinMasterService, "", constants.MasterIgn, constants.MasterMergeIgn); err != nil {
			return err
		}
	}

	if len(ig.ClusterAsset.Worker) > 0 {
		if err := ig.generateNodeIgnition(constants.Worker, constants.JoinWorkerService, "", constants.WorkerIgn, constants.WorkerMergeIgn); err != nil {
			return err
		}
	}

	return nil
}

func (ig *Ignition) generateNodeIgnition(nodeType, service string, yamlPath string, ignFilename string, mergeIgnFilename string) error {
	tmpl := newTemplate(ig.ClusterAsset, &igntypes.Config{}, append(enabledServices, service), enabledFiles)
	if yamlPath != "" {
		tmpl.enabledFiles = append(tmpl.enabledFiles, yamlPath)
	}

	if err := tmpl.GenerateBootConfig(); err != nil {
		return err
	}

	// Merge certificates directly into the configuration
	if nodeType == constants.Controlplane {
		for _, cert := range ig.ClusterAsset.Master[0].Certs {
			ignFile := fileWithContents(cert.Path, cert.Mode, cert.Content)
			tmpl.config.Storage.Files = appendFiles(tmpl.config.Storage.Files, ignFile)
		}
	}

	savePath := bootconfig.GetSavePath(ig.ClusterAsset.Cluster_ID)
	if err := bootconfig.SaveJSON(tmpl.config, savePath, ignFilename); err != nil {
		return err
	}
	mergeIgnFile := generateMergeIgnition(ig.BootstrapBaseurl, ignFilename)
	if err := bootconfig.SaveJSON(mergeIgnFile, savePath, mergeIgnFilename); err != nil {
		return err
	}

	ignData, err := bootconfig.Marshal(tmpl.config)
	if err != nil {
		logrus.WithError(err).Errorf("failed to Marshal ignition config for %s node", nodeType)
		return err
	}

	switch nodeType {
	case constants.Controlplane:
		ig.ClusterAsset.BootConfig.Controlplane = asset.BootFile{
			Content:   ignData,
			Path:      filepath.Join(savePath, ignFilename),
			MergePath: filepath.Join(savePath, mergeIgnFilename),
		}
	case constants.Master:
		ig.ClusterAsset.BootConfig.Master = asset.BootFile{
			Content:   ignData,
			Path:      filepath.Join(savePath, ignFilename),
			MergePath: filepath.Join(savePath, mergeIgnFilename),
		}
	case constants.Worker:
		ig.ClusterAsset.BootConfig.Worker = asset.BootFile{
			Content:   ignData,
			Path:      filepath.Join(savePath, ignFilename),
			MergePath: filepath.Join(savePath, mergeIgnFilename),
		}
	}

	return nil
}

func fileWithContents(path string, mode int, contents []byte) igntypes.File {
	return igntypes.File{
		Node: igntypes.Node{
			Path:      path,
			Overwrite: ignutil.BoolToPtr(true),
		},
		FileEmbedded1: igntypes.FileEmbedded1{
			Mode: &mode,
			Contents: igntypes.Resource{
				Source: ignutil.StrToPtr(dataurl.EncodeBytes(contents)),
			},
		},
	}
}

func appendFiles(files []igntypes.File, file igntypes.File) []igntypes.File {
	for i, f := range files {
		if f.Node.Path == file.Node.Path {
			files[i] = file
			return files
		}
	}
	files = append(files, file)
	return files
}
