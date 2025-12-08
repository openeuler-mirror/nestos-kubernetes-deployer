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
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"nestos-kubernetes-deployer/cmd/command"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/kubeclient"
)

func NewUpgradeCommand() *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade your cluster to a newer version",
		Long:  "",
		RunE:  runUpgradeCmd,
	}
	command.SetupUpgradeCmdOpts(upgradeCmd)

	return upgradeCmd
}

func runUpgradeCmd(cmd *cobra.Command, args []string) error {
	if err := configmanager.Initial(&opts.Opts); err != nil {
		logrus.Debugf("Failed to initialize configuration parameters: %v", err)
		return err
	}
	clusterConfig, err := configmanager.GetClusterConfig(opts.Opts.ClusterID)
	if err != nil {
		logrus.Debugf("Failed to get cluster config using the cluster id: %v", err)
		return err
	}

	if err := upgradeCluster(clusterConfig); err != nil {
		return err
	}
	return nil
}

func upgradeCluster(clusterConfig *asset.ClusterAsset) error {

	// Define the YAML data for the Custom Resource (CR)
	yamlData := fmt.Sprintf(`
apiVersion: housekeeper.io/v1alpha1
kind: Update
metadata:
  name: housekeeper-upgrade
  namespace: housekeeper-system
spec:
  osImageURL: %s
  kubeVersion: %s
  evictPodForce: %t
  maxUnavailable: %d
`, clusterConfig.Housekeeper.OSImageURL, clusterConfig.Housekeeper.KubeVersion, clusterConfig.Housekeeper.EvictPodForce, clusterConfig.Housekeeper.MaxUnavailable)

	if err := kubeclient.ApplyYAML([]byte(yamlData)); err != nil {
		logrus.Debugf("Failed to deploy Custom Resource: %v", err)
		return err
	}

	logrus.Debugf("Custom Resource deployed successfully.")
	return nil
}
