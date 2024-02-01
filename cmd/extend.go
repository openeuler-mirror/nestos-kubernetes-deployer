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
	"nestos-kubernetes-deployer/cmd/command"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/infra"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewExtendCommand() *cobra.Command {
	extendCmd := &cobra.Command{
		Use:   "extend",
		Short: "Extend worker nodes of kubernetes cluster",
		RunE:  runExtendCmd,
	}
	command.SetupExtendCmdOpts(extendCmd)

	return extendCmd
}

func runExtendCmd(cmd *cobra.Command, args []string) error {
	clusterID, err := cmd.Flags().GetString("cluster-id")
	if err != nil {
		logrus.Errorf("Failed to get cluster-id: %v", err)
		return err
	}
	if clusterID == "" {
		logrus.Errorf("cluster-id is not provided: %v", err)
	}

	num, err := cmd.Flags().GetUint("num")
	if err != nil {
		logrus.Errorf("Failed to get assets directory: %v", err)
		return err
	}

	if err := configmanager.Initial(&opts.Opts); err != nil {
		logrus.Errorf("Failed to initialize configuration parameters: %v", err)
		return err
	}

	clusterConfig, err := configmanager.GetClusterConfig(clusterID)
	if err != nil {
		logrus.Errorf("Failed to get cluster config using the cluster id: %v", err)
		return err
	}
	extendArray(clusterConfig, int(num))

	if err := extendCluster(clusterConfig); err != nil {
		logrus.Errorf("Failed to extend %s cluster: %v", clusterID, err)
		return err
	}
	if err := configmanager.Persist(); err != nil {
		logrus.Errorf("Failed to persist the cluster asset: %v", err)
		return err
	}
	logrus.Infof("To access 'cluster-id:%s' cluster using 'kubectl', run 'export KUBECONFIG=%s'", clusterID, clusterConfig.AdminKubeConfig)

	return nil
}

func extendArray(c *asset.ClusterAsset, count int) {
	num := len(c.Worker)
	for i := 0; i < count; i++ {
		hostname := fmt.Sprintf("k8s-worker%02d", num+i+1)
		c.Worker = append(c.Worker, asset.NodeAsset{
			Hostname: hostname,
			IP:       "",
			HardwareInfo: asset.HardwareInfo{
				CPU:  c.Worker[i].CPU,
				RAM:  c.Worker[i].RAM,
				Disk: c.Worker[i].Disk,
			},
			Ignitions: c.Worker[i].Ignitions,
		})
	}
}

func extendCluster(conf *asset.ClusterAsset) error {
	// regenerate worker.tf
	var worker infra.Infra
	if err := worker.Generate(conf, "worker"); err != nil {
		logrus.Errorf("Failed to generate worker terraform file")
		return err
	}

	persistDir := configmanager.GetPersistDir()
	workerInfra := infra.InstanceCluster(persistDir, conf.Cluster_ID, "worker", uint(len(conf.Worker)))
	if err := workerInfra.Deploy(); err != nil {
		logrus.Errorf("Failed to deploy worker nodes:%v", err)
		return err
	}

	return nil
}
