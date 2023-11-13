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
package cmd

import (
	"nestos-kubernetes-deployer/cmd/command"

	"github.com/spf13/cobra"
)

func NewDeployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a kubernetes cluster",
		RunE:  runDeployCmd,
	}

	cmd.PersistentFlags().StringVar(&command.ClusterOpts.ClusterId, "cluster-id", "", "clusterID of kubernetes cluster")
	cmd.PersistentFlags().StringVar(&command.ClusterOpts.GatherDeployOpts.SSHKey, "sshkey", "", "Path to SSH private keys that should be used for authentication.")
	cmd.PersistentFlags().StringVar(&command.ClusterOpts.Platform, "platform", "", "Select the infrastructure platform to deploy the cluster")

	// cmd.AddCommand(deploy.NewDeployMasterCommand())
	// cmd.AddCommand(deploy.NewDeployWorkerCommand())

	return cmd
}

func runDeployCmd(cmd *cobra.Command, args []string) error {

	return nil
}

// 生成部署集群所需配置数据
func runInstallconfig() error {

	return nil
}

func runDeployCluster() error {
	return nil
}

// 等待集群安装完成
func waitForClusterComplete(config string) error {

	return nil
}

// check 集群running状态
func checkPod() error {
	return nil
}
