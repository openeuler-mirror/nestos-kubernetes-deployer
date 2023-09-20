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
	// system
	MasterHostName = "k8s-master"
	WorkerHostName = "k8s-worker"
	Username       = "root"
	Password       = "$1$yoursalt$UGhjCXAJKpWWpeN8xsF.c/"

	// repo
	Registry = "" //registry.cn-hangzhou.aliyuncs.com/google_containers

	// infra
	Platform = "openstack"

	// master size
	MasterVcpus = 4
	MasterRam   = 8192
	MasterDisk  = 100

	// worker size
	WorkerVcpus = 2
	WorkerRam   = 4096
	WorkerDisk  = 100

	// openstack
	Openstack_UserName         = ""
	Openstack_Password         = ""
	Openstack_Tenant_name      = ""
	Openstack_Auth_url         = ""
	Openstack_Region           = ""
	Openstack_Internal_network = ""
	Openstack_External_network = ""
	Openstack_Master_ip        = []string{"10.1.10.51", "", ""}
	Openstack_Glance_Name      = ""
	Availability_zone          = ""

	KubernetesVersion = "" //v1.23.10
	ServiceSubnet     = "10.96.0.0/16"
	PodSubnet         = "10.100.0.0/16"

	// worker
	APIServerEndpoint = ""
	Token             = "" //abcdef.0123456789abcdef
	TlsBootstrapToken = ""

	MasterNode = "master"
	WorkerNode = "worker"

	Master_Count = 3
	Worker_Count = 2
	SSHKey       = ""

	PauseImageTag   = "" //3.6
	CorednsImageTag = "" //v1.8.6
	ReleaseImageURl = ""
	CertificateKey  = ""
)
