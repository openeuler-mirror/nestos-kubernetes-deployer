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
	"errors"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager/globalconfig"
	"os"

	"gopkg.in/yaml.v2"
)

func InitClusterAsset(globalAsset *globalconfig.GlobalConfig, infraAsset InfraAsset, opts *opts.OptionsList) (*ClusterAsset, error) {
	clusterAsset := &ClusterAsset{}

	if opts.File != "" {
		// Parse configuration file.
		configData, err := os.ReadFile(opts.File)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(configData, clusterAsset); err != nil {
			return nil, err
		}
	}

	// cluster info
	setStringValue(&clusterAsset.Cluster_ID, opts.ClusterID, "default cluster id")
	setStringValue(&clusterAsset.Kubernetes_Version, opts.KubeVersion, "default k8s version")

	// bind info
	// infra platform
	switch opts.Platform {
	case "openstack", "Openstack", "OpenStack":
		openstackAsset, ok := infraAsset.(*OpenStackAsset)
		if !ok {
			return nil, errors.New("unsupported platform")
		}
		setStringValue(&clusterAsset.OpenStack_Auth_URL, openstackAsset.Auth_URL, "default openstack auth url")
	}

	// subordinate info
	// master node
	if opts.MasterCount != 0 {
		clusterAsset.Master_Count = opts.MasterCount
	} else if clusterAsset.Master_Count == 0 {
		clusterAsset.Master_Count = 3
	}
	for i := 0; i < opts.MasterCount; i++ {
		master_node := InitNodeAsset(opts)
		clusterAsset.Master_Node = append(clusterAsset.Master_Node, master_node)
	}
	// worker node
	if opts.WorkerCount != 0 {
		clusterAsset.Worker_Count = opts.WorkerCount
	} else if clusterAsset.Worker_Count == 0 {
		clusterAsset.Worker_Count = 3
	}
	for i := 0; i < opts.WorkerCount; i++ {
		worker_node := InitNodeAsset(opts)
		clusterAsset.Worker_Node = append(clusterAsset.Worker_Node, worker_node)
	}

	return clusterAsset, nil
}

// Sets a value of the string type, using the parameter value if the command line argument exists,
// otherwise using the default value.
func setStringValue(target *string, value string, defaultValue string) {
	if value != "" {
		*target = value
	} else if *target == "" {
		*target = defaultValue
	}
}

// Sets a value of type integer, using the parameter value if the command line argument exists,
// otherwise using the default value.
func setIntValue(target *int, value int, defaultValue int) {
	if value != 0 {
		*target = value
	} else if *target == 0 {
		*target = defaultValue
	}
}

// ========== Structure method ==========

type ClusterAsset struct {
	Cluster_ID string

	OpenStack_UserName          string
	OpenStack_Password          string
	OpenStack_Tenant_Name       string
	OpenStack_Auth_URL          string
	OpenStack_Region            string
	OpenStack_Internal_Network  string
	OpenStack_External_Network  string
	OpenStack_Master_IP         []string
	OpenStack_Glance_Name       string
	OpenStack_Availability_Zone string
	OpenStack_UserData          string
	OpenStack_Volume            string

	Master_Count int
	Worker_Count int
	Master_Node  []NodeAsset
	Worker_Node  []NodeAsset

	Kubernetes
	Housekeeper
}

type Kubernetes struct {
	Kubernetes_Version string
	ApiServer_Endpoint string
	Insecure_Registry  string
	Pause_Image        string
	Release_Image_URL  string

	Network
}

type Network struct {
	Service_Subnet        string
	Pod_Subnet            string
	CoreDNS_Image_Version string
}

type Housekeeper struct {
	Operator_Image_URL   string
	Controller_Image_URL string
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
