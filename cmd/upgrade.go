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
	"errors"
	"fmt"
	"nestos-kubernetes-deployer/cmd/command"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/kubeclient"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

func getFlagString(cmd *cobra.Command, flagName string) string {
	flagValue, err := cmd.Flags().GetString(flagName)
	if err != nil {
		logrus.Errorf("Failed to get %s parameter: %v", flagName, err)
	}
	return flagValue
}

func runUpgradeCmd(cmd *cobra.Command, args []string) error {
	clusterId := getFlagString(cmd, "cluster-id")
	kubeVersion := getFlagString(cmd, "kube-version")
	imageURL := getFlagString(cmd, "imageurl")
	if clusterId == "" {
		return errors.New("cluster-id is required")
	}

	if kubeVersion == "" {
		return errors.New("kube-version is required")
	}

	if imageURL == "" {
		return errors.New("imageurl is required")
	}

	if err := configmanager.Initial(&opts.Opts); err != nil {
		logrus.Errorf("Failed to initialize configuration parameters: %v", err)
		return err
	}
	clusterConfig, err := configmanager.GetClusterConfig(clusterId)
	if err != nil {
		logrus.Errorf("Failed to get cluster config using the cluster id: %v", err)
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

	adminconfig := filepath.Join(configmanager.GetPersistDir(), clusterConfig.ClusterID, "admin.config")
	if err := kubeclient.ApplyHousekeeperCR(yamlData, adminconfig); err != nil {
		logrus.Errorf("Failed to deploy Custom Resource: %v", err)
		return err
	}

	logrus.Info("Custom Resource deployed successfully.")
	return nil
}
