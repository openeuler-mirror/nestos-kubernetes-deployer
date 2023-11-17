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

package asset

import (
	"nestos-kubernetes-deployer/pkg/configmanager/globalconfig"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// ========== Structure method ==========

type ClusterAsset struct {
	// cluster info
	ClusterID         string
	KubernetesVersion string

	// bind info
	OpenStackAsset

	// subordinate info
	Master_Count int
	Worker_Count int
	Master_Node  []NodeAsset
	Worker_Node  []NodeAsset
}

func InitClusterAsset(globalAsset *globalconfig.GlobalAsset, cmd *cobra.Command) (*ClusterAsset, error) {
	clusterAsset := &ClusterAsset{}

	configFile, _ := cmd.Flags().GetString("cluster-config-file")
	if configFile != "" {
		// Parse configuration file.
		configData, err := os.ReadFile(configFile)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(configData, clusterAsset); err != nil {
			return nil, err
		}
	}

	// cluster info
	setStringValue(&clusterAsset.ClusterID, cmd, "clusterid", "default cluster id")
	setStringValue(&clusterAsset.KubernetesVersion, cmd, "kubernetes-version", "default k8s version")

	// subordinate info
	// master node
	master_count, _ := cmd.Flags().GetInt("master-count")
	if master_count != 0 {
		clusterAsset.Master_Count = master_count
	} else if clusterAsset.Master_Count == 0 {
		clusterAsset.Master_Count = 3
	}
	for i := 0; i < master_count; i++ {
		master_node := InitNodeAsset(cmd, "master")
		clusterAsset.Master_Node = append(clusterAsset.Master_Node, master_node)
	}
	// worker node
	worker_count, _ := cmd.Flags().GetInt("worker-count")
	if worker_count != 0 {
		clusterAsset.Worker_Count = worker_count
	} else if clusterAsset.Worker_Count == 0 {
		clusterAsset.Worker_Count = 3
	}
	for i := 0; i < worker_count; i++ {
		worker_node := InitNodeAsset(cmd, "worker")
		clusterAsset.Worker_Node = append(clusterAsset.Worker_Node, worker_node)
	}

	return clusterAsset, nil
}

// Sets a value of the string type, using the parameter value if the command line argument exists,
// otherwise using the default value.
func setStringValue(target *string, cmd *cobra.Command, flagName string, defaultValue string) {
	value, _ := cmd.Flags().GetString(flagName)
	if value != "" {
		*target = value
	} else if *target == "" {
		*target = defaultValue
	}
}

// Sets a value of type integer, using the parameter value if the command line argument exists,
// otherwise using the default value.
func setIntValue(target *int, cmdValue int, defaultValue int) {
	if cmdValue != 0 {
		*target = cmdValue
	} else if *target == 0 {
		*target = defaultValue
	}
}

// TODO: Delete deletes the cluster asset.
func (ca *ClusterAsset) Delete() error {
	return nil
}

// TODO: Persist persists the cluster asset.
func (ca *ClusterAsset) Persist() error {
	// Serialize the cluster asset to yaml.
	clusterData, err := yaml.Marshal(ca)
	if err != nil {
		return err
	}

	err = os.WriteFile("cluster_config.yaml", clusterData, 0644)
	if err != nil {
		return err
	}

	return nil
}
