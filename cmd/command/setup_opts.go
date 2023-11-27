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
	flags.StringVarP(&opts.Opts.DeployConfig, "file", "f", "./deploy-config.yaml", "location of cluster deploy config file, default ./deploy-config.yaml")
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
	flags.StringVarP(&opts.Opts.Upgrade.KubeVersion, "kube-version", "", "", "Choose a specific kubernetes version for upgrading")
	flags.BoolVarP(&opts.Opts.Upgrade.EvictPodForce, "force", "f", false, "Force evict pod")
	flags.IntVarP(&opts.Opts.Upgrade.MaxUnavailable, "maxunavailable", "n", 2, "Number of nodes that are upgraded at the same time")
	flags.StringVarP(&opts.Opts.Upgrade.KubeConfigFile, "kubeconfig", "", "./auth/config", "kubeconfig file access path")
	flags.StringVarP(&opts.Opts.Upgrade.OSImageURL, "imageurl", "", "", "The address of the container image to use for upgrading")
}
