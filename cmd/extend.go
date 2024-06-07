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
	"context"
	"errors"
	"fmt"
	"nestos-kubernetes-deployer/cmd/command"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/infraasset"
	"nestos-kubernetes-deployer/pkg/constants"
	"nestos-kubernetes-deployer/pkg/httpserver"
	"nestos-kubernetes-deployer/pkg/infra"
	"nestos-kubernetes-deployer/pkg/kubeclient"
	"nestos-kubernetes-deployer/pkg/osmanager"
	"nestos-kubernetes-deployer/pkg/terraform"
	"nestos-kubernetes-deployer/pkg/tftpserver"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
	cleanup := command.SetuploggerHook(opts.Opts.RootOptDir)
	defer cleanup()

	clusterID, err := cmd.Flags().GetString("cluster-id")
	if err != nil {
		logrus.Errorf("Failed to get cluster-id: %v", err)
		return err
	}
	if clusterID == "" {
		logrus.Errorf("cluster-id is not provided: %v", err)
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

	num, err := cmd.Flags().GetUint("num")
	if err != nil {
		platform := strings.ToLower(clusterConfig.Platform)
		if platform != "pxe" && platform != "ipxe" {
			logrus.Errorf("Failed to get the number of extended nodes: %v", err)
			return err
		}
	}

	if err := extendCluster(clusterConfig, num); err != nil {
		logrus.Errorf("Failed to extend %s cluster: %v", clusterID, err)
		return err
	}

	logrus.Infof("The cluster nodes are extended successfully")

	return nil
}

func extendCluster(conf *asset.ClusterAsset, num uint) error {
	httpService := httpserver.NewHTTPService(configmanager.GetBootstrapIgnPort())
	defer httpService.Stop()

	data, err := os.ReadFile(conf.BootConfig.Worker.Path)
	if err != nil {
		logrus.Errorf("error reading boot config file: %v", err)
		return err
	}

	osMgr := osmanager.NewOSManager(conf)
	if osMgr.IsNestOS() {
		httpService.AddFileToCache(constants.WorkerIgn, data)
	}
	if osMgr.IsGeneralOS() {
		if strings.ToLower(conf.Platform) == "pxe" || strings.ToLower(conf.Platform) == "ipxe" {
			httpService.AddFileToCache(constants.Worker+constants.KickstartSuffix, data)
		}
	}

	httpService.AddFileToCache(constants.WorkerIgn, data)

	if err := configmanager.Persist(); err != nil {
		logrus.Errorf("Failed to persist the cluster asset: %v", err)
		return err
	}

	p := infra.InfraPlatform{}
	switch strings.ToLower(conf.Platform) {
	case "libvirt":
		httpserver.StartHTTPService(httpService)

		// regenerate worker.tf
		var worker terraform.Infra
		if err := worker.Generate(conf, "worker"); err != nil {
			logrus.Errorf("Failed to generate worker terraform file")
			return err
		}

		libvirtWorker := &infra.Libvirt{
			PersistDir: configmanager.GetPersistDir(),
			ClusterID:  conf.ClusterID,
			Node:       "worker",
			Count:      uint(len(conf.Worker)),
		}
		p.SetInfra(libvirtWorker)
		if err := p.Extend(); err != nil {
			logrus.Errorf("Failed to extend worker nodes:%v", err)
			return err
		}
	case "openstack":
		httpserver.StartHTTPService(httpService)

		// regenerate worker.tf
		var worker terraform.Infra
		if err := worker.Generate(conf, "worker"); err != nil {
			logrus.Errorf("Failed to generate worker terraform file")
			return err
		}

		openstackWorker := &infra.OpenStack{
			PersistDir: configmanager.GetPersistDir(),
			ClusterID:  conf.ClusterID,
			Node:       "worker",
			Count:      uint(len(conf.Worker)),
		}
		p.SetInfra(openstackWorker)
		if err := p.Extend(); err != nil {
			logrus.Errorf("Failed to extend worker nodes:%v", err)
			return err
		}
	case "pxe":
		pxeConfig := conf.InfraPlatform.(*infraasset.PXEAsset)
		httpService.Port = pxeConfig.HTTPServerPort
		httpService.DirPath = pxeConfig.HTTPRootDir
		httpserver.StartHTTPService(httpService)

		tftpService := tftpserver.NewTFTPService(pxeConfig.IP, pxeConfig.TFTPServerPort, pxeConfig.TFTPRootDir)
		go func() {
			select {
			case <-httpService.Ch:
				logrus.Info("tftp server stop")
				tftpService.Stop()
				return
			}
		}()
		go func() {
			if err := tftpService.Start(); err != nil {
				logrus.Errorf("error starting http service: %v", err)
				return
			}
		}()
		defer tftpService.Stop()
	case "ipxe":
		ipxeConfig := conf.InfraPlatform.(*infraasset.IPXEAsset)
		httpService.Port = ipxeConfig.Port
		httpService.DirPath = ipxeConfig.OSInstallTreePath
		fileContent, err := os.ReadFile(ipxeConfig.FilePath)
		if err != nil {
			return err
		}
		httpService.AddFileToCache(constants.IPXECfg, fileContent)
		httpserver.StartHTTPService(httpService)

	default:
		return errors.New("unsupported platform")
	}

	if err := checkNodesReady(context.Background(), conf, int(num)); err != nil {
		return err
	}

	return nil
}

// checkNodesReady waits for all nodes to be ready
func checkNodesReady(ctx context.Context, conf *asset.ClusterAsset, num int) error {
	clientset, err := kubeclient.CreateClient(conf.Kubernetes.AdminKubeConfig)
	if err != nil {
		logrus.Errorf("error creating Kubernetes client: %v", err)
		return err
	}

	// Get the current number of ready nodes
	readyNodesCount, err := getReadyNodesCount(ctx, clientset)
	if err != nil {
		logrus.Errorf("error getting current ready nodes count: %v", err)
		return err
	}
	allNodeNums := readyNodesCount + num

	// Wait for nodes to be ready
	timeout := 30 * time.Minute
	err = waitForMinimumReadyNodes(ctx, clientset, allNodeNums, timeout)
	if err != nil {
		logrus.Errorf("error waiting for nodes to be ready: %v", err)
		return err
	}

	return nil
}

// getReadyNodesCount returns the number of ready nodes in the cluster
func getReadyNodesCount(ctx context.Context, clientset *kubernetes.Clientset) (int, error) {
	nodeList, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to list nodes: %v", err)
	}

	readyNodesCount := 0
	for _, node := range nodeList.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == v1.NodeReady && condition.Status == v1.ConditionTrue {
				readyNodesCount++
				break
			}
		}
	}

	return readyNodesCount, nil
}

func waitForMinimumReadyNodes(ctx context.Context, clientset *kubernetes.Clientset, requiredReadyNodes int, timeout time.Duration) error {
	logrus.Infof("Waiting for cluster extend nodes to be ready...")
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		time.Sleep(10 * time.Second)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			readyNodeCount, err := getReadyNodesCount(ctx, clientset)
			if err != nil {
				return err
			}

			if readyNodeCount >= requiredReadyNodes {
				return nil
			}
		}
	}
}
