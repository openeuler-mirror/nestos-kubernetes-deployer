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

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	ClusterAsset *asset.ClusterAsset
}

func (w *Worker) GenerateFiles() error {
	wtd := ignition.GetTmplData(w.ClusterAsset)

	for _, worker := range w.ClusterAsset.Worker {
		config := &igntypes.Config{}
		wtd.NodeName = worker.Hostname
		generateFile := ignition.Common{
			NodeType:        "worker",
			TmplData:        wtd,
			EnabledServices: ignition.EnabledServices,
			Config:          config,
			UserName:        worker.UserName,
			SSHKey:          worker.SSHKey,
			PassWord:        worker.Password,
		}

		// Generate Ignition data
		if err := generateFile.Generate(); err != nil {
			logrus.Errorf("failed to generate %s ignition file: %v", worker.UserName, err)
			return err
		}
		data, err := ignition.Marshal(generateFile.Config)
		if err != nil {
			logrus.Errorf("failed to Marshal ignition config: %v", err)
			return err
		}

		// Assign the Ignition data to the Worker node
		worker.Ign_Data = string(data)
	}

	return nil
}
