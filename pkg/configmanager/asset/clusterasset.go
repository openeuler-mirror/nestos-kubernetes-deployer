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
	"fmt"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager/globalconfig"
	"os"

	"github.com/clarketm/json"

	"gopkg.in/yaml.v2"
)

func InitClusterAsset(globalAsset *globalconfig.GlobalConfig, infraAsset InfraAsset, opts *opts.OptionsList) (*ClusterAsset, error) {
	clusterAsset := &ClusterAsset{}

	if opts.ClusterConfigFile != "" {
		// Parse configuration file.
		configData, err := os.ReadFile(opts.ClusterConfigFile)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(configData, clusterAsset); err != nil {
			return nil, err
		}
	}

	// cluster info
	setStringValue(&clusterAsset.Cluster_ID, opts.ClusterID, "cluster")

	// bind info
	// infra platform
	switch opts.Platform {
	case "openstack", "Openstack", "OpenStack":
		openstackAsset, ok := infraAsset.(*OpenStackAsset)
		if !ok {
			return nil, errors.New("unsupported platform")
		}
		clusterAsset.InfraPlatform = openstackAsset
	case "libvirt", "Libvirt":
		libvirtAsset, ok := infraAsset.(*LibvirtAsset)
		if !ok {
			return nil, errors.New("unsupported platform")
		}
		clusterAsset.InfraPlatform = libvirtAsset
	default:
		return nil, errors.New("unsupported platform")
	}

	// subordinate info
	// master node
	for i, master_node := range clusterAsset.Master {
		setStringValue(&master_node.Hostname, opts.Master.Hostname[i], fmt.Sprintf("master%.2d", i))
		setIntValue(&master_node.HardwareInfo.CPU, opts.Master.CPU, 4)
		setIntValue(&master_node.HardwareInfo.RAM, opts.Master.RAM, 8)
		setIntValue(&master_node.HardwareInfo.Disk, opts.Master.Disk, 100)
		setStringValue(&master_node.UserName, opts.Master.UserName, "root")
		setStringValue(&master_node.Password, opts.Master.Password, "")
		setStringValue(&master_node.SSHKey, opts.Master.SSHKey, "")
		setStringValue(&master_node.IP, opts.Master.IP[i], "")

		if opts.Master.IgnFilePath[i] != "" {
			ignData, err := json.Marshal(opts.Master.IgnFilePath[i])
			if err != nil {
				return nil, err
			}
			master_node.Ign_Data = ignData
		}

		clusterAsset.Master = append(clusterAsset.Master, master_node)
	}
	// worker node
	for i, worker_node := range clusterAsset.Worker {
		setStringValue(&worker_node.Hostname, opts.Worker.Hostname[i], fmt.Sprintf("worker%.2d", i))
		setIntValue(&worker_node.HardwareInfo.CPU, opts.Worker.CPU, 4)
		setIntValue(&worker_node.HardwareInfo.RAM, opts.Worker.RAM, 8)
		setIntValue(&worker_node.HardwareInfo.Disk, opts.Worker.Disk, 100)
		setStringValue(&worker_node.UserName, opts.Worker.UserName, "root")
		setStringValue(&worker_node.Password, opts.Worker.Password, "")
		setStringValue(&worker_node.SSHKey, opts.Worker.SSHKey, "")
		setStringValue(&worker_node.IP, opts.Worker.IP[i], "")

		if opts.Worker.IgnFilePath[i] != "" {
			ignData, err := json.Marshal(opts.Worker.IgnFilePath[i])
			if err != nil {
				return nil, err
			}
			worker_node.Ign_Data = ignData
		}

		clusterAsset.Worker = append(clusterAsset.Worker, worker_node)
	}

	setStringValue(&clusterAsset.Kubernetes.Kubernetes_Version, opts.KubeVersion, "")
	setStringValue(&clusterAsset.Kubernetes.ApiServer_Endpoint, opts.ApiServerEndpoint, "")
	setStringValue(&clusterAsset.Kubernetes.Insecure_Registry, opts.InsecureRegistry, "")
	setStringValue(&clusterAsset.Kubernetes.Pause_Image, opts.PauseImage, "")
	setStringValue(&clusterAsset.Kubernetes.Release_Image_URL, opts.ReleaseImageUrl, "")
	setStringValue(&clusterAsset.Kubernetes.Network.Service_Subnet, opts.NetWork.ServiceSubnet, "")
	setStringValue(&clusterAsset.Kubernetes.Network.Pod_Subnet, opts.NetWork.PodSubnet, "")
	setStringValue(&clusterAsset.Kubernetes.Network.CoreDNS_Image_Version, opts.NetWork.DNS.ImageVersion, "")

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
	Platform   string

	InfraPlatform
	Master []NodeAsset
	Worker []NodeAsset
	Kubernetes
	Housekeeper
	CertAsset
}

type InfraPlatform interface {
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
	KubeVersion          string
	EvictPodForce        bool
	MaxUnavailable       int
	OSImageURL           string
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
