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
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	wait "k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewDeployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a kubernetes cluster",
		RunE:  runDeployCmd,
	}

	cmd.PersistentFlags().StringVar(&command.ClusterOpts.ClusterId, "cluster-id", "", "clusterID of kubernetes cluster")
	cmd.PersistentFlags().StringVar(&command.ClusterOpts.GatherDeployOpts.SSHKey, "sshkey", "", "Path to SSH private keys that should be used for authentication.")
	cmd.PersistentFlags().StringVar(&command.ClusterOpts.Platform, "platform", "", "Select the infrastructure platform to deploy the cluster")

	// cmd.AddCommand(deploy.NewDeployMasterCommand())
	// cmd.AddCommand(deploy.NewDeployWorkerCommand())

	return cmd
}

func runDeployCmd(cmd *cobra.Command, args []string) error {

	//todo：部署集群

	configPath := filepath.Join(command.RootOptDir, "auth", "kubeconfig")
	if err := checkClusterState(configPath); err != nil {
		logrus.Error("Cluster deploy timeout!")
		return err
	}
	return nil
}

// 生成部署集群所需配置数据
func runInstallconfig() error {

	return nil
}

func runDeployCluster() error {

	return nil
}

func checkClusterState(kubeconfigPath string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		logrus.Errorf("error to load kubeconfig: %v", err)
		return err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Errorf("failed to create a kubernetes client: %v", err)
		return err
	}

	if err := waitForAPIReady(client); err != nil {
		return err
	}

	if err := waitForPodsRunning(client); err != nil {
		return err
	}
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

	err = apiContext.Err()
	if err != nil && err != context.Canceled {
		logrus.Errorf("Failed to waiting for kubernetes API: %v", err)
		return err
	}

	return nil
}

func waitForPodsRunning(client *kubernetes.Clientset) error {
	waitDuration := 10 * time.Minute
	waitCtx, cancel := context.WithTimeout(context.Background(), waitkDuration)
	logrus.Infof("Waiting up to %v for the Kubernetes Pods running ...", waitkDuration)
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
