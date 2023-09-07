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

package nkd

// define kubeadm default config
var (
	// cluster
	NkdClusterName = "example nkd cluster"

	// system
	MasterHostName = "master"
	WorkerHostName = "worker"
	Username       = "root"
	Password       = "$1$yoursalt$UGhjCXAJKpWWpeN8xsF.c/"

	// repo
	Secret   = []map[string]string{{"repousre": ""}}
	Registry = "registry.cn-hangzhou.aliyuncs.com/google_containers"

	// infra
	Platform = "openstack"

	// size
	Vcpus = 4
	Ram   = 8192
	Disk  = 128

	// openstack
	Openstack_UserName         = ""
	Openstack_Password         = ""
	Openstack_Tenant_name      = ""
	Openstack_Auth_url         = ""
	Openstack_Region           = ""
	Openstack_Internal_network = ""
	Openstack_External_network = ""
	Openstack_Master_ip        = []string{"10.1.10.51", "10.1.10.52", "10.1.10.53"}
	Openstack_Flavor_Name      = ""
	Openstack_Glance_Name      = ""

	// bootstrapTokens
	BootstrapTokensGroups = []string{"system:bootstrappers:kubeadm:default-node-token"}
	BootstrapTokensToken  = "abcdef.0123456789abcdef"
	// DefaultTokenDuration  = 24 * time.Hour
	DefaultTokenDuration = "24h0m0s"
	DefaultUsages        = []string{"signing", "authentication"}

	// TypeMeta
	DefaultapiVersion = "kubeadm.k8s.io/v1beta3"
	Kind              = "InitConfiguration"

	// localAPIEndpoint
	AdvertiseAddress       = "1.2.3.4"
	BindPort         int32 = 6643

	// nodeRegistration
	CriSocket       = "/var/run/isulad.sock"
	ImagePullPolicy = "IfNotPresent"
	Name            = "node"
	Taints          = []Taint{}

	// apiServer
	TimeoutForControlPlane = "4m0s"

	// ClusterConfiguration
	CertificatesDir = "/etc/kubernetes/pki"
	ClusterName     = "kubernetes"

	// etcd
	LocalDir          = "/var/lib/etcd"
	ImageRepository   = "registry.cn-hangzhou.aliyuncs.com/google_containers"
	KubernetesVersion = "1.23.10"
	DnsDomain         = "cluster.local"
	ServiceSubnet     = "10.96.0.0/16"
	PodSubnet         = "10.100.0.0/16"

	// worker
	APIServerEndpoint        = "10.1.10.51:6443"
	Token                    = "abcdef.0123456789abcdef"
	UnsafeSkipCAVerification = true
	WorkerDiscoverTimeout    = "5m0s"
	TlsBootstrapToken        = "abcdef.0123456789abcdef"
	CaCertPath               = "/etc/kubernetes/pki/ca.crt"

	MasterNode = "master"
	WorkerNode = "worker"

	Master_Count = 3
	Worker_Count = 3
	SSHKey       = ""

	// containerdaemon
	PauseImageTag   = "3.6"
	CorednsImageTag = "v1.8.6"
	ReleaseImageURl = ""
)
