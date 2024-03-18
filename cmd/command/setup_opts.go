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

package command

import (
	"nestos-kubernetes-deployer/cmd/command/opts"

	"github.com/spf13/cobra"
)

func SetupDeployCmdOpts(deployCmd *cobra.Command) {
	flags := deployCmd.Flags()
	flags.StringVarP(&opts.Opts.ClusterConfigFile, "file", "f", "", "Location of cluster deploy config file")
	flags.StringVarP(&opts.Opts.ClusterID, "cluster-id", "", "", "Cluster ID")
	flags.StringVar(&opts.Opts.Arch, "arch", "", "kubernetes Cluster deployment architecture")
	flags.StringVarP(&opts.Opts.Platform, "platform", "", "", "Select the infrastructure platform to deploy the cluster")
	flags.StringVarP(&opts.Opts.UserName, "username", "", "", "User to login the node")
	flags.StringVarP(&opts.Opts.Password, "password", "", "", "Password to login the node")
	flags.StringVarP(&opts.Opts.SSHKey, "sshkey", "", "", "Set nodes ssh private key for authentication")
	flags.StringArrayVarP(&opts.Opts.Master.Hostname, "master-hostname", "", []string{}, "Set master hostnames")
	flags.UintVar(&opts.Opts.Master.CPU, "master-cpu", 0, "Set master CPU (units: cores)")
	flags.UintVar(&opts.Opts.Master.RAM, "master-ram", 0, "Set master RAM (units: MB)")
	flags.UintVar(&opts.Opts.Master.Disk, "master-disk", 0, "Set master disk size (units: GB)")
	flags.StringArrayVarP(&opts.Opts.Master.IP, "master-ips", "", []string{}, "Set master IPs")
	flags.StringArrayVarP(&opts.Opts.Worker.Hostname, "worker-hostname", "", []string{}, "Set worker hostnames")
	flags.UintVar(&opts.Opts.Worker.CPU, "worker-cpu", 0, "Set worker CPU (units: cores)")
	flags.UintVar(&opts.Opts.Worker.RAM, "worker-ram", 0, "Set worker RAM (units: MB)")
	flags.UintVar(&opts.Opts.Worker.Disk, "worker-disk", 0, "Set worker disk size (units: GB)")
	flags.StringArrayVarP(&opts.Opts.Worker.IP, "worker-ips", "", []string{}, "Set worker IPs")
	flags.StringVarP(&opts.Opts.Runtime, "runtime", "", "", "Specify container runtime type (docker, isulad, containerd or crio)")
	flags.StringVarP(&opts.Opts.ImageRegistry, "image-registry", "", "", "Specify the registry address for pulling the Kubernetes component container image")
	flags.StringVarP(&opts.Opts.PauseImage, "pause-image", "", "", "Specify the image for the pause container")
	flags.StringVarP(&opts.Opts.ReleaseImageUrl, "release-image-url", "", "", "Specify the URL of the NestOS container image that contains the Kubernetes component. Only supports the qcow2 format.")
	flags.StringVarP(&opts.Opts.KubeVersion, "kubeversion", "", "", "Specify the version of Kubernetes to deploy")
	flags.UintVarP(&opts.Opts.KubernetesAPIVersion, "kubernetes-apiversion", "", 0,
		"Sets the Kubernetes API version. Acceptable reference values:\n"+
			"  - 1 for Kubernetes versions < v1.15.0,\n"+
			"  - 2 for Kubernetes versions between v1.15.0 and v1.22.0,\n"+
			"  - 3 for Kubernetes versions >= v1.22.0")
	flags.StringVarP(&opts.Opts.Token, "token", "", "", "Specify the authentication token for accessing resources")
	flags.StringVarP(&opts.Opts.CertificateKey, "CertificateKey", "", "", "Specifies the certificate key to be added to the master node")
	flags.StringVarP(&opts.Opts.NetWork.ServiceSubnet, "service-subnet", "", "", "Specify the subnet for Kubernetes services")
	flags.StringVarP(&opts.Opts.NetWork.PodSubnet, "pod-subnet", "", "", "Specify the subnet for Kubernetes Pods.")
	flags.StringVarP(&opts.Opts.NetWork.Plugin, "network-plugin-url", "", "", "URL for the network plugin")
	flags.StringVarP(&opts.Opts.Housekeeper.ControllerImageUrl, "controller-image-url", "", "", "Specify the URL of the container image for the housekeeper controller component")
	flags.StringVarP(&opts.Opts.Housekeeper.OperatorImageUrl, "operator-image-url", "", "", "Specify the URL of the container image for the housekeeper operator component")
	flags.BoolVarP(&opts.Opts.DeployHousekeeper, "deploy-housekeeper", "", false, "Deploy the Housekeeper Operator.")
	flags.StringVarP(&opts.Opts.NKD.BootstrapIgnHost, "bootstrap-ign-host", "", "", "Specify the host for Bootstrap Ignition")
	flags.StringVarP(&opts.Opts.NKD.BootstrapIgnPort, "bootstrap-ign-port", "", "", "Specify the port for Bootstrap Ignition")
}

func SetupDestroyCmdOpts(destroyCmd *cobra.Command) {
	flags := destroyCmd.Flags()
	flags.StringVarP(&opts.Opts.ClusterID, "cluster-id", "", "", "Cluster ID")
}

