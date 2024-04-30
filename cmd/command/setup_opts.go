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
	flags.StringVarP(&opts.Opts.ClusterConfigFile, "file", "f", "", "Location of the cluster deploy config file")
	flags.StringVarP(&opts.Opts.ClusterID, "cluster-id", "", "", "Unique identifier for the cluster")
	flags.StringVar(&opts.Opts.Arch, "arch", "", "Architecture for Kubernetes cluster deployment (e.g., amd64 or arm64)")
	flags.StringVarP(&opts.Opts.Platform, "platform", "", "", "Infrastructure platform for deploying the cluster (supports 'libvirt' or 'openstack')")
	flags.StringVarP(&opts.Opts.OSImage.Type, "os-type", "", "", "Operating system type for Kubernetes cluster deployment (e.g., nestos or openeuler)")
	flags.StringVarP(&opts.Opts.UserName, "username", "", "", "User name for node login")
	flags.StringVarP(&opts.Opts.Password, "password", "", "", "Password for node login")
	flags.StringVarP(&opts.Opts.SSHKey, "sshkey", "", "", "SSH key file path used for node authentication (default: ~/.ssh/id_rsa.pub)")
	flags.StringArrayVarP(&opts.Opts.Master.Hostname, "master-hostname", "", []string{}, "Hostnames of master nodes (e.g., --master-hostname [master-01] --master-hostname [master-02] ...)")
	flags.UintVar(&opts.Opts.Master.CPU, "master-cpu", 0, "CPU allocation for master nodes (units: cores)")
	flags.UintVar(&opts.Opts.Master.RAM, "master-ram", 0, "RAM allocation for master nodes (units: MB)")
	flags.UintVar(&opts.Opts.Master.Disk, "master-disk", 0, "Disk size allocation for master nodes (units: GB)")
	flags.StringArrayVarP(&opts.Opts.Master.IP, "master-ips", "", []string{}, "IP addresses of master nodes (e.g., --master-ips [master-ip-01] --master-ips [master-ip-02] ...)")
	flags.StringArrayVarP(&opts.Opts.Worker.Hostname, "worker-hostname", "", []string{}, "Hostnames of worker nodes (e.g., --worker-hostname [worker-01] --worker-hostname [worker-02] ...)")
	flags.UintVar(&opts.Opts.Worker.CPU, "worker-cpu", 0, "CPU allocation for worker nodes (units: cores)")
	flags.UintVar(&opts.Opts.Worker.RAM, "worker-ram", 0, "RAM allocation for worker nodes (units: MB)")
	flags.UintVar(&opts.Opts.Worker.Disk, "worker-disk", 0, "Disk size allocation for worker nodes (units: GB)")
	flags.StringArrayVarP(&opts.Opts.Worker.IP, "worker-ips", "", []string{}, "IP addresses of worker nodes (e.g., --worker-ips [worker-ip-01] --worker-ips [worker-ip-02] ...)")
	flags.StringVarP(&opts.Opts.Runtime, "runtime", "", "", "Container runtime type (docker, isulad or crio)")
	flags.StringVarP(&opts.Opts.ImageRegistry, "image-registry", "", "", "Registry address for Kubernetes component container images")
	flags.StringVarP(&opts.Opts.PauseImage, "pause-image", "", "", "Image for the pause container (e.g., pause:TAG)")
	flags.StringVarP(&opts.Opts.ReleaseImageUrl, "release-image-url", "", "", "URL of the NestOS container image containing Kubernetes component")
	flags.StringVarP(&opts.Opts.KubeVersion, "kubeversion", "", "", "Version of Kubernetes to deploy")
	flags.UintVarP(&opts.Opts.KubernetesAPIVersion, "kubernetes-apiversion", "", 0,
		"Sets the Kubernetes API version. Acceptable reference values:\n"+
			"  - 1 for Kubernetes versions < v1.15.0,\n"+
			"  - 2 for Kubernetes versions >= v1.15.0 && < v1.22.0,\n"+
			"  - 3 for Kubernetes versions >= v1.22.0")
	flags.StringVarP(&opts.Opts.Token, "token", "", "", "Used to validate the cluster information obtained from the control plane, with non-control plane nodes used for joining the cluster")
	flags.StringVarP(&opts.Opts.CertificateKey, "certificateKey", "", "", "The key that is used for decryption of certificates after they are downloaded from the secret upon joining a new master node.(the certificate key is a hex encoded string that is an AES key of size 32 bytes)")
	flags.StringVarP(&opts.Opts.NetWork.ServiceSubnet, "service-subnet", "", "", "Subnet used by Kubernetes services. (default: 10.96.0.0/16)")
	flags.StringVarP(&opts.Opts.NetWork.PodSubnet, "pod-subnet", "", "", "Subnet used for Kubernetes Pods. (default: 10.244.0.0/16)")
	flags.StringVarP(&opts.Opts.NetWork.Plugin, "network-plugin-url", "", "", "The deployment yaml URL of the network plugin")
	flags.StringVarP(&opts.Opts.Housekeeper.ControllerImageUrl, "controller-image-url", "", "", "URL of the container image for the housekeeper controller component")
	flags.StringVarP(&opts.Opts.Housekeeper.OperatorImageUrl, "operator-image-url", "", "", "URL of the container image for the housekeeper operator component")
	flags.BoolVarP(&opts.Opts.DeployHousekeeper, "deploy-housekeeper", "", false, "Deploy the Housekeeper Operator. (default: false)")
	flags.StringVarP(&opts.Opts.NKD.BootstrapIgnHost, "bootstrap-ign-host", "", "", "Ignition service address (domain name or IP)")
	flags.StringVarP(&opts.Opts.NKD.BootstrapIgnPort, "bootstrap-ign-port", "", "", "Ignition service port (default: 9080)")
	flags.StringVarP(&opts.Opts.PreHookScript, "prehook-script", "", "", "Specify a script file or directory to execute before cluster deployment as hooks")
	flags.StringVarP(&opts.Opts.PostHookYaml, "posthook-yaml", "", "", "Specify a YAML file or directory to apply after cluster deployment using 'kubectl apply'")
}

