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
			return runPrintDefaultConfig()
		},
	}
	return cmd
}

func runPrintDefaultConfig() error {
	internalconfig := &nkd.Nkd{}
	DefaultedStaticInitConfiguration(internalconfig)
	conf, err := yaml.Marshal(&internalconfig)
	fmt.Println(string(conf))
	if err != nil {
		return err
	}
	return err
}

// return internal Nkd with static defaults
func DefaultedStaticInitConfiguration(internalconfig *nkd.Nkd) *nkd.Nkd {
	cluster := nkd.Cluster{Name: nkd.NkdClusterName}

	system1 := nkd.System{Hostname: nkd.Hostname1,
		Username: nkd.Username,
		Password: nkd.Password,
	}

	repo := nkd.Repo{Secret: nkd.Secret,
		Registry: nkd.Registry,
	}

	openstack := nkd.Openstack{
		User_name:            nkd.Openstack_UserName,
		Password:             nkd.Openstack_Password,
		Tenant_name:          nkd.Openstack_Tenant_name,
		Auth_url:             nkd.Openstack_Auth_url,
		Region:               nkd.Openstack_Region,
		Master_instance_name: nkd.Openstack_MasterNodeName,
		Worker_instance_name: nkd.Openstack_WorkerNodeName,
		Internal_network:     nkd.Openstack_Internal_network,
		External_network:     nkd.Openstack_External_network,
		Master_ip:            nkd.Openstack_Master_ip,
		Worker_ip:            nkd.Openstack_Worker_ip,
		Glance:               nkd.Openstack_Glance_Name,
		Flavor:               nkd.Openstack_Flavor_Name,
	}

	vmsize := nkd.Vmsize{
		Master: nkd.Size{
			Vcpus: nkd.Vcpus,
			Ram:   nkd.Ram,
			Disk:  nkd.Disk,
		},
		Worker: nkd.Size{
			Vcpus: nkd.Vcpus,
			Ram:   nkd.Ram,
			Disk:  nkd.Disk,
		},
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
	internalconfig.Kubeadm = kubeadm
	internalconfig.Cluster = cluster
	internalconfig.Infra = infra
	internalconfig.System = []nkd.System{system1}
	internalconfig.Repo = repo

	return internalconfig
}
