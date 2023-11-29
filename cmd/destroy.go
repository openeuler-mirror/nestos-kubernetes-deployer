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
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewDestroyCommand() *cobra.Command {
	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy a kubernetes cluster",
		RunE:  runDestroyCmd,
	}
	command.SetupDestroyCmdOpts(destroyCmd)

	return destroyCmd
}

func runDestroyCmd(cmd *cobra.Command, args []string) error {
	clusterId, err := cmd.Flags().GetString("cluster-id")
	if err != nil {
		logrus.Errorf("Failed to get cluster-id: %v", err)
		return err
	}

	if err := configmanager.Initial(&opts.Opts); err != nil {
		logrus.Errorf("Failed to initialize configuration parameters: %v", err)
		return err
	}
	config, err := configmanager.GetClusterConfig(clusterId)
	if err != nil {
		logrus.Errorf("Failed to get cluster config using the cluster id: %v", err)
		return err
	}

	if err := destroyCluster(config); err != nil {
		return err
	}
	return nil
}

func destroyCluster(config *asset.ClusterAsset) error {
	/*调用TF模块接口，销毁集群*/
	return nil
}
