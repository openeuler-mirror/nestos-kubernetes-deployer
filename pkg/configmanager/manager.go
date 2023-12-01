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

package configmanager

import (
	"errors"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/globalconfig"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Set global data
var GlobalConfig *globalconfig.GlobalConfig
var ClusterAsset = map[string]*asset.ClusterAsset{}

// var InfraAsset = map[string]*asset.InfraAsset{}

func Initial(opts *opts.OptionsList) error {
	// Init global asset
	globalConfig, err := globalconfig.InitGlobalConfig(opts)
	if err != nil {
		return err
	}
	GlobalConfig = globalConfig

	fileData := &asset.ClusterAsset{}
	if opts.ClusterConfigFile != "" {
		// Parse configuration file.
		configData, err := os.ReadFile(opts.ClusterConfigFile)
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal(configData, fileData); err != nil {
			return err
		}
	}

	// Init infra asset
	infraAsset, err := asset.InitInfraAsset(fileData, opts)
	if err != nil {
		return err
	}

	// Init cluster asset
	clusterAsset, err := fileData.InitClusterAsset(infraAsset, opts)
	if err != nil {
		return err
	}
	ClusterAsset[fileData.Cluster_ID] = clusterAsset

	return nil
}

func GetGlobalConfig() (*globalconfig.GlobalConfig, error) {
	return GlobalConfig, nil
}

func GetPersistDir() string {
	return GlobalConfig.PersistDir
}

func GetClusterConfig(clusterID string) (*asset.ClusterAsset, error) {
	clusterConfig, ok := ClusterAsset[clusterID]
	if !ok {
		return nil, errors.New("ClusterID not found")
	}

	return clusterConfig, nil
}

func Persist() error {
	// Get persist dir
	persistDir := GetPersistDir()

	// Persist cluster
	for _, clusterAsset := range ClusterAsset {
		clusterDir := filepath.Join(persistDir, clusterAsset.Cluster_ID)
		if err := os.MkdirAll(clusterDir, 0644); err != nil {
			return err
		}

		if err := clusterAsset.Persist(clusterDir); err != nil {
			return err
		}
	}

	return nil
}

func Delete(clusterID string) error {
	// Get persist dir
	persistDir := GetPersistDir()

	clusterAsset, err := GetClusterConfig(clusterID)
	if err != nil {
		return err
	}

	if err := clusterAsset.Delete(filepath.Join(persistDir, clusterID)); err != nil {
		return err
	}

	return nil
}
