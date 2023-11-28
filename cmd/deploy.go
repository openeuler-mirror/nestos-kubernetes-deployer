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
	"nestos-kubernetes-deployer/cmd/command"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/cert"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/constants"
	"nestos-kubernetes-deployer/pkg/ignition"
	"nestos-kubernetes-deployer/pkg/ignition/machine"
	"nestos-kubernetes-deployer/pkg/infra"
	"nestos-kubernetes-deployer/pkg/kubeclient"
	"nestos-kubernetes-deployer/pkg/utils"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	wait "k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

type crdTmplData struct {
	operatorImageUrl   string
	controllerImageUrl string
}

type Certificates []ignition.StorageContent

func NewDeployCommand() *cobra.Command {
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a kubernetes cluster",
		RunE:  runDeployCmd,
	}
	command.SetupDeployCmdOpts(deployCmd)

	return deployCmd
}

func runDeployCmd(cmd *cobra.Command, args []string) error {
	if err := configmanager.Initial(&opts.Opts); err != nil {
		logrus.Errorf("Failed to initialize configuration parameters: %v", err)
		return err
	}
	config, err := configmanager.GetClusterConfig("clusterId")
	if err != nil {
		logrus.Errorf("Failed to get cluster config using the cluster id: %v", err)
		return err
	}

	if err := deployCluster(config); err != nil {
		return err
	}
	err = configmanager.Persist()

	return nil
}

func deployCluster(conf *asset.ClusterAsset) error {
	if err := getClusterDeployConfig(conf); err != nil {
		return err
	}
	if err := createCluster(conf); err != nil {
		return err
	}

	configPath := filepath.Join(opts.RootOptDir, "auth", "kubeconfig")
	if err := checkClusterState(configPath); err != nil {
		logrus.Error("Cluster deploy timeout!")
		return err
	}

	/*调用配置管理模块接口，获取crdTmplData数据*/

	if err := deployOperator( /**/ ); err != nil {
		logrus.Errorf("Failed to deploy operator: %v", err)
		return err
	}

	return nil
}

func getClusterDeployConfig(conf *asset.ClusterAsset) error {
	certs, err := generateCerts(conf)
	if err != nil {
		logrus.Errorf("Error generating certificate files: %v", err)
		return err
	}

	if err := generateIgnition(conf, certs); err != nil {
		logrus.Errorf("Error generating ignition files: %v", err)
		return err
	}

	if err := generateTF(conf); err != nil {
		logrus.Errorf("Error generating terraform files: %v", err)
		return err
	}

	return nil
}

func generateCerts(conf *asset.ClusterAsset) (Certificates, error) {
	rootCA, err := cert.GenerateRootCA()
	if err != nil {
		logrus.Errorf("Error generating root CA:%v", err)
		return nil, err
	}
	etcdCA, err := cert.GenerateEtcdCA()
	if err != nil {
		logrus.Errorf("Error generating etcd CA:%v", err)
		return nil, err
	}
	frontProxyCA, err := cert.GenerateFrontProxyCA()
	if err != nil {
		logrus.Errorf("Error generating front proxy CA:%v", err)
		return nil, err
	}

	var certs Certificates
	ignition.UpsertStorageFiles(certs, constants.CaCrt, int(constants.CertFileMode), rootCA.CertRaw)
	ignition.UpsertStorageFiles(certs, constants.CaKey, int(constants.CertFileMode), rootCA.KeyRaw)
	ignition.UpsertStorageFiles(certs, constants.EtcdCaCrt, int(constants.CertFileMode), etcdCA.CertRaw)
	ignition.UpsertStorageFiles(certs, constants.EtcdCaKey, int(constants.CertFileMode), etcdCA.KeyRaw)
	ignition.UpsertStorageFiles(certs, constants.FrontProxyCaCrt, int(constants.CertFileMode), frontProxyCA.CertRaw)
	ignition.UpsertStorageFiles(certs, constants.FrontProxyCaKey, int(constants.CertFileMode), frontProxyCA.KeyRaw)

	return certs, nil
}

