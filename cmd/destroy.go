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
	"nestos-kubernetes-deployer/pkg/configmanager/manager"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var destroyClusterOpts struct {
	clusterID string
}

func NewDestroyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy a kubernetes cluster",
		RunE:  runDestroyCmd,
	}
	cmd.PersistentFlags().StringVarP(&destroyClusterOpts.clusterID, "cluster-id", "", "", "cluster ID")

	// cmd.AddCommand(destroy.NewDestroyMasterCommand())
	// cmd.AddCommand(destroy.NewDestroyWorkerCommand())

	return cmd
}

func runDestroyCmd(cmd *cobra.Command, args []string) error {
	if err := manager.Initial(cmd); err != nil {
		logrus.Errorf("Failed to initialize configuration parameters: %v", err)
		return err
	}
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	conf, err := manager.GetClusterConfig(clusterID)
	if err != nil {
		return err
	}

	if err := destroyCluster(conf); err != nil {
		return err
	}
	return nil
}

func destroyCluster( /*输入：TF配置*/ ) error {
	/*调用TF模块接口，销毁集群*/
	return nil
}
