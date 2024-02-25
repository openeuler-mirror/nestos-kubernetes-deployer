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
	"nestos-kubernetes-deployer/pkg/utils"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
func setUIntValue(target *uint, value uint, defaultValue uint) {
	if value != 0 {
		*target = value
	} else if *target == 0 {
		*target = defaultValue
	}
}

// Generate token
func generateToken() string {
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
	return string(token)
}

func setMasterConfigs(mc []NodeAsset, opts *opts.MasterConfig) []NodeAsset {
	var confs []NodeAsset
	if len(mc) >= len(opts.IP) {
		confs = append(confs, mc...)
		for i, v := range opts.IP {
			confs[i].IP = v
			confs[i].Hostname = opts.Hostname[i]
		}
	} else {
		confs = append(confs, mc...)
		for i, _ := range opts.IP {
			if i < len(mc) {
				confs[i].IP = opts.IP[i]
				confs[i].Hostname = opts.Hostname[i]
				continue
			}
			confs = append(confs, NodeAsset{
				IP:       opts.IP[i],
				Hostname: opts.Hostname[i],
				HardwareInfo: HardwareInfo{
					CPU:  4,
					RAM:  8192,
					Disk: 50,
				},
			})
		}
	}
	return confs
}

func setWorkerHostname(wc []NodeAsset, opts *opts.WorkerConfig) []NodeAsset {
	var confs []NodeAsset
	if len(wc) >= len(opts.Hostname) {
		confs = append(confs, wc...)
		for i, v := range opts.Hostname {
			confs[i].Hostname = v
		}
	} else {
		confs = append(confs, wc...)
		for i, _ := range opts.Hostname {
			if i < len(wc) {
				confs[i].Hostname = opts.Hostname[i]
				continue
			}
			confs = append(confs, NodeAsset{
				IP:       "",
				Hostname: opts.Hostname[i],
				HardwareInfo: HardwareInfo{
					CPU:  4,
					RAM:  8192,
					Disk: 50,
				},
			})
		}
	}
	return confs
}

// ========== Structure method ==========

type ClusterAsset struct {
	Cluster_ID   string
	Architecture string
	Platform     string
	InfraPlatform
	UserName string
	Password string
	SSHKey   string
	Master   []NodeAsset
	Worker   []NodeAsset
	Runtime  string `yaml:"runtime"` //后续考虑增加os层面的配置管理，并将runtime放入OS层面的配置中
	Kubernetes
	Housekeeper
	CertAsset
}

type InfraPlatform interface {
}

type Kubernetes struct {
	KubernetesVersion    string `yaml:"kubernetes-version"`
	KubernetesAPIVersion string `yaml:"kubernetes-apiversion"`
	ApiServerEndpoint    string `yaml:"apiserver-endpoint"`
	ImageRegistry        string `yaml:"image-registry"`
	PauseImage           string `yaml:"pause-image"`
	ReleaseImageURL      string `yaml:"release-image-url"`
	Token                string
	AdminKubeConfig      string
	CertificateKey       string

	Network
}

type Network struct {
	ServiceSubnet string `yaml:"service-subnet"`
	PodSubnet     string `yaml:"pod-subnet"`
	Plugin        string
}

type Housekeeper struct {
	DeployHousekeeper  bool
	OperatorImageUrl   string
	ControllerImageUrl string
	KubeVersion        string `json:"-" yaml:"-"`
	EvictPodForce      bool   `json:"-" yaml:"-"`
	MaxUnavailable     uint   `json:"-" yaml:"-"`
	OSImageURL         string `json:"-" yaml:"-"`
}

