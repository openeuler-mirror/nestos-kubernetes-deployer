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
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	corev1 "k8s.io/api/core/v1"
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

	httpService := httpserver.NewHTTPService(configmanager.GetBootstrapIgnPort())
	defer httpService.Stop()

	if strings.ToLower(clusterConfig.Platform) == "pxe" || strings.ToLower(clusterConfig.Platform) == "ipxe" {
		if err := extendCluster(clusterConfig, httpService); err != nil {
			logrus.Errorf("Failed to extend %s cluster: %v", clusterID, err)
			return err
		}
		return nil
	}

	num, err := cmd.Flags().GetUint("num")
	if err != nil {
		logrus.Errorf("Failed to get the number of extended nodes: %v", err)
		return err
	}

	newHostnames := extendArray(clusterConfig, int(num))

	if err := extendCluster(clusterConfig, httpService); err != nil {
		logrus.Errorf("Failed to extend %s cluster: %v", clusterID, err)
		return err
	}

	logrus.Infof("Waiting for cluster extend nodes to be ready...")
	if err := checkNodesReady(clusterConfig, newHostnames); err != nil {
		return err
	}

	logrus.Infof("The cluster id:%s node is extended successfully", clusterID)

	return nil
}

func extendArray(c *asset.ClusterAsset, count int) []string {
	num := len(c.Worker)
	var newHostnames []string
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
		})
		newHostnames = append(newHostnames, hostname)
	}
	return newHostnames
}

func extendCluster(conf *asset.ClusterAsset, httpService *httpserver.HTTPService) error {
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
			httpService.AddFileToCache(constants.WorkerKS, data)
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
		go func() {
			if err := httpService.Start(); err != nil {
				logrus.Errorf("error starting http service: %v", err)
				return
			}
		}()

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
		go func() {
			if err := httpService.Start(); err != nil {
				logrus.Errorf("error starting http service: %v", err)
				return
			}
		}()

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
		pxe := &infra.PXE{
			IP:             pxeConfig.IP,
			HTTPServerPort: pxeConfig.HTTPServerPort,
			HTTPRootDir:    pxeConfig.HTTPRootDir,
			TFTPServerPort: pxeConfig.TFTPServerPort,
			TFTPRootDir:    pxeConfig.TFTPRootDir,
			HTTPService:    httpService,
		}
		p.SetInfra(pxe)
		if err := p.Extend(); err != nil {
			logrus.Errorf("Failed to extend worker nodes:%v", err)
			return err
		}
	case "ipxe":
		ipxeConfig := conf.InfraPlatform.(*infraasset.IPXEAsset)
		ipxe := &infra.IPXE{
			Port:              ipxeConfig.Port,
			FilePath:          ipxeConfig.FilePath,
			OSInstallTreePath: ipxeConfig.OSInstallTreePath,
			HTTPService:       httpService,
		}
		p.SetInfra(ipxe)
		if err := p.Extend(); err != nil {
			logrus.Errorf("Failed to extend worker nodes:%v", err)
			return err
		}
	default:
		return errors.New("unsupported platform")
	}

	return nil
}

// waitUntilNodesReady waits until all nodes are ready within a given timeout
func waitUntilNodesReady(ctx context.Context, clientset *kubernetes.Clientset, nodeNames []string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		time.Sleep(10 * time.Second)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			allNodesReady := true
			for _, nodeName := range nodeNames {
				node, err := clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
				if err != nil {
					allNodesReady = false
					break
				}
				for _, condition := range node.Status.Conditions {
					if condition.Type == corev1.NodeReady && condition.Status != corev1.ConditionTrue {
						allNodesReady = false
						break
					}
				}
			}
			if allNodesReady {
				return nil
			}
		}
	}
}

// checkNodesReady waits for all nodes to be ready
func checkNodesReady(conf *asset.ClusterAsset, nodeNames []string) error {
	clientset, err := kubeclient.CreateClient(conf.Kubernetes.AdminKubeConfig)
	if err != nil {
		logrus.Errorf("error creating Kubernetes client: %v", err)
		return err
	}

	// Wait for nodes to be ready
	timeout := 30 * time.Minute
	err = waitUntilNodesReady(context.Background(), clientset, nodeNames, timeout)
	if err != nil {
		logrus.Errorf("error waiting for nodes to be ready: %v", err)
		return err
	}

	return nil
}
