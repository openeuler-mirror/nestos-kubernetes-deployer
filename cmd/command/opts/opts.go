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

var OptionsList struct {
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
	Master            []*HostConfig
	Worker            []*HostConfig
	NetWork           NetworkConfig
	Housekeeper       HousekeeperConfig
}

type HostConfig struct {
	Name string
	Ip   string
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

var Upgrade UpgradeOpts