func SetupDestroyCmdOpts(destroyCmd *cobra.Command) {
	flags := destroyCmd.Flags()
	flags.StringVarP(&opts.Opts.ClusterID, "cluster-id", "", "", "Unique identifier for the cluster")
}

func SetupUpgradeCmdOpts(upgradeCmd *cobra.Command) {
	flags := upgradeCmd.Flags()
	flags.StringVarP(&opts.Opts.ClusterID, "cluster-id", "", "", "Unique identifier for the cluster")
	flags.StringVarP(&opts.Opts.Housekeeper.KubeVersion, "kube-version", "", "", "Choose a specific kubernetes version for upgrading")
	flags.BoolVarP(&opts.Opts.Housekeeper.EvictPodForce, "force", "", false, "Force eviction of pods even if unsafe. This may result in data loss or service disruption, use with caution (default: false)")
	flags.UintVarP(&opts.Opts.Housekeeper.MaxUnavailable, "maxunavailable", "", 0, "Number of nodes that are upgraded at the same time (default: 2)")
	flags.StringVarP(&opts.Opts.KubeConfigFile, "kubeconfig", "", "", "Specify the access path to the Kubeconfig file")
	flags.StringVarP(&opts.Opts.Housekeeper.OSImageURL, "imageurl", "", "", "The address of the container image to use for upgrading")
}

func SetupExtendCmdOpts(extendCmd *cobra.Command) {
	flags := extendCmd.Flags()
	flags.StringVarP(&opts.Opts.ClusterID, "cluster-id", "", "", "Unique identifier for the cluster")
	flags.UintVarP(&opts.Opts.ExtendCount, "num", "n", 0, "The number of extend worker nodes")
}

func SetupTemplateCmdOpts(templateCmd *cobra.Command) {
	flags := templateCmd.Flags()
	flags.StringVarP(&opts.Opts.ClusterConfigFile, "output", "o", "", "Generates a default configuration template at the specified location")
}
