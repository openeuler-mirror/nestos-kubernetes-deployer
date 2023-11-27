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
	setStringValue(&clusterAsset.Kubernetes.Kubernetes_Version, opts.KubeVersion, "default k8s version")

	// bind info
	// infra platform
	switch opts.Platform {
	case "openstack", "Openstack", "OpenStack":
		openstackAsset, ok := infraAsset.(*OpenStackAsset)
		if !ok {
			return nil, errors.New("unsupported platform")
		}
		setStringValue(&clusterAsset.OpenStack.Auth_URL, openstackAsset.Auth_URL, "default openstack auth url")
	}

	// subordinate info
	// master node
	if opts.MasterCount != 0 {
		clusterAsset.Master.Count = opts.MasterCount
	} else if clusterAsset.Master.Count == 0 {
		clusterAsset.Master.Count = 3
	}
	for i := 0; i < opts.MasterCount; i++ {
		master_node := &NodeAsset{
			Hostname: opts.MasterConfig[i].Hostname,
			HardwareInfo: HardwareInfo{
				CPU:  opts.MasterConfig[i].CPU,
				RAM:  opts.MasterConfig[i].RAM,
				Disk: opts.MasterConfig[i].Disk,
			},
			UserName: opts.MasterConfig[i].UserName,
			Password: opts.MasterConfig[i].Password,
			SSHKey:   opts.MasterConfig[i].SSHKey,
			IP:       opts.MasterConfig[i].IP,
			Ign_Data: opts.MasterConfig[i].Ign_Data,
		}
		if len(clusterAsset.Master.NodeAsset) == 0 {
			clusterAsset.Master.NodeAsset = append(clusterAsset.Master.NodeAsset, *master_node)
		}
	}
	// worker node
	if opts.WorkerCount != 0 {
		clusterAsset.Worker.Count = opts.WorkerCount
	} else if clusterAsset.Worker.Count == 0 {
		clusterAsset.Worker.Count = 3
	}
	for i := 0; i < opts.WorkerCount; i++ {
		worker_node := &NodeAsset{
			Hostname: opts.WorkerConfig[i].Hostname,
			HardwareInfo: HardwareInfo{
				CPU:  opts.WorkerConfig[i].CPU,
				RAM:  opts.WorkerConfig[i].RAM,
				Disk: opts.WorkerConfig[i].Disk,
			},
			UserName: opts.WorkerConfig[i].UserName,
			Password: opts.WorkerConfig[i].Password,
			SSHKey:   opts.WorkerConfig[i].SSHKey,
			IP:       opts.WorkerConfig[i].IP,
			Ign_Data: opts.WorkerConfig[i].Ign_Data,
		}
		if len(clusterAsset.Worker.NodeAsset) == 0 {
			clusterAsset.Worker.NodeAsset = append(clusterAsset.Worker.NodeAsset, *worker_node)
		}
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

	Platform
	Master
	Worker
	Kubernetes
	Housekeeper
}

type Platform struct {
	OpenStack
	Libvirt
}

type OpenStack struct {
	UserName          string
	Password          string
	Tenant_Name       string
	Auth_URL          string
	Region            string
	Internal_Network  string
	External_Network  string
	Glance_Name       string
	Availability_Zone string
}

type Libvirt struct {
}

type Master struct {
	Count     int
	NodeAsset []NodeAsset
}

type Worker struct {
	Count     int
	NodeAsset []NodeAsset
}

type Kubernetes struct {
	Kubernetes_Version string
	ApiServer_Endpoint string
	Insecure_Registry  string
	Pause_Image        string
	Release_Image_URL  string
	Token              string

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
