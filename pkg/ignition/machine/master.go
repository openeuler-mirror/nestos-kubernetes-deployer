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
	"os"
	"path/filepath"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/sirupsen/logrus"
)

type Master struct {
	ClusterAsset *asset.ClusterAsset
}

func (m *Master) GenerateFiles() error {
	sshkeyContent, err := os.ReadFile(m.ClusterAsset.SSHKey)
	if err != nil {
		logrus.Debug("Error to read sshkey content")
	}
	//Get template dependency configuration
	mtd := ignition.GetTmplData(m.ClusterAsset)
	for i, master := range m.ClusterAsset.Master {
		nodeType := "controlplane"
		if i > 0 {
			nodeType = "master"
		}
		mtd.NodeName = master.Hostname
		generateFile := ignition.Common{
			UserName:        m.ClusterAsset.UserName,
			SSHKey:          string(sshkeyContent),
			PassWord:        m.ClusterAsset.Password,
			NodeType:        nodeType,
			TmplData:        mtd,
			EnabledServices: ignition.EnabledServices,
			Config:          &igntypes.Config{},
		}

		// Generate Ignition data
		if err := generateFile.Generate(); err != nil {
			logrus.Errorf("failed to generate %s ignition file: %v", master.Hostname, err)
			return err
		}

		// Merge certificates into ignition.Config
		if i == 0 {
			for _, file := range master.Certs {
				ignFile := ignition.FileWithContents(file.Path, file.Mode, file.Content)
				generateFile.Config.Storage.Files = ignition.AppendFiles(generateFile.Config.Storage.Files, ignFile)
			}
		}

		//Assign the Ignition path to the Master node
		filePath := filepath.Join(configmanager.GetPersistDir(), m.ClusterAsset.Cluster_ID, "ignition")
		fileName := master.Hostname + ".ign"
		m.ClusterAsset.Master[i].Ign_Path = filepath.Join(filePath, fileName)

		ignition.SaveFile(generateFile.Config, filePath, fileName)
	}

	return nil
}
