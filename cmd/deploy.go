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
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"

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
	"nestos-kubernetes-deployer/pkg/tftpserver"
	"nestos-kubernetes-deployer/pkg/utils"
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
	kubeSystemNS = "kube-system"
)

func runDeployCmd(cmd *cobra.Command, args []string) error {
	cleanup := command.SetuploggerHook(opts.Opts.RootOptDir)
	defer cleanup()

	if err := validateDeployConfig(); err != nil {
		logrus.Debugf("Deploy configuration validation failed: %v", err)
		return err
	}

	// Initialize configuration parameters
	config, err := getClusterConfig(&opts.Opts)
	if err != nil {
		logrus.Debugf("Failed to get cluster configuration: %v", err)
		return err
	}

	if err := createCluster(config); err != nil {
		logrus.Debugf("Failed to create cluster: %v", err)
		return err
	}

	logrus.Info("Cluster deployed successfully!")
	logrus.Infof("To access 'cluster-id:%s' cluster using 'kubectl', run 'export KUBECONFIG=%s'", opts.Opts.ClusterID, config.AdminKubeConfig)
	return nil
}

func validateDeployConfig() error {
	//sync default kubeconfig path
	utils.SetDefaultKubeConfigPath(opts.Opts.RootOptDir, opts.Opts.ClusterID)

	// Check if clusterConfigFile already exists
	defClusterConfigFilePath := filepath.Join(opts.Opts.RootOptDir, opts.Opts.ClusterID, configmanager.GetClusterConfigFileName())
	if _, err := os.Stat(defClusterConfigFilePath); err == nil {
		logrus.Debugf("cluster ID: %s already exists", opts.Opts.ClusterID)
		return fmt.Errorf("cluster ID: %s already exists", opts.Opts.ClusterID)
	}

	return nil
}

func getClusterConfig(options *opts.OptionsList) (*asset.ClusterAsset, error) {
	if err := configmanager.Initial(options); err != nil {
		logrus.Debugf("Failed to initialize configuration parameters: %v", err)
		return nil, err
	}

	config, err := configmanager.GetClusterConfig(opts.Opts.ClusterID)
	if err != nil {
		logrus.Debugf("Failed to get cluster config using the cluster id: %v", err)
		return nil, err
	}
	return config, nil
}

func createInfraInstance(platform, nodeType string, count uint, conf *asset.ClusterAsset) infra.Infrastructure {
	switch strings.ToLower(platform) {
	case "libvirt":
		return &infra.Libvirt{
			PersistDir: configmanager.GetPersistDir(),
			ClusterID:  conf.ClusterID,
			Node:       nodeType,
			Count:      count,
		}
	case "openstack":
		return &infra.OpenStack{
			PersistDir: configmanager.GetPersistDir(),
			ClusterID:  conf.ClusterID,
			Node:       nodeType,
			Count:      count,
		}
	default:
		return nil
	}
}

func deployNodes(platform infra.Infrastructure, nodeType string) error {
	p := infra.InfraPlatform{}
	p.SetInfra(platform)
	if err := p.Deploy(); err != nil {
		logrus.Debugf("Failed to deploy %s nodes: %v", nodeType, err)
		return err
	}
	return nil
}

