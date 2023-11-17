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

package manager

import (
	"errors"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/cluster"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/global"

	"github.com/spf13/cobra"
)

func Initial(cmd *cobra.Command) error {
	// Init global asset.
	globalAsset := &global.GlobalAsset{}
	if err := globalAsset.Initial(cmd); err != nil {
		return err
	}

	// Init cluster asset.
	clusterAsset := &cluster.ClusterAsset{}
	if err := clusterAsset.Initial(cmd); err != nil {
		return err
	}

	return nil
}

func Search() error {
	return nil
}

func GetGlobalConfig() (*global.GlobalAsset, error) {
	return global.GlobalConfig, nil
}

func GetClusterConfig(clusterID string) (*cluster.ClusterAsset, error) {
	clusterConfig, ok := cluster.ClusterConfig[clusterID]
	if !ok {
		return nil, errors.New("ClusterID not found")
	}

	return clusterConfig, nil
}

func Persist() error {
	// Persist global asset.
	globalConfig, err := GetGlobalConfig()
	if err != nil {
		return err
	}
	if err := globalConfig.Persist(); err != nil {
		return err
	}

	// Persist cluster asset.
	for _, clusterConfig := range cluster.ClusterConfig {
		if err := clusterConfig.Persist(); err != nil {
			return err
		}
	}

	return nil
}

func Delete() error {
	return nil
}
