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
package machine

import (
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/ignition"
	"nestos-kubernetes-deployer/pkg/utils"
	"os"
	"path/filepath"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/sirupsen/logrus"
)

const (
	MasterIgnFilename            = "master.ign"
	ControlplaneIgnFilename      = "controlplane.ign"
	masterMergeIgnFilename       = "master-merge.ign"
	controlplaneMergeIgnFilename = "controlplane-merge.ign"
)

type Master struct {
	ClusterAsset     *asset.ClusterAsset
	BootstrapBaseurl string
}

func (m *Master) GenerateFiles() error {
	sshkeyContent, err := os.ReadFile(m.ClusterAsset.SSHKey)
	if err != nil {
		logrus.Debug("Failed to read sshkey content:", err)
		return err
	}

	// Get template dependency configuration
	masterTemplateData, err := ignition.GetTmplData(m.ClusterAsset)
	if err != nil {
		return err
	}
	ignitionDir := filepath.Join(configmanager.GetPersistDir(), m.ClusterAsset.Cluster_ID, "ignition")

	for i, master := range m.ClusterAsset.Master {
		nodeType := getNodeTypeName(i)
		masterTemplateData.NodeName = master.Hostname

		generateFile := ignition.Common{
			UserName:        m.ClusterAsset.UserName,
			SSHKey:          string(sshkeyContent),
			PassWord:        m.ClusterAsset.Password,
			NodeType:        nodeType,
			TmplData:        masterTemplateData,
			EnabledServices: ignition.EnabledServices,
			Config:          &igntypes.Config{},
		}

		// Generate Ignition data
		if err := generateFile.Generate(); err != nil {
			logrus.Errorf("failed to generate %s ignition file: %v", master.Hostname, err)
			return err
		}

		filename := MasterIgnFilename
		mergeFilename := masterMergeIgnFilename
		if i == 0 {
			filename = ControlplaneIgnFilename
			mergeFilename = controlplaneMergeIgnFilename
			mergeCertificatesIntoConfig(generateFile.Config, master.Certs)
		}

		if len(m.ClusterAsset.ShellFiles) > 0 {
			ignition.MergeHookFilesIntoConfig(generateFile.Config, m.ClusterAsset.ShellFiles)
		}

		m.ClusterAsset.Master[i].Ignitions.CreateIgnPath = filepath.Join(ignitionDir, filename)
		m.ClusterAsset.Master[i].Ignitions.MergeIgnPath = filepath.Join(ignitionDir, mergeFilename)

		if err := ignition.SaveFile(generateFile.Config, ignitionDir, filename); err != nil {
			return err
		}

		mergerConfig := ignition.GenerateMergeIgnition(m.BootstrapBaseurl, filename)
		if err := ignition.SaveFile(mergerConfig, ignitionDir, mergeFilename); err != nil {
			return err
		}

		data, err := ignition.Marshal(generateFile.Config)
		if err != nil {
			logrus.WithError(err).Error("Failed to Marshal ignition config")
			return err
		}
		m.ClusterAsset.Master[i].CreateIgnContent = data
	}

	return nil
}

func getNodeTypeName(index int) string {
	if index == 0 {
		return "controlplane"
	}
	return "master"
}

// Merge certificates into ignition.Config
func mergeCertificatesIntoConfig(config *igntypes.Config, certs []utils.StorageContent) {
	for _, file := range certs {
		ignFile := ignition.FileWithContents(file.Path, file.Mode, file.Content)
		config.Storage.Files = ignition.AppendFiles(config.Storage.Files, ignFile)
	}
}