func generateIgnition(conf *asset.ClusterAsset, certFiles []ignition.StorageContent) error {
	master := &machine.Master{
		ClusterAsset:   conf,
		StorageContent: certFiles,
		IgnFiles:       []ignition.IgnFile{},
	}
	if err := master.GenerateFiles(); err != nil {
		logrus.Errorf("Failed to generate master ignition file: %v", err)
		return err
	}

	worker := &machine.Worker{
		ClusterAsset: conf,
		IgnFiles:     []ignition.IgnFile{},
	}
	if err := worker.GenerateFiles(); err != nil {
		logrus.Errorf("Failed to generate worker ignition file: %v", err)
		return err
	}

	return nil
}

func generateTF(conf *asset.ClusterAsset) error {
	// generate master.tf
	var master infra.Infra
	if err := master.Generate(conf, "master"); err != nil {
		logrus.Errorf("Failed to generate master terraform file")
		return err
	}
	// generate worker.tf
	var worker infra.Infra
	if err := worker.Generate(conf, "worker"); err != nil {
		logrus.Errorf("Failed to generate worker terraform file")
		return err
	}
	return nil
}

func createCluster(conf *asset.ClusterAsset) error {

	/*应用集群配置文件部署集群*/
	return nil
}

func checkClusterState(kubeconfigPath string) error {
	client, err := kubeclient.CreateClient(kubeconfigPath)
	if err != nil {
		logrus.Errorf("failed to create kubernetes client %v", err)
		return err
	}
	if err := waitForAPIReady(client); err != nil {
		logrus.Errorf("failed while waiting for Kubernetes API to be ready: %v", err)
		return err
	}
	if err := waitForPodsRunning(client); err != nil {
		logrus.Errorf("failed while waiting for pods to be in 'Running' state: %v", err)
		return err
	}
	return nil
}

func deployOperator(folderPath string, client *kubernetes.Clientset) error {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		logrus.Errorf("Error reading folder: %v", err)
		return err
	}
	// 实例化crdTmplData
	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())
		data, err := utils.FetchAndUnmarshalUrl(filePath, "" /*获取crdTmplData数据*/)
		if err != nil {
			logrus.Errorf("Error to get file content: %v", err)
			return err
		}
		if err := kubeclient.ApplyResource(client, "", string(data)); err != nil {
			logrus.Errorf("Error to apply crd resource: %v", err)
			return err
		}
	}

	return nil
}

func waitForAPIReady(client *kubernetes.Clientset) error {
	apiTimeout := 10 * time.Minute
	ctx := context.Background()
	apiContext, cancel := context.WithTimeout(ctx, apiTimeout)
	logrus.Infof("Waiting up to %v for the Kubernetes API at %s...", apiTimeout, config.Host)
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

func waitForPodsRunning(client *kubernetes.Clientset) error {
	waitDuration := 10 * time.Minute
	waitCtx, cancel := context.WithTimeout(context.Background(), waitDuration)
	logrus.Infof("Waiting up to %v for the Kubernetes Pods running ...", waitDuration)
	defer cancel()

	wait.Until(func() {
		pods, err := client.CoreV1().Pods("kube-system").List(waitCtx, metav1.ListOptions{})
		if err != nil {
			logrus.Errorf("Failed to list Pods: %v", err)
			return err
		}
		allRunning := true
		for _, pod := range pods.Items {
			if pod.Status.Phase != corev1.PodRunning {
				allRunning = false
				logrus.Infof("Pod %s is not running. Current phase: %s", pod.Name, pod.Status.Phase)
				break
			}
		}
		if allRunning {
			logrus.Info("All Pods are running")
			cancel()
		}
	}, 5*time.Second, waitCtx.Done())

	err := waitCtx.Err()
	if err != nil && err != context.Canceled {
		logrus.Errorf("Failed to wait for Pods to be running: %v", err)
		return err
	}
	return nil
}
