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

package phases

import (
	"fmt"

	"nestos-kubernetes-deployer/app/apis/nkd"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func NewPrintDefaultNkdConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "print",
		Short: "use this command to print nkd config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
			// return runPrintDefaultConfig()
		},
	}
	cmd.AddCommand(newPrintMasterDefaultConfigCommand())
	cmd.AddCommand(newPrintWorkerDefaultConfigCommand())
	return cmd
}

func newPrintMasterDefaultConfigCommand() *cobra.Command {
	return newCommandPrintDefaultNodeConfig("master")
}

func newPrintWorkerDefaultConfigCommand() *cobra.Command {
	return newCommandPrintDefaultNodeConfig("worker")
}

func newCommandPrintDefaultNodeConfig(node string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s", node),
		Short: fmt.Sprintf("use this command to init %s default config", node),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPrintDefaultConfig(node)
		},
	}
	return cmd
}

func runPrintDefaultConfig(node string) error {
	if node == "master" {
		internalconfig := &nkd.Master{}
		DefaultedStaticMasterConfiguration(internalconfig)
		conf, err := yaml.Marshal(&internalconfig)
		fmt.Println(string(conf))
		if err != nil {
			return err
		}
	} else if node == "worker" {
		internalconfig := &nkd.Worker{}
		DefaultedStaticWorkerConfiguration(internalconfig)
		conf, err := yaml.Marshal(&internalconfig)
		fmt.Println(string(conf))
		if err != nil {
			return err
		}
	}
	return nil
}

func DefaultedStaticWorkerConfiguration(internalconfig *nkd.Worker) *nkd.Worker {
	repo := nkd.Repo{
		Secret:   nkd.Secret,
		Registry: nkd.Registry,
	}
	openstack := nkd.Openstack{
		User_name:        nkd.Openstack_UserName,
		Password:         nkd.Openstack_Password,
		Tenant_name:      nkd.Openstack_Tenant_name,
		Auth_url:         nkd.Openstack_Auth_url,
		Region:           nkd.Openstack_Region,
		Internal_network: nkd.Openstack_Internal_network,
		External_network: nkd.Openstack_External_network,
		Glance:           nkd.Openstack_Glance_Name,
		Flavor:           nkd.Openstack_Flavor_Name,
	}

	system1 := nkd.System{
		Count:    nkd.Master_Count,
		Ips:      nkd.Openstack_Master_ip,
		HostName: nkd.HostName,
		Username: nkd.Username,
		Password: nkd.Password,
	}

	vmsize := nkd.Size{
		Vcpus: nkd.Vcpus,
		Ram:   nkd.Ram,
		Disk:  nkd.Disk,
	}

	infra := nkd.Infra{
		Platform:  nkd.Platform,
		Openstack: openstack,
		Vmsize:    vmsize,
	}

	bootstrapTokenDiscovery := nkd.BootstrapTokenDiscovery{
		APIServerEndpoint:        nkd.APIServerEndpoint,
		Token:                    nkd.Token,
		UnsafeSkipCAVerification: nkd.UnsafeSkipCAVerification,
	}

	discover := nkd.Discovery{
		BootstrapToken:    &bootstrapTokenDiscovery,
		Timeout:           nkd.WorkerDiscoverTimeout,
		TlsBootstrapToken: nkd.TlsBootstrapToken,
	}

	nodeRegistrationOptions := nkd.NodeRegistrationOptions{
		CRISocket:       nkd.CriSocket,
		ImagePullPolicy: nkd.PullPolicy(nkd.ImagePullPolicy),
		Name:            nkd.Name,
		Taints:          nil,
	}

	worker := nkd.WorkerK8s{
		Discovery:        discover,
		CaCertPath:       nkd.CaCertPath,
		NodeRegistration: nodeRegistrationOptions,
	}

	containerdaemon := nkd.ContainerDaemon{
		Pause: nkd.Pauseversion,
	}

	internalconfig.Node = nkd.WorkerNode
	internalconfig.Repo = repo
	internalconfig.System = system1
	internalconfig.Infra = infra
	internalconfig.Worker = worker
	internalconfig.ContainerDaemon = containerdaemon
	return nil
}

