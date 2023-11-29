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
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/ignition"
	"nestos-kubernetes-deployer/pkg/utils"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/sirupsen/logrus"
)

type Master struct {
	ClusterAsset   *asset.ClusterAsset
	StorageContent []utils.StorageContent
}

func (m *Master) GenerateFiles() error {
	mtd := ignition.GetTmplData(m.ClusterAsset)
	generateFile := ignition.Common{
		UserName:        m.ClusterAsset.Master.NodeAsset[0].UserName,
		SSHKey:          m.ClusterAsset.Master.NodeAsset[0].SSHKey,
		PassWord:        m.ClusterAsset.Master.NodeAsset[0].Password,
		NodeType:        "controlplane",
		TmplData:        mtd,
		EnabledServices: ignition.EnabledServices,
		Config:          &igntypes.Config{},
	}
	if err := generateFile.Generate(); err != nil {
		logrus.Errorf("failed to generate %s ignition file: %v", m.ClusterAsset.Master.NodeAsset[0].UserName, err)
		return err
	}
	for _, file := range m.StorageContent {
		ignFile := ignition.FileWithContents(file.Path, file.Mode, file.Content)
		generateFile.Config.Storage.Files = ignition.AppendFiles(generateFile.Config.Storage.Files, ignFile)
	}
	data, err := ignition.Marshal(generateFile.Config)
	if err != nil {
		logrus.Errorf("failed to Marshal ignition config: %v", err)
		return err
	}
	m.ClusterAsset.Master.NodeAsset[0].Ign_Data = data
	for i := 1; i < m.ClusterAsset.Master.Count; i++ {
		generateFile.UserName = m.ClusterAsset.Master.NodeAsset[i].UserName
		generateFile.SSHKey = m.ClusterAsset.Master.NodeAsset[i].SSHKey
		generateFile.PassWord = m.ClusterAsset.Master.NodeAsset[i].Password
		generateFile.NodeType = "master"
		if err := generateFile.Generate(); err != nil {
			logrus.Errorf("failed to generate %s ignition file: %v", m.ClusterAsset.Master.NodeAsset[i].UserName, err)
			return err
		}
		data, err := ignition.Marshal(generateFile.Config)
		if err != nil {
			logrus.Errorf("failed to Marshal ignition config: %v", err)
			return err
		}
		m.ClusterAsset.Master.NodeAsset[i].Ign_Data = data
	}

	return nil
}
