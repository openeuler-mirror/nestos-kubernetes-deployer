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
	"nestos-kubernetes-deployer/pkg/configmanager/asset/cluster"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/global"

	"github.com/spf13/cobra"
)

func Initial(cmd *cobra.Command) error {
	var globalAsset *global.GlobalAsset
	var clusterAsset *cluster.ClusterAsset

	err := globalAsset.Initial(cmd)
	if err != nil {
		return err
	}

	err = clusterAsset.Initial(cmd)
	if err != nil {
		return err
	}

	return nil
}

func Search() error {
	return nil
}

func Persist() error {
	return nil
}

func Delete() error {
	return nil
}
