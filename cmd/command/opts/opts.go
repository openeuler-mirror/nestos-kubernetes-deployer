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

package opts

var (
	RootOptDir string
)

var Opts OptionsList

type OptionsList struct {
	File              string
	ClusterID         string
	Password          string
	SSHKey            string
	Platform          string
	DeployConfig      string
	ApiServerEndpoint string
	InsecureRegistry  string
	PauseImage        string
	ReleaseImageUrl   string
	KubeVersion       string
	MasterCount       int
	MasterConfig      []NodeConfig
	WorkerCount       int
	WorkerConfig      []NodeConfig
	NetWork           NetworkConfig
	Housekeeper       HousekeeperConfig
	Upgrade           UpgradeOpts
}

type NodeConfig struct {
	Hostname string
	CPU      int
	RAM      int
	Disk     int
	UserName string
	Password string
	SSHKey   string
	IP       string
	Ign_Data string
}

type NetworkConfig struct {
	ServiceSubnet string
	PodSubnet     string
	DNS           DnsConfig
}

type DnsConfig struct {
	ImageVersion string //coredns
}

type HousekeeperConfig struct {
	OperatorImageUrl   string
	ControllerImageUrl string
}

type UpgradeOpts struct {
	KubeVersion    string
	EvictPodForce  bool
	MaxUnavailable int
	KubeConfigFile string
	OSImageURL     string
}