func createCluster(conf *asset.ClusterAsset) error {
	platform := strings.ToLower(conf.Platform)

	httpService := httpserver.NewHTTPService(configmanager.GetBootstrapIgnPort())
	defer httpService.Stop()

	osMgr := osmanager.NewOSManager(conf)
	if err := osMgr.GenerateOSConfig(); err != nil {
		logrus.Debugf("Failed to generate OS config: %v", err)
		return err
	}

	if osMgr.IsNestOS() {
		if err := addIgnitionFiles(httpService, conf); err != nil {
			return err
		}
	}
	if osMgr.IsGeneralOS() && len(conf.Master) > 0 {
		certs, _ := cert.CertsToBytes(conf.Master[0].Certs)
		if err := httpService.AddFileToCache(constants.CertsFiles, certs); err != nil {
			return err
		}

		if len(conf.Kubernetes.RpmPackagePath) > 0 {
			httpService.PackageDir = conf.Kubernetes.RpmPackagePath
		}

		if platform == "pxe" || platform == "ipxe" {
			if err := addKickstartFiles(httpService, conf); err != nil {
				return fmt.Errorf("error adding kickstart file to cache: %v", err)
			}
		}
	}

	if err := configmanager.Persist(); err != nil {
		logrus.Debugf("Failed to persist cluster asset: %v", err)
		return err
	}

	switch platform {
	case "libvirt", "openstack":
		httpserver.StartHTTPService(httpService)

		masterInfra := createInfraInstance(platform, "master", uint(len(conf.Master)), conf)
		if masterInfra == nil {
			return fmt.Errorf("unsupported platform: %s", platform)
		}
		if err := deployNodes(masterInfra, "master"); err != nil {
			logrus.Debugf("Failed to deploy master nodes: %v", err)
			return err
		}

		workerInfra := createInfraInstance(platform, "worker", uint(len(conf.Worker)), conf)
		if workerInfra == nil {
			return fmt.Errorf("unsupported platform: %s", platform)
		}
		if err := deployNodes(workerInfra, "worker"); err != nil {
			logrus.Debugf("Failed to deploy worker nodes: %v", err)
			return err
		}
	case "pxe":
		pxeConfig, ok := conf.InfraPlatform.(*infraasset.PXEAsset)
		if !ok || pxeConfig == nil {
			return fmt.Errorf("infra platform mismatch: expected *infraasset.PXEAsset, got %T", conf.InfraPlatform)
		}
		httpService.Port = pxeConfig.HTTPServerPort
		httpService.DirPath = pxeConfig.HTTPRootDir
		httpserver.StartHTTPService(httpService)

		tftpService := tftpserver.NewTFTPService(pxeConfig.IP, pxeConfig.TFTPServerPort, pxeConfig.TFTPRootDir)
		go func() {
			if err := tftpService.Start(); err != nil {
				logrus.Debugf("error starting http service: %v", err)
				return
			}
			<-httpService.Ch
			logrus.Debug("tftp server stop")
			tftpService.Stop()
		}()

	case "ipxe":
		ipxeConfig, ok := conf.InfraPlatform.(*infraasset.IPXEAsset)
		if !ok || ipxeConfig == nil {
			return fmt.Errorf("infra platform mismatch: expected *infraasset.IPXEAsset, got %T", conf.InfraPlatform)
		}
		httpService.Port = ipxeConfig.Port
		httpService.DirPath = ipxeConfig.OSInstallTreePath
		fileContent, err := os.ReadFile(ipxeConfig.FilePath)
		if err != nil {
			return err
		}
		if err := httpService.AddFileToCache(constants.IPXECfg, fileContent); err != nil {
			return fmt.Errorf("error adding ipxe config file to cache: %v", err)
		}
		httpserver.StartHTTPService(httpService)
	default:
		return errors.New("unsupported platform")
	}

	if err := clusterCreatePost(conf); err != nil {
		return err
	}
	return nil
}

func clusterCreatePost(conf *asset.ClusterAsset) error {
	kubeClient, err := kubeclient.CreateClient(conf.Kubernetes.AdminKubeConfig)
	if err != nil {
		logrus.Debugf("Failed to create kubernetes client %v", err)
		return err
	}

	if err := waitForAPIReady(kubeClient); err != nil {
		logrus.Debugf("Failed while waiting for Kubernetes API to be ready: %v", err)
		return err
	}
	// Set kubeconfig environment variable
	err = os.Setenv("KUBECONFIG", conf.Kubernetes.AdminKubeConfig)
	if err != nil {
		logrus.Debugf("Failed to set environment variable KUBECONFIG: %v", err)
		return err
	}

	if err := waitForCoreAPIsReady(kubeClient.Discovery(), 60*time.Second); err != nil {
		logrus.Debugf("APIs not ready in time, proceeding with partial discovery: %v", err)
		return err
	}

	// Apply network plugin
	if err := applyNetworkPlugin(conf.Network.Plugin, conf.IsNestOS); err != nil {
		logrus.Debugf("Failed to apply network plugin: %v", err)
		return err
	}
	logrus.Debug("Network plugin deployment completed successfully.")

	if conf.Housekeeper.DeployHousekeeper {
		logrus.Debug("Starting deployment of Housekeeper...")
		if err := deployHousekeeper(conf.Housekeeper); err != nil {
			logrus.Debugf("Failed to deploy operator: %v", err)
			return err
		}
		logrus.Debug("Housekeeper deployment completed successfully.")
	}

	// Wait for pods to be ready
	if err := waitForPodsReady(kubeClient); err != nil {
		logrus.Debugf("Failed while waiting for pods to be in 'Ready' state: %v", err)
		return err
	}
	return nil
}

func waitForAPIReady(client *kubernetes.Clientset) error {
	apiTimeout := 60 * time.Minute
	ctx := context.Background()
	apiContext, cancel := context.WithTimeout(ctx, apiTimeout)
	logrus.Debugf("Waiting up to %v for the Kubernetes API ready...", apiTimeout)
	defer cancel()

	discovery := client.Discovery()
	wait.Until(func() {
		version, err := discovery.ServerVersion()
		if err == nil {
			logrus.Debugf("The Kubernetes API %s up", version)
			cancel()
		} else {
			logrus.Debugf("Still waiting for Kubernetes API ready: %v", err)
		}
	}, 2*time.Second, apiContext.Done())

	err := apiContext.Err()
	if err != nil && err != context.Canceled {
		logrus.Debugf("Failed to waiting for kubernetes API: %v", err)
		return err
	}
	return nil
}