func (clusterAsset *ClusterAsset) InitClusterAsset(infraAsset InfraAsset, opts *opts.OptionsList) (*ClusterAsset, error) {
	// bind info
	// infra platform
	clusterAsset.InfraPlatform = infraAsset

	cf, err := GetDefaultClusterConfig(clusterAsset.Architecture)
	if err != nil {
		return nil, err
	}

	// set node config
	if len(clusterAsset.Master) == 0 {
		clusterAsset.Master = append(clusterAsset.Master, cf.Master...)
	}
	if len(opts.Master.Hostname) != len(opts.Master.IP) {
		return nil, fmt.Errorf("the number of configuration parameters master hostname and ip should be the same")
	}
	if len(opts.Master.IP) != 0 {
		clusterAsset.Master = setMasterConfigs(clusterAsset.Master, &opts.Master)
	}
	if opts.Master.CPU != 0 {
		for i, _ := range clusterAsset.Master {
			clusterAsset.Master[i].HardwareInfo.CPU = opts.Master.CPU
		}
	}
	if opts.Master.RAM != 0 {
		for i, _ := range clusterAsset.Master {
			clusterAsset.Master[i].HardwareInfo.RAM = opts.Master.RAM
		}
	}
	if opts.Master.Disk != 0 {
		for i, _ := range clusterAsset.Master {
			clusterAsset.Master[i].HardwareInfo.Disk = opts.Master.Disk
		}
	}

	// set worker node config
	if len(clusterAsset.Worker) == 0 {
		clusterAsset.Worker = append(clusterAsset.Worker, cf.Worker...)
	}
	// set worker hostname
	if len(opts.Worker.Hostname) != 0 {
		clusterAsset.Worker = setWorkerHostname(clusterAsset.Worker, &opts.Worker)
	}
	// set worker IPs
	if len(opts.Worker.IP) != 0 {
		if len(opts.Worker.Hostname) != len(opts.Worker.IP) {
			return nil, fmt.Errorf("the number of configuration parameters worker hostname and ip should be the same")
		}
		for i, _ := range opts.Worker.IP {
			clusterAsset.Worker[i].IP = opts.Worker.IP[i]
		}
	}

	if opts.Worker.CPU != 0 {
		for i, _ := range clusterAsset.Worker {
			clusterAsset.Worker[i].HardwareInfo.CPU = opts.Worker.CPU
		}
	}
	if opts.Worker.RAM != 0 {
		for i, _ := range clusterAsset.Worker {
			clusterAsset.Worker[i].HardwareInfo.RAM = opts.Worker.RAM
		}
	}
	if opts.Worker.Disk != 0 {
		for i, _ := range clusterAsset.Worker {
			clusterAsset.Worker[i].HardwareInfo.Disk = opts.Worker.Disk
		}
	}

	// cluster info
	setStringValue(&clusterAsset.Cluster_ID, opts.ClusterID, cf.Cluster_ID)
	setStringValue(&clusterAsset.UserName, opts.UserName, cf.UserName)
	setStringValue(&clusterAsset.Password, opts.Password, cf.Password)
	setStringValue(&clusterAsset.SSHKey, opts.SSHKey, cf.SSHKey)
	setStringValue(&clusterAsset.Kubernetes.KubernetesVersion, opts.KubeVersion, cf.KubernetesVersion)
	setStringValue(&clusterAsset.Runtime, opts.Runtime, cf.Runtime)
	setStringValue(&clusterAsset.Kubernetes.ApiServerEndpoint, opts.ApiServerEndpoint, cf.ApiServerEndpoint)
	setStringValue(&clusterAsset.Kubernetes.ImageRegistry, opts.ImageRegistry, cf.ImageRegistry)
	setStringValue(&clusterAsset.Kubernetes.PauseImage, opts.PauseImage, cf.PauseImage)
	setStringValue(&clusterAsset.Kubernetes.ReleaseImageURL, opts.ReleaseImageUrl, cf.ReleaseImageURL)
	setStringValue(&clusterAsset.Kubernetes.CertificateKey, opts.CertificateKey, opts.CertificateKey)
	setStringValue(&clusterAsset.Kubernetes.Token, opts.Token, cf.Token)
	setStringValue(&clusterAsset.Kubernetes.Network.ServiceSubnet, opts.NetWork.ServiceSubnet, cf.ServiceSubnet)
	setStringValue(&clusterAsset.Kubernetes.Network.PodSubnet, opts.NetWork.PodSubnet, cf.Network.PodSubnet)
	setStringValue(&clusterAsset.Kubernetes.Network.Plugin, opts.NetWork.Plugin, cf.Network.Plugin)
	apiVersion, err := utils.GetKubernetesApiVersion(opts.KubernetesAPIVersion)
	if err != nil {
		logrus.Errorf("Error getting kubernetes api version: %v\n", err)
		return nil, err
	}
	setStringValue(&clusterAsset.Kubernetes.KubernetesAPIVersion, apiVersion, cf.KubernetesAPIVersion)

	if clusterAsset.Housekeeper.DeployHousekeeper || opts.Housekeeper.DeployHousekeeper {
		setStringValue(&clusterAsset.Housekeeper.OperatorImageUrl, opts.Housekeeper.OperatorImageUrl, cf.OperatorImageUrl)
		setStringValue(&clusterAsset.Housekeeper.ControllerImageUrl, opts.Housekeeper.ControllerImageUrl, cf.ControllerImageUrl)
		setStringValue(&clusterAsset.Housekeeper.KubeVersion, opts.Housekeeper.KubeVersion, "")
		setStringValue(&clusterAsset.Housekeeper.OSImageURL, opts.Housekeeper.OSImageURL, "")
		setUIntValue(&clusterAsset.Housekeeper.MaxUnavailable, opts.Housekeeper.MaxUnavailable, cf.MaxUnavailable)
		clusterAsset.Housekeeper.EvictPodForce = opts.Housekeeper.EvictPodForce
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

func GetDefaultClusterConfig(arch string) (*ClusterAsset, error) {
	var (
		OperatorImageUrl   string
		ControllerImageUrl string
	)
	switch arch {
	case "amd64", "x86_64":
		OperatorImageUrl = "hub.oepkgs.net/nestos/housekeeper/amd64/housekeeper-operator-manager:0.1.0"
		ControllerImageUrl = "hub.oepkgs.net/nestos/housekeeper/amd64/housekeeper-controller-manager:0.1.0"
	case "arm64", "aarch64":
		OperatorImageUrl = "hub.oepkgs.net/nestos/housekeeper/arm64/housekeeper-operator-manager:0.1.0"
		ControllerImageUrl = "hub.oepkgs.net/nestos/housekeeper/arm64/housekeeper-controller-manager:0.1.0"
	default:
		return nil, errors.New("unsupported architecture")
	}

	return &ClusterAsset{
		Cluster_ID:   "cluster",
		Architecture: arch,
		Platform:     "libvirt",
		UserName:     "root",
		Password:     "$1$yoursalt$UGhjCXAJKpWWpeN8xsF.c/",
		SSHKey:       utils.GetDefaultPubKeyPath(),
		Master: []NodeAsset{
			{
				Hostname: "k8s-master01",
				HardwareInfo: HardwareInfo{
					CPU:  4,
					RAM:  8192,
					Disk: 50,
				},
				IP: "192.168.132.11",
			},
		},
		Worker: []NodeAsset{
			{
				Hostname: "k8s-worker01",
				HardwareInfo: HardwareInfo{
					CPU:  4,
					RAM:  8192,
					Disk: 50,
				},
				IP: "",
			},
		},
		Runtime: "isulad",
		Kubernetes: Kubernetes{
			KubernetesVersion:    "v1.23.10",
			KubernetesAPIVersion: "v1beta3",
			ApiServerEndpoint:    utils.GetApiServerEndpoint("192.168.132.11"),
			ImageRegistry:        "k8s.gcr.io",
			PauseImage:           "pause:3.6",
			ReleaseImageURL:      "",
			Token:                generateToken(),
			CertificateKey:       "a301c9c55596c54c5d4c7173aa1e3b6fd304130b0c703bb23149c0c69f94b8e0",
			Network: Network{
				ServiceSubnet: "10.96.0.0/16",
				PodSubnet:     "10.244.0.0/16",
				Plugin:        "https://projectcalico.docs.tigera.io/archive/v3.22/manifests/calico.yaml",
			},
		},
		Housekeeper: Housekeeper{
			OperatorImageUrl:   OperatorImageUrl,
			ControllerImageUrl: ControllerImageUrl,
			MaxUnavailable:     2,
		},
	}, nil
}
