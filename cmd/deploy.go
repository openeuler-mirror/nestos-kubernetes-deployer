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
	"io"
	"nestos-kubernetes-deployer/cmd/command"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/data"
	"nestos-kubernetes-deployer/pkg/cert"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/infraasset"
	"nestos-kubernetes-deployer/pkg/constants"
	"nestos-kubernetes-deployer/pkg/httpserver"
	"nestos-kubernetes-deployer/pkg/infra"
	"nestos-kubernetes-deployer/pkg/kubeclient"
	"nestos-kubernetes-deployer/pkg/osmanager"
	"nestos-kubernetes-deployer/pkg/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	wait "k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

func NewDeployCommand() *cobra.Command {
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a kubernetes cluster",
		RunE:  runDeployCmd,
	}
	command.SetupDeployCmdOpts(deployCmd)

	return deployCmd
}

const (
	clusterID         = "cluster"
	clusterConfigFile = "cluster_config.yaml"
	namespace         = "housekeeper-system"
)

func runDeployCmd(cmd *cobra.Command, args []string) error {
	cleanup := command.SetuploggerHook(opts.Opts.RootOptDir)
	defer cleanup()

	if err := validateDeployConfig(); err != nil {
		return err
	}

	// Initialize configuration parameters
	config, err := getClusterConfig(&opts.Opts)
	if err != nil {
		return err
	}

	if err := deployCluster(config); err != nil {
		logrus.Errorf("Failed to deploy %s cluster: %v", clusterID, err)
		return err
	}

	logrus.Infof("To access 'cluster-id:%s' cluster using 'kubectl', run 'export KUBECONFIG=%s'", clusterID, config.AdminKubeConfig)
	return nil
}

func validateDeployConfig() error {
	opts.Opts.ClusterID = clusterID
	clusterConfigFile := filepath.Join(opts.Opts.RootOptDir, opts.Opts.ClusterID, clusterConfigFile)
	// Check if clusterConfigFile already exists
	if _, err := os.Stat(clusterConfigFile); err == nil {
		logrus.Debugf("cluster ID: %s already exists", opts.Opts.ClusterID)
		return fmt.Errorf("cluster ID: %s already exists", opts.Opts.ClusterID)
	}

	// Check if kubectl is installed
	if !kubeclient.IsKubectlInstalled() {
		logrus.Debug("kubectl is not installed")
		return fmt.Errorf("kubectl is not installed")
	}
	return nil
}

func getClusterConfig(options *opts.OptionsList) (*asset.ClusterAsset, error) {
	if err := configmanager.Initial(options); err != nil {
		logrus.Errorf("Failed to initialize configuration parameters: %v", err)
		return nil, err
	}

	config, err := configmanager.GetClusterConfig(clusterID)
	if err != nil {
		logrus.Errorf("Failed to get cluster config using the cluster id: %v", err)
		return nil, err
	}
	return config, nil
}

func deployCluster(conf *asset.ClusterAsset) error {
	hs := httpserver.NewHTTPService(configmanager.GetBootstrapIgnPort())
	defer hs.Stop()

	if err := createCluster(conf, hs); err != nil {
		logrus.Errorf("Failed to create cluster: %v", err)
		return err
	}

	kubeClient, err := kubeclient.CreateClient(conf.Kubernetes.AdminKubeConfig)
	if err != nil {
		logrus.Errorf("Failed to create kubernetes client %v", err)
		return err
	}

	if err := waitForAPIReady(kubeClient); err != nil {
		logrus.Errorf("Failed while waiting for Kubernetes API to be ready: %v", err)
		return err
	}
	// Set kubeconfig environment variable
	os.Setenv("KUBECONFIG", conf.Kubernetes.AdminKubeConfig)

	// Apply network plugin
	if err := applyNetworkPlugin(conf.Network.Plugin, conf.IsNestOS); err != nil {
		logrus.Errorf("Failed to apply network plugin: %v", err)
		return err
	}
	logrus.Info("Network plugin deployment completed successfully.")

	if conf.Housekeeper.DeployHousekeeper {
		logrus.Info("Starting deployment of Housekeeper...")
		if err := deployHousekeeper(conf.Housekeeper, conf.Kubernetes.AdminKubeConfig); err != nil {
			logrus.Errorf("Failed to deploy operator: %v", err)
			return err
		}
		logrus.Info("Housekeeper deployment completed successfully.")
	}

	// Wait for pods to be ready
	if err := waitForPodsReady(kubeClient); err != nil {
		logrus.Errorf("Failed while waiting for pods to be in 'Ready' state: %v", err)
		return err
	}

	logrus.Info("Cluster deployment completed successfully!")
	return nil
}

