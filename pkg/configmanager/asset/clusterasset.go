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
	"fmt"
	mrand "math/rand"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Sets a value of the string type, using the parameter value if the command line argument exists,
// otherwise using the default value.
func setStringValue(target *string, value string, defaultValue string) {
	if value != "" {
		*target = value
	} else if *target == "" {
		*target = defaultValue
	}
}

// Check whether the value is provided
func checkStringValue(target *string, value string, paramName string) error {
	if value != "" {
		*target = value
	} else if *target == "" {
		return errors.Errorf("%s is unprovided", paramName)
	}

	return nil
}

// Sets a value of type integer, using the parameter value if the command line argument exists,
// otherwise using the default value.
func setIntValue(target *uint, value uint, defaultValue uint) {
	if value != 0 {
		*target = value
	} else if *target == 0 {
		*target = defaultValue
	}
}

// Generate token
func generateToken() (string, error) {
	// Generate a character set for lowercase letters and numbers.
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	charsetLength := len(charset)

	mrand.Seed(time.Now().UnixNano())

	const tokenLength = 6
	token := make([]byte, tokenLength)
	for i := range token {
		token[i] = charset[mrand.Intn(charsetLength)]
	}

	// Add a dot.
	token = append(token, '.')

	// Generate a 16-bit random string.
	const randomPartLength = 16
	randomPart := make([]byte, randomPartLength)
	for i := range randomPart {
		randomPart[i] = charset[mrand.Intn(charsetLength)]
	}

	token = append(token, randomPart...)
	return string(token), nil
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
	Image_Registry     string
	Pause_Image        string
	Release_Image_URL  string
	Token              string
	AdminKubeConfig    string
	CertificateKey     string

	Network
}

type Network struct {
	Service_Subnet        string
	Pod_Subnet            string
	CoreDNS_Image_Version string
}

type Housekeeper struct {
	DeployHousekeeper    bool
	Operator_Image_URL   string
	Controller_Image_URL string
	KubeVersion          string
	EvictPodForce        bool
	MaxUnavailable       uint
	OSImageURL           string
}