func waitForPodsReady(client *kubernetes.Clientset) error {
	waitDuration := 20 * time.Minute
	waitCtx, cancel := context.WithTimeout(context.Background(), waitDuration)
	defer cancel()
	logrus.Debugf("Waiting up to %v for the Kubernetes Pods ready ...", waitDuration)

	err := wait.PollImmediate(10*time.Second, waitDuration, func() (bool, error) {
		pods, err := client.CoreV1().Pods(kubeSystemNS).List(waitCtx, metav1.ListOptions{})
		if err != nil {
			logrus.Debugf("Failed to list Pods: %v", err)
			return false, nil
		}
		allReady := true
		for _, pod := range pods.Items {
			for _, condition := range pod.Status.Conditions {
				if condition.Type == corev1.PodReady && condition.Status != corev1.ConditionTrue {
					allReady = false
					logrus.Debugf("Pod %s in namespace %s is not in Ready state", pod.Name, pod.Namespace)
					break
				}
			}
		}

		if allReady {
			// logrus.Infof("All Pods in namespace %s are in Ready state", namespace)
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		logrus.Debugf("failed to wait for Pods to be Ready: %v", err)
		return fmt.Errorf("failed to wait for Pods to be Ready: %v", err)
	}
	return nil
}

func waitForCoreAPIsReady(dc discovery.DiscoveryInterface, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	requiredGroups := []string{"apps", "apiextensions.k8s.io", "policy", "rbac.authorization.k8s.io"}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for core API groups to be ready")
		case <-ticker.C:
			serverGroups, err := dc.ServerGroups()
			if err != nil {
				logrus.Debugf("Failed to get server groups, retrying: %v", err)
				continue
			}

			found := make(map[string]bool)
			for _, g := range serverGroups.Groups {
				found[g.Name] = true
			}

			allReady := true
			for _, rg := range requiredGroups {
				if !found[rg] {
					allReady = false
					break
				}
			}

			if allReady {
				logrus.Debug("All required API groups are ready")
				return nil
			}

			logrus.Debugf("Waiting for API groups: %v (current: %v)", requiredGroups, getGroupNames(serverGroups))
		}
	}
}

func getGroupNames(sg *metav1.APIGroupList) []string {
	names := make([]string, len(sg.Groups))
	for i, g := range sg.Groups {
		names[i] = g.Name
	}
	return names
}

func deployHousekeeper(tmplData interface{}) error {
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
		if err != nil {
			return err
		}
		err = kubeclient.ApplyYAML(data)
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
		client := &http.Client{Timeout: time.Duration(60) * time.Second}
		response, err := client.Get(pluginConfigPath)
		if err != nil {
			return fmt.Errorf("failed to fetch network plugin configuration from URL %s: %w", pluginConfigPath, err)
		}
		defer response.Body.Close()

		if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
			return fmt.Errorf("unexpected HTTP status %d when fetching %s", response.StatusCode, pluginConfigPath)
		}

		content, err = io.ReadAll(response.Body)
		if err != nil {
			logrus.Debugf("Failed to read content from HTTP response: %v", err)
			return err
		}
	} else {
		// Read the content from the local file
		content, err = os.ReadFile(pluginConfigPath)
		if err != nil {
			logrus.Debugf("Failed to read network plugin configuration file: %v", err)
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

	if err := kubeclient.ApplyYAML(content); err != nil {
		logrus.Debugf("Failed to apply network plugin configuration: %v", err)
		return err
	}

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

func addKickstartFiles(httpService *httpserver.HTTPService, conf *asset.ClusterAsset) error {
	// Only one master node
	if err := httpService.AddFileToCache(conf.Master[0].Hostname+constants.KickstartSuffix, conf.BootConfig.Controlplane.Content); err != nil {
		return fmt.Errorf("error adding control plane kickstart file to cache: %v", err)
	}

	// multiple master nodes
	n := len(conf.Master)
	if n > 1 {
		for i := 1; i < n; i++ {
			if err := httpService.AddFileToCache(conf.Master[i].Hostname+constants.KickstartSuffix, conf.BootConfig.KickstartMaster[i-1].Content); err != nil {
				return fmt.Errorf("error adding master kickstart file to cache: %v", err)
			}
		}
	}

	if err := httpService.AddFileToCache(constants.Worker+constants.KickstartSuffix, conf.BootConfig.Worker.Content); err != nil {
		return fmt.Errorf("error adding worker kickstart file to cache: %v", err)
	}

	return nil
}