func createCluster(conf *asset.ClusterAsset, httpService *httpserver.HTTPService) error {
	osMgr := osmanager.NewOSManager(conf)
	if err := osMgr.GenerateOSConfig(); err != nil {
		logrus.Errorf("Error generating OS config: %v", err)
		return err
	}

	if osMgr.IsNestOS() {
		if err := addIgnitionFiles(httpService, conf); err != nil {
			return err
		}
	}
	if osMgr.IsOpenEuler() && len(conf.Master) > 0 {
		certs, _ := cert.CertsToBytes(conf.Master[0].Certs)
		if err := httpService.AddFileToCache(constants.CertsFiles, certs); err != nil {
			return err
		}
	}

	if err := configmanager.Persist(); err != nil {
		logrus.Errorf("Failed to persist the cluster asset: %v", err)
		return err
	}

	p := infra.InfraPlatform{}
	switch strings.ToLower(conf.Platform) {
	case "libvirt":
		if err := httpService.Start(); err != nil {
			return fmt.Errorf("error starting http service: %v", err)
		}

		libvirtMaster := &infra.Libvirt{
			PersistDir: configmanager.GetPersistDir(),
			ClusterID:  conf.ClusterID,
			Node:       "master",
			Count:      uint(len(conf.Master)),
		}

		p.SetInfra(libvirtMaster)
		if err := p.Deploy(); err != nil {
			logrus.Errorf("Failed to deploy master nodes:%v", err)
			return err
		}

		libvirtWorker := &infra.Libvirt{
			PersistDir: configmanager.GetPersistDir(),
			ClusterID:  conf.ClusterID,
			Node:       "worker",
			Count:      uint(len(conf.Master)),
		}
		p.SetInfra(libvirtWorker)
		if err := p.Deploy(); err != nil {
			logrus.Errorf("Failed to deploy worker nodes:%v", err)
			return err
		}
	case "openstack":
		if err := httpService.Start(); err != nil {
			return fmt.Errorf("error starting http service: %v", err)
		}

		openstackMaster := &infra.OpenStack{
			PersistDir: configmanager.GetPersistDir(),
			ClusterID:  conf.ClusterID,
			Node:       "master",
			Count:      uint(len(conf.Master)),
		}
		p.SetInfra(openstackMaster)
		if err := p.Deploy(); err != nil {
			logrus.Errorf("Failed to deploy master nodes:%v", err)
			return err
		}

		openstackWorker := &infra.OpenStack{
			PersistDir: configmanager.GetPersistDir(),
			ClusterID:  conf.ClusterID,
			Node:       "worker",
			Count:      uint(len(conf.Master)),
		}
		p.SetInfra(openstackWorker)
		if err := p.Deploy(); err != nil {
			logrus.Errorf("Failed to deploy worker nodes:%v", err)
			return err
		}
	case "pxe":
		pxeConfig := conf.InfraPlatform.(*infraasset.PXEAsset)
		pxe := &infra.PXE{
			HTTPServerPort: pxeConfig.HTTPServerPort,
			HTTPRootDir:    pxeConfig.HTTPRootDir,
			TFTPServerIP:   pxeConfig.TFTPServerIP,
			TFTPServerPort: pxeConfig.TFTPServerPort,
			TFTPRootDir:    pxeConfig.TFTPRootDir,
			HTTPService:    httpService,
		}
		p.SetInfra(pxe)
		if err := p.Deploy(); err != nil {
			logrus.Errorf("Failed to deploy PXE:%v", err)
			return err
		}
	case "ipxe":
		ipxeConfig := conf.InfraPlatform.(*infraasset.IPXEAsset)
		ipxe := &infra.IPXE{
			IPXEPort:              ipxeConfig.IPXEPort,
			IPXEFilePath:          ipxeConfig.IPXEFilePath,
			IPXEOSInstallTreePath: ipxeConfig.IPXEOSInstallTreePath,
			HTTPService:           httpService,
		}
		p.SetInfra(ipxe)
		if err := p.Deploy(); err != nil {
			logrus.Errorf("Failed to deploy IPXE:%v", err)
			return err
		}
	default:
		return errors.New("unsupported platform")
	}

	return nil
}

