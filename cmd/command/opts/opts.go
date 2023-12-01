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
	ClusterConfigFile string
	KubeConfigFile    string
	NKD               NKDConfig
	InfraPlatform

	ClusterID string
	Platform  string

	Master      NodeConfig
	Worker      NodeConfig
	ExtendCount int

	ApiServerEndpoint string
	InsecureRegistry  string
	PauseImage        string
	ReleaseImageUrl   string
	KubeVersion       string
	Token             string

	NetWork NetworkConfig
	Housekeeper
}

type NKDConfig struct {
	Log_Level string
}

type InfraPlatform struct {
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

type NodeConfig struct {
	Hostname    []string
	CPU         int
	RAM         int
	Disk        int
	UserName    string
	Password    string
	SSHKey      string
	IP          []string
	IgnFilePath []string
}

type NetworkConfig struct {
	ServiceSubnet string
	PodSubnet     string
	DNS           DnsConfig
}

type DnsConfig struct {
	ImageVersion string //coredns
}

type Housekeeper struct {
	OperatorImageUrl   string
	ControllerImageUrl string
	KubeVersion        string
	EvictPodForce      bool
	MaxUnavailable     int
	OSImageURL         string
}