func SetupUpgradeCmdOpts(upgradeCmd *cobra.Command) {
	flags := upgradeCmd.Flags()
	flags.StringVarP(&opts.Opts.ClusterID, "cluster-id", "", "", "Cluster ID")
	flags.StringVarP(&opts.Opts.Housekeeper.KubeVersion, "kube-version", "", "", "Choose a specific kubernetes version for upgrading")
	flags.BoolVarP(&opts.Opts.Housekeeper.EvictPodForce, "force", "", false, "Force eviction of pods even if unsafe. This may result in data loss or service disruption, use with caution")
	flags.UintVarP(&opts.Opts.Housekeeper.MaxUnavailable, "maxunavailable", "", 0, "Number of nodes that are upgraded at the same time")
	flags.StringVarP(&opts.Opts.KubeConfigFile, "kubeconfig", "", "", "Specify the access path to the Kubeconfig file")
	flags.StringVarP(&opts.Opts.Housekeeper.OSImageURL, "imageurl", "", "", "The address of the container image to use for upgrading")
}

func SetupExtendCmdOpts(extendCmd *cobra.Command) {
	flags := extendCmd.Flags()
	flags.StringVarP(&opts.Opts.ClusterID, "cluster-id", "", "", "Cluster ID")
	flags.UintVarP(&opts.Opts.ExtendCount, "num", "n", 0, "The number of extend worker nodes")
}

func SetupTemplateCmdOpts(templateCmd *cobra.Command) {
	flags := templateCmd.Flags()
	flags.StringVarP(&opts.Opts.ClusterConfigFile, "file", "f", "", "Location of cluster deploy config file")
	flags.StringVarP(&opts.Opts.ClusterID, "cluster-id", "", "", "Cluster ID")
	flags.StringVar(&opts.Opts.Arch, "arch", "", "kubernetes Cluster deployment architecture")
	flags.StringVarP(&opts.Opts.Platform, "platform", "", "", "Select the infrastructure platform to deploy the cluster")
	flags.StringVarP(&opts.Opts.UserName, "username", "", "", "User to login the node")
	flags.StringVarP(&opts.Opts.Password, "password", "", "", "Password to login the node")
	flags.StringVarP(&opts.Opts.SSHKey, "sshkey", "", "", "Set nodes ssh private key for authentication")
	flags.StringArrayVarP(&opts.Opts.Master.Hostname, "master-hostname", "", []string{}, "Set master hostnames")
	flags.UintVar(&opts.Opts.Master.CPU, "master-cpu", 0, "Set master CPU (units: cores)")
	flags.UintVar(&opts.Opts.Master.RAM, "master-ram", 0, "Set master RAM (units: MB)")
	flags.UintVar(&opts.Opts.Master.Disk, "master-disk", 0, "Set master disk size (units: GB)")
	flags.StringArrayVarP(&opts.Opts.Master.IP, "master-ips", "", []string{}, "Set master IPs")
	flags.StringArrayVarP(&opts.Opts.Worker.Hostname, "worker-hostname", "", []string{}, "Set worker hostnames")
	flags.UintVar(&opts.Opts.Worker.CPU, "worker-cpu", 0, "Set worker CPU (units: cores)")
	flags.UintVar(&opts.Opts.Worker.RAM, "worker-ram", 0, "Set worker RAM (units: MB)")
	flags.UintVar(&opts.Opts.Worker.Disk, "worker-disk", 0, "Set worker disk size (units: GB)")
	flags.StringArrayVarP(&opts.Opts.Worker.IP, "worker-ips", "", []string{}, "Set worker IPs")
	flags.StringVarP(&opts.Opts.ImageRegistry, "image-registry", "", "", "Specify the registry address for pulling the Kubernetes component container image")
	flags.StringVarP(&opts.Opts.PauseImage, "pause-image", "", "", "Specify the image for the pause container")
	flags.StringVarP(&opts.Opts.ReleaseImageUrl, "release-image-url", "", "", "Specify the URL of the NestOS container image that contains the Kubernetes component. Only supports the qcow2 format.")
	flags.StringVarP(&opts.Opts.KubeVersion, "kubeversion", "", "", "Specify the version of Kubernetes to deploy")
	flags.StringVarP(&opts.Opts.Token, "token", "", "", "Specify the authentication token for accessing resources")
	flags.StringVarP(&opts.Opts.CertificateKey, "CertificateKey", "", "", "Specifies the certificate key to be added to the master node")
	flags.StringVarP(&opts.Opts.NetWork.ServiceSubnet, "service-subnet", "", "", "Specify the subnet for Kubernetes services")
	flags.StringVarP(&opts.Opts.NetWork.PodSubnet, "pod-subnet", "", "", "Specify the subnet for Kubernetes Pods.")
	flags.StringVarP(&opts.Opts.NetWork.Plugin, "network-plugin-url", "", "", "URL for the network plugin")
	flags.StringVarP(&opts.Opts.Housekeeper.ControllerImageUrl, "controller-image-url", "", "", "Specify the URL of the container image for the housekeeper controller component")
	flags.StringVarP(&opts.Opts.Housekeeper.OperatorImageUrl, "operator-image-url", "", "", "Specify the URL of the container image for the housekeeper operator component")
	flags.BoolVarP(&opts.Opts.DeployHousekeeper, "deploy-housekeeper", "", false, "Deploy the Housekeeper Operator.")
}