func waitForAPIReady(client *kubernetes.Clientset) error {
	apiTimeout := 60 * time.Minute
	ctx := context.Background()
	apiContext, cancel := context.WithTimeout(ctx, apiTimeout)
	logrus.Infof("Waiting up to %v for the Kubernetes API ready...", apiTimeout)
	defer cancel()

	discovery := client.Discovery()
	wait.Until(func() {
		version, err := discovery.ServerVersion()
		if err == nil {
			logrus.Infof("The Kubernetes API %s up", version)
			cancel()
		} else {
			logrus.Debugf("Still waiting for Kubernetes API ready: %v", err)
		}
	}, 2*time.Second, apiContext.Done())

	err := apiContext.Err()
	if err != nil && err != context.Canceled {
		logrus.Errorf("Failed to waiting for kubernetes API: %v", err)
		return err
	}
	return nil
}

func waitForPodsReady(client *kubernetes.Clientset) error {
	waitDuration := 20 * time.Minute
	namespace := "kube-system"
	waitCtx, cancel := context.WithTimeout(context.Background(), waitDuration)
	defer cancel()
	logrus.Infof("Waiting up to %v for the Kubernetes Pods ready ...", waitDuration)

	err := wait.PollImmediate(10*time.Second, waitDuration, func() (bool, error) {
		pods, err := client.CoreV1().Pods(namespace).List(waitCtx, metav1.ListOptions{})
		if err != nil {
			logrus.Errorf("Failed to list Pods: %v", err)
			return false, nil
		}
		allReady := true
		for _, pod := range pods.Items {
			for _, condition := range pod.Status.Conditions {
				if condition.Type == corev1.PodReady && condition.Status != corev1.ConditionTrue {
					allReady = false
					logrus.Infof("Pod %s in namespace %s is not in Ready state", pod.Name, pod.Namespace)
					break
				}
			}
		}

		if allReady {
			logrus.Infof("All Pods in namespace %s are in Ready state", namespace)
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		logrus.Errorf("Failed to wait for Pods to be Ready: %v", err)
		return err
	}
	return nil
}

func deployHousekeeper(tmplData interface{}, kubeconfig string) error {
	dir, err := data.Assets.Open("housekeeper")
	if err != nil {
		return err
	}
	defer dir.Close()
	child, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, childInfo := range child {
		filePath := filepath.Join("housekeeper", childInfo.Name())
		data, err := utils.FetchAndUnmarshalUrl(filePath, tmplData)

		switch childInfo.Name() {
		case "1housekeeper.io_updates.yaml":
			err = kubeclient.DeployCRD(string(data), kubeconfig)
		case "2namespace.yaml":
			err = kubeclient.DeployNamespace(string(data), kubeconfig)
		case "3role.yaml":
			err = kubeclient.DeployClusterRole(string(data), kubeconfig)
		case "4role_binding.yaml":
			err = kubeclient.DeployClusterRoleBinding(string(data), kubeconfig)
		case "5deployment.yaml.template":
			err = kubeclient.DeployDeployment(string(data), kubeconfig, namespace)
		case "6daemonset.yaml.template":
			err = kubeclient.DeployDaemonSet(string(data), kubeconfig, namespace)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func applyNetworkPlugin(pluginConfigPath string, isNestOS bool) error {
	var content []byte
	var err error

	// Check if the pluginConfigPath is an HTTP(S) link or a local file path
	if strings.HasPrefix(pluginConfigPath, "http://") || strings.HasPrefix(pluginConfigPath, "https://") {
		response, err := http.Get(pluginConfigPath)
		if err != nil {
			logrus.Errorf("Failed to fetch network plugin configuration from URL: %v", err)
			return err
		}
		defer response.Body.Close()
		content, err = io.ReadAll(response.Body)
		if err != nil {
			logrus.Errorf("Failed to read content from HTTP response: %v", err)
			return err
		}
	} else {
		// Read the content from the local file
		content, err = os.ReadFile(pluginConfigPath)
		if err != nil {
			logrus.Errorf("Failed to read network plugin configuration file: %v", err)
			return err
		}
	}

	// 在类似NestOS 或者 Fedora CoreOS 这类不可变基础设施中，目录/usr为只读目录。在支持FlexVolume时，默认路径为
	// "/usr/libexec/kubernetes/kubelet-plugins"，而 FlexVolume 的目录必须是可写入的，
	// 该功能特性才能正常工作，为了解决这个问题将/usr目录修改为可写目录/opt.
	// Check if the content contains "/usr/libexec/kubernetes/kubelet-plugins"
	if isNestOS && strings.Contains(string(content), "/usr/libexec/kubernetes/kubelet-plugins") {
		content = []byte(strings.ReplaceAll(string(content),
			"/usr/libexec/kubernetes/kubelet-plugins",
			"/opt/libexec/kubernetes/kubelet-plugins"))
	}

	// Save the modified content to a file in the "/tmp" directory with a fixed name
	tmpFilePath := "/tmp/modified-plugin-config.yaml"

	err = os.WriteFile(tmpFilePath, content, 0644)
	if err != nil {
		logrus.Errorf("Failed to write content to file: %v", err)
		return err
	}

	// Apply the modified configuration using kubeclient
	if err := kubeclient.RunKubectlApplyWithYaml(tmpFilePath); err != nil {
		logrus.Errorf("Failed to apply network plugin configuration: %v", err)
		return err
	}

	// removal of the temporary file
	defer func() {
		if err := os.Remove(tmpFilePath); err != nil {
			logrus.Errorf("Failed to remove temporary file: %v", err)
		}
	}()

	return nil
}

func addIgnitionFiles(httpService *httpserver.HTTPService, conf *asset.ClusterAsset) error {
	// Ignition files are divided into three types:
	// control plane ignition files for initializing the cluster,
	// master ignition files for master node joining the cluster,
	// and worker ignition files for worker node joining the cluster.

	// Only one master node
	if err := httpService.AddFileToCache(constants.ControlplaneIgn, conf.BootConfig.Controlplane.Content); err != nil {
		return fmt.Errorf("error adding control plane ignition file to cache: %v", err)
	}

	// multiple master nodes
	if len(conf.Master) > 1 {
		if err := httpService.AddFileToCache(constants.MasterIgn, conf.BootConfig.Master.Content); err != nil {
			return fmt.Errorf("error adding master ignition file to cache: %v", err)
		}
	}

	if err := httpService.AddFileToCache(constants.WorkerIgn, conf.BootConfig.Worker.Content); err != nil {
		return fmt.Errorf("error adding worker ignition file to cache: %v", err)
	}

	return nil
}
