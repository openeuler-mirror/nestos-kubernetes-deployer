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
	flags.StringVarP(&opts.Opts.ClusterConfigFile, "file", "f", "./deploy-config.yaml", "location of cluster deploy config file, default ./deploy-config.yaml")
	flags.StringVarP(&opts.Opts.ClusterID, "cluster-id", "", "", "cluster id")
	flags.StringVarP(&opts.Opts.SSHKey, "sshkey", "", "", "path to SSH private keys that should be used for authentication.")
	flags.StringVarP(&opts.Opts.Platform, "platform", "", "", "select the infrastructure platform to deploy the cluster")
}

func SetupDestroyCmdOpts(destroyCmd *cobra.Command) {
	flags := destroyCmd.Flags()
	flags.StringVarP(&opts.Opts.ClusterID, "cluster-id", "", "", "cluster id")
}

func SetupUpgradeCmdOpts(upgradeCmd *cobra.Command) {
	flags := upgradeCmd.Flags()
	flags.StringVarP(&opts.Opts.ClusterID, "cluster-id", "", "", "cluster id")
	flags.StringVarP(&opts.Opts.Housekeeper.KubeVersion, "kube-version", "", "", "Choose a specific kubernetes version for upgrading")
	flags.BoolVarP(&opts.Opts.Housekeeper.EvictPodForce, "force", "", false, "Force eviction of pods even if unsafe. This may result in data loss or service disruption, use with caution")
	flags.IntVarP(&opts.Opts.Housekeeper.MaxUnavailable, "maxunavailable", "", 2, "Number of nodes that are upgraded at the same time")
	flags.StringVarP(&opts.Opts.KubeConfigFile, "kubeconfig", "", "/etc/nkd/pki/kubeconfig/admin.conf", "kubeconfig file access path")
	flags.StringVarP(&opts.Opts.Housekeeper.OSImageURL, "imageurl", "", "", "The address of the container image to use for upgrading")
}

func SetupExtendCmdOpts(extendCmd *cobra.Command) {
	flags := extendCmd.Flags()
	flags.StringVarP(&opts.Opts.ClusterID, "cluster-id", "", "", "cluster id")
	flags.IntVarP(&opts.Opts.Worker.Count, "num", "n", 0, "The number of expanded worker nodes")
}
