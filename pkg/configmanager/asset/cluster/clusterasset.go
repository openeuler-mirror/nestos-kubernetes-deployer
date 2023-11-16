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

package cluster

import (
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var ClusterConfig *ClusterAsset

// ========== Package method ==========

func GetClusterConfig() (*ClusterAsset, error) {
	return ClusterConfig, nil
}

// ========== Structure method ==========

type ClusterAsset struct {
	Node
	KubernetesVersion string
}

type Node struct {
	Count int
}

// TODO: Initial inits the cluster asset.
func (ca *ClusterAsset) Initial(cmd *cobra.Command) error {
	configFile, _ := cmd.Flags().GetString("cluster-config-file")

	if configFile != "" {
		// Parse configuration file.
		configData, err := os.ReadFile(configFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(configData, ca)
		if err != nil {
			return err
		}
	}

	kubernetes_version, _ := cmd.Flags().GetString("kubernetes-version")
	if kubernetes_version != "" {
		ca.KubernetesVersion = kubernetes_version
	} else {
		ca.KubernetesVersion = "default k8s version"
	}

	ClusterConfig = ca

	return nil
}

// TODO: Delete deletes the cluster asset.
func (ca *ClusterAsset) Delete() error {
	return nil
}

// TODO: Persist persists the cluster asset.
func (ca *ClusterAsset) Persist() error {
	// TODO: Serialize the cluster asset to json or yaml.
	clusterConfig, err := yaml.Marshal(ca)
	if err != nil {
		return err
	}

	err = os.WriteFile("cluster_config.yaml", clusterConfig, 0644)
	if err != nil {
		return err
	}

	return nil
}