func (clusterAsset *ClusterAsset) InitClusterAsset(infraAsset InfraAsset, opts *opts.OptionsList) (*ClusterAsset, error) {
	// cluster info
	setStringValue(&clusterAsset.Cluster_ID, opts.ClusterID, "cluster")

	// bind info
	// infra platform
	clusterAsset.InfraPlatform = infraAsset

	// subordinate info
	// master node
	if len(clusterAsset.Master) == 0 {
		clusterAsset.Master = make([]NodeAsset, 3)
	}
	for i := 0; i < len(clusterAsset.Master); i++ {
		master_node := &clusterAsset.Master[i]

		setStringValue(&master_node.Hostname, opts.Master.Hostname[i], fmt.Sprintf("master%.2d", i))
		setIntValue(&master_node.HardwareInfo.CPU, opts.Master.CPU, 4)
		setIntValue(&master_node.HardwareInfo.RAM, opts.Master.RAM, 8)
		setIntValue(&master_node.HardwareInfo.Disk, opts.Master.Disk, 100)
		setStringValue(&master_node.UserName, opts.Master.UserName, "root")
		setStringValue(&master_node.Password, opts.Master.Password, "")
		setStringValue(&master_node.SSHKey, opts.Master.SSHKey, "")

		opts.Master.IP = make([]string, len(clusterAsset.Master))
		if err := checkStringValue(&master_node.IP, opts.Master.IP[i], fmt.Sprintf("master-ip[%d]", i)); err != nil {
			return nil, err
		}

		opts.Master.IgnFilePath = make([]string, len(clusterAsset.Master))
		if err := checkStringValue(&master_node.Ign_Path, opts.Master.IgnFilePath[i], fmt.Sprintf("master-ign[%d]", i)); err != nil {
			return nil, err
		}
	}
	// worker node
	if len(clusterAsset.Worker) == 0 {
		clusterAsset.Worker = make([]NodeAsset, 3)
	}
	for i := 0; i < len(clusterAsset.Worker); i++ {
		worker_node := &clusterAsset.Worker[i]

		setStringValue(&worker_node.Hostname, opts.Worker.Hostname[i], fmt.Sprintf("worker%.2d", i))
		setIntValue(&worker_node.HardwareInfo.CPU, opts.Worker.CPU, 4)
		setIntValue(&worker_node.HardwareInfo.RAM, opts.Worker.RAM, 8)
		setIntValue(&worker_node.HardwareInfo.Disk, opts.Worker.Disk, 100)
		setStringValue(&worker_node.UserName, opts.Worker.UserName, "root")
		setStringValue(&worker_node.Password, opts.Worker.Password, "")
		setStringValue(&worker_node.SSHKey, opts.Worker.SSHKey, "")

		opts.Worker.IP = make([]string, len(clusterAsset.Worker))
		if err := checkStringValue(&worker_node.IP, opts.Master.IP[i], fmt.Sprintf("worker-ip[%d]", i)); err != nil {
			return nil, err
		}

		opts.Worker.IgnFilePath = make([]string, len(clusterAsset.Worker))
		if err := checkStringValue(&worker_node.Ign_Path, opts.Worker.IgnFilePath[i], fmt.Sprintf("worker-ign[%d]", i)); err != nil {
			return nil, err
		}
	}

	setStringValue(&clusterAsset.Kubernetes.Kubernetes_Version, opts.KubeVersion, "")
	setStringValue(&clusterAsset.Kubernetes.ApiServer_Endpoint, opts.ApiServerEndpoint, "")
	setStringValue(&clusterAsset.Kubernetes.Image_Registry, opts.ImageRegistry, "")
	setStringValue(&clusterAsset.Kubernetes.Pause_Image, opts.PauseImage, "")
	setStringValue(&clusterAsset.Kubernetes.Release_Image_URL, opts.ReleaseImageUrl, "")
	setStringValue(&clusterAsset.Kubernetes.CertificateKey, opts.CertificateKey, "")

	token, err := generateToken()
	if err != nil {
		return nil, err
	}
	setStringValue(&clusterAsset.Kubernetes.Token, opts.Token, token)

	setStringValue(&clusterAsset.Kubernetes.Network.Service_Subnet, opts.NetWork.ServiceSubnet, "")
	setStringValue(&clusterAsset.Kubernetes.Network.Pod_Subnet, opts.NetWork.PodSubnet, "")
	setStringValue(&clusterAsset.Kubernetes.Network.CoreDNS_Image_Version, opts.NetWork.DNS.ImageVersion, "")

	if clusterAsset.Housekeeper.DeployHousekeeper || opts.Housekeeper.DeployHousekeeper {
		if err := checkStringValue(&clusterAsset.Housekeeper.Operator_Image_URL, opts.Housekeeper.OperatorImageUrl, "operator-image-url"); err != nil {
			return nil, err
		}
		if err := checkStringValue(&clusterAsset.Housekeeper.Controller_Image_URL, opts.Housekeeper.ControllerImageUrl, "controller-image-url"); err != nil {
			return nil, err
		}
		if err := checkStringValue(&clusterAsset.Housekeeper.KubeVersion, opts.Housekeeper.KubeVersion, "kubeversion"); err != nil {
			return nil, err
		}
		if opts.Housekeeper.EvictPodForce {
			clusterAsset.Housekeeper.EvictPodForce = true
		}
		setIntValue(&clusterAsset.Housekeeper.MaxUnavailable, opts.Housekeeper.MaxUnavailable, 2)
		if err := checkStringValue(&clusterAsset.Housekeeper.OSImageURL, opts.Housekeeper.OSImageURL, "imageurl"); err != nil {
			return nil, err
		}
	}

	return clusterAsset, nil
}

func (clusterAsset *ClusterAsset) Delete(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return nil
}

func (clusterAsset *ClusterAsset) Persist(dir string) error {
	// Serialize the cluster asset to yaml.
	clusterData, err := yaml.Marshal(clusterAsset)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dir+"/cluster_config.yaml", clusterData, 0644); err != nil {
		return err
	}

	return nil
}
