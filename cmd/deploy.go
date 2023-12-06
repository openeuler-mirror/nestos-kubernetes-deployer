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
	var clusterID = "cluster"
	opts.Opts.ClusterID = clusterID
	if err := configmanager.Initial(&opts.Opts); err != nil {
		logrus.Errorf("Failed to initialize configuration parameters: %v", err)
		return err
	}
	config, err := configmanager.GetClusterConfig(clusterID)
	if err != nil {
		logrus.Errorf("Failed to get cluster config using the cluster id: %v", err)
		return err
	}

	if err := deployCluster(config); err != nil {
		logrus.Errorf("Failed to deploy %s cluster: %v", clusterID, err)
		return err
	}
	if err := configmanager.Persist(); err != nil {
		logrus.Errorf("Failed to persist the cluster asset: %v", err)
		return err
	}

	return nil
}

func deployCluster(conf *asset.ClusterAsset) error {
	if err := getClusterDeployConfig(conf); err != nil {
		logrus.Errorf("Failed to get cluster deploy config: %v", err)
		return err
	}
	if err := createCluster(conf); err != nil {
		logrus.Errorf("Failed to create cluster: %v", err)
		return err
	}

	configPath := conf.Kubernetes.AdminKubeConfig
	kubeClient, err := kubeclient.CreateClient(configPath)
	if err != nil {
		logrus.Errorf("failed to create kubernetes client %v", err)
		return err
	}
	if err := checkClusterState(kubeClient); err != nil {
		logrus.Error("Cluster deploy timeout!")
		return err
	}

	if conf.Housekeeper.DeployHousekeeper {
		logrus.Info("Starting deployment of Housekeeper...")
		if err := deployOperator(conf.Housekeeper, kubeClient); err != nil {
			logrus.Errorf("Failed to deploy operator: %v", err)
			return err
		}
		logrus.Info("Housekeeper deployment completed successfully.")
	}

	return nil
}

func getClusterDeployConfig(conf *asset.ClusterAsset) error {
	if err := generateCerts(conf); err != nil {
		logrus.Errorf("Error generating certificate files: %v", err)
		return err
	}

	if err := generateIgnition(conf); err != nil {
		logrus.Errorf("Error generating ignition files: %v", err)
		return err
	}

	if err := generateTF(conf); err != nil {
		logrus.Errorf("Error generating terraform files: %v", err)
		return err
	}

	return nil
}

func generateCerts(conf *asset.ClusterAsset) error {
	// Generate CA certificates
	caCerts, err := cert.GenerateCAFiles(conf.Cluster_ID)
	if err != nil {
		logrus.Errorf("Error generating CA files: %v", err)
		return err
	}

	sameCerts, err := cert.GenerateCertFilesAllSame(conf.Cluster_ID)
	if err != nil {
		return err
	}

	// Generate certificates for each Master node
	for i, master := range conf.Master {
		var masterCerts []utils.StorageContent

		certs, err := cert.GenerateCertFilesForNode(&master)
		if err != nil {
			logrus.Errorf("Error generating certificate files for Master %d: %v", i, err)
			return err
		}
		masterCerts = append(masterCerts, caCerts...)
		masterCerts = append(masterCerts, sameCerts...)
		masterCerts = append(masterCerts, certs...)

		// Assign the certificates to the corresponding Master node
		conf.Master[i].Certs = masterCerts
	}

	return nil
}

func generateIgnition(conf *asset.ClusterAsset) error {
	master := &machine.Master{
		ClusterAsset: conf,
	}
	if err := master.GenerateFiles(); err != nil {
		logrus.Errorf("Failed to generate master ignition file: %v", err)
		return err
	}

	worker := &machine.Worker{
		ClusterAsset: conf,
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
	persistDir := configmanager.GetPersistDir()
	masterInfra := infra.InstanceCluster(persistDir, conf.Cluster_ID, "master", len(conf.Master))
	if err := masterInfra.Deploy(); err != nil {
		logrus.Errorf("Failed to deploy master nodes:%v", err)
		return err
	}
	workerInfra := infra.InstanceCluster(persistDir, conf.Cluster_ID, "worker", len(conf.Worker))
	if err := workerInfra.Deploy(); err != nil {
		logrus.Errorf("Failed to deploy worker nodes:%v", err)
		return err
	}

	return nil
}

func checkClusterState(client *kubernetes.Clientset) error {
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

func deployOperator(tmplData interface{}, client *kubernetes.Clientset) error {
	folderPath := "housekeeper/"
	files, err := os.ReadDir(folderPath)
	if err != nil {
		logrus.Errorf("Error reading folder: %v", err)
		return err
	}
	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())
		data, err := utils.FetchAndUnmarshalUrl(filePath, tmplData)
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

func waitForPodsRunning(client *kubernetes.Clientset) error {
	waitDuration := 10 * time.Minute
	waitCtx, cancel := context.WithTimeout(context.Background(), waitDuration)
	logrus.Infof("Waiting up to %v for the Kubernetes Pods running ...", waitDuration)
	defer cancel()

	wait.Until(func() {
		pods, err := client.CoreV1().Pods("kube-system").List(waitCtx, metav1.ListOptions{})
		if err != nil {
			logrus.Errorf("Failed to list Pods: %v", err)
			return
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
