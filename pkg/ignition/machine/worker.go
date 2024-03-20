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

const (
	WorkerIgnFilename      = "worker.ign"
	workerMergeIgnFilename = "worker-merge.ign"
)

type Worker struct {
	ClusterAsset      *asset.ClusterAsset
	Bootstrap_baseurl string
}

func (w *Worker) GenerateFiles() error {
	sshkeyContent, err := os.ReadFile(w.ClusterAsset.SSHKey)
	if err != nil {
		logrus.Debug("Failed to read sshkey content:", err)
		return err
	}

	workerTemplateData, err := ignition.GetTmplData(w.ClusterAsset)
	if err != nil {
		return err
	}
	generateFile := ignition.Common{
		UserName:        w.ClusterAsset.UserName,
		SSHKey:          string(sshkeyContent),
		PassWord:        w.ClusterAsset.Password,
		NodeType:        "worker",
		TmplData:        workerTemplateData,
		EnabledServices: ignition.EnabledServices,
		Config:          &igntypes.Config{},
	}

	// Generate Ignition data
	if err := generateFile.Generate(); err != nil {
		logrus.Errorf("failed to generate %s ignition file: %v", w.ClusterAsset.Worker[0].Hostname, err)
		return err
	}

	ignitionDir := filepath.Join(configmanager.GetPersistDir(), w.ClusterAsset.Cluster_ID, "ignition")

	if err := ignition.SaveFile(generateFile.Config, ignitionDir, WorkerIgnFilename); err != nil {
		return err
	}

	mergerConfig := ignition.GenerateMergeIgnition(w.Bootstrap_baseurl, WorkerIgnFilename)
	if err := ignition.SaveFile(mergerConfig, ignitionDir, workerMergeIgnFilename); err != nil {
		return err
	}

	data, err := ignition.Marshal(generateFile.Config)
	if err != nil {
		logrus.Errorf("failed to Marshal ignition config: %v", err)
		return err
	}

	for i, _ := range w.ClusterAsset.Worker {
		w.ClusterAsset.Worker[i].Ignitions.CreateIgnPath = filepath.Join(ignitionDir, WorkerIgnFilename)
		w.ClusterAsset.Worker[i].Ignitions.MergeIgnPath = filepath.Join(ignitionDir, workerMergeIgnFilename)
		w.ClusterAsset.Worker[i].CreateIgnContent = data
	}

	return nil
}