// return internal Nkd with static defaults
func DefaultedStaticMasterConfiguration(internalconfig *nkd.Master) *nkd.Master {
	cluster := nkd.Cluster{Name: nkd.NkdClusterName}

	system1 := nkd.System{
		Count:    nkd.Master_Count,
		Ips:      nkd.Openstack_Master_ip,
		HostName: nkd.HostName,
		Username: nkd.Username,
		Password: nkd.Password,
	}

	repo := nkd.Repo{
		Secret:   nkd.Secret,
		Registry: nkd.Registry,
	}

	openstack := nkd.Openstack{
		User_name:        nkd.Openstack_UserName,
		Password:         nkd.Openstack_Password,
		Tenant_name:      nkd.Openstack_Tenant_name,
		Auth_url:         nkd.Openstack_Auth_url,
		Region:           nkd.Openstack_Region,
		Internal_network: nkd.Openstack_Internal_network,
		External_network: nkd.Openstack_External_network,
		Glance:           nkd.Openstack_Glance_Name,
		Flavor:           nkd.Openstack_Flavor_Name,
	}

	vmsize := nkd.Size{
		Vcpus: nkd.Vcpus,
		Ram:   nkd.Ram,
		Disk:  nkd.Disk,
	}

	infra := nkd.Infra{
		Platform:  nkd.Platform,
		Openstack: openstack,
		Vmsize:    vmsize,
	}

	apiServer := nkd.APIServer{
		TimeoutForControlPlane: nkd.TimeoutForControlPlane,
	}

	bootstrapToken := nkd.BootstrapToken{
		Token:  nkd.BootstrapTokensToken,
		Groups: nkd.BootstrapTokensGroups,
		TTL:    nkd.DefaultTokenDuration,
		Usages: nkd.DefaultUsages,
	}

	typemeta := nkd.TypeMeta{
		ApiVersion: nkd.DefaultapiVersion,
	}

	localAPIEndpoint := nkd.APIEndpoint{
		AdvertiseAddress: nkd.AdvertiseAddress,
		BindPort:         nkd.BindPort,
	}

	NodeRegistrationOptions := nkd.NodeRegistrationOptions{
		CRISocket:       nkd.CriSocket,
		ImagePullPolicy: nkd.PullPolicy(nkd.ImagePullPolicy),
		Name:            nkd.Name,
		Taints:          nil,
	}

	ClusterConfiguration := nkd.ClusterConfiguration{
		TypeMeta:          typemeta,
		CertificatesDir:   nkd.CertificatesDir,
		ClusterName:       nkd.ClusterName,
		Etcd:              nkd.Etcd{Local: &nkd.LocalEtcd{DataDir: nkd.LocalDir}},
		ImageRepository:   nkd.ImagePullPolicy,
		KubernetesVersion: nkd.KubernetesVersion,
		Networking:        nkd.Networking{DNSDomain: nkd.DnsDomain, ServiceSubnet: nkd.ServiceSubnet},
		APIServer:         apiServer,
	}

	kubeadm := nkd.Kubeadm{
		ClusterConfiguration: ClusterConfiguration,
		TypeMeta:             typemeta,
		BootstrapTokens:      []nkd.BootstrapToken{bootstrapToken},
		LocalAPIEndpoint:     localAPIEndpoint,
		NodeRegistration:     NodeRegistrationOptions,
	}

	containerdaemon := nkd.ContainerDaemon{
		Pause: nkd.Pauseversion,
	}

	internalconfig.Node = nkd.MasterNode
	internalconfig.Kubeadm = kubeadm
	internalconfig.Cluster = cluster
	internalconfig.Infra = infra
	internalconfig.System = system1
	internalconfig.Repo = repo
	internalconfig.ContainerDaemon = containerdaemon

	return internalconfig
}
