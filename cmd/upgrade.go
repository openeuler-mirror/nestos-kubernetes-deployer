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
	"fmt"
	"nestos-kubernetes-deployer/cmd/command"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/kubeclient"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	wait "k8s.io/apimachinery/pkg/util/wait"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	clusterId, err := cmd.Flags().GetString("cluster-id")
	if err != nil {
		logrus.Errorf("Failed to get cluster-id: %v", err)
		return err
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
	loopTimeout := 2 * time.Minute
	dynamicClient, err := kubeclient.CreateDynamicClient("/***/")

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

	var unstructuredObj unstructured.Unstructured
	err = yaml.Unmarshal([]byte(yamlData), &unstructuredObj)
	if err != nil {
		logrus.Errorf("Error unmarshalling YAML: %v\n", err)
		return err
	}

	// Create or Update CR
	resource := schema.GroupVersionResource{
		Group:    "housekeeper.io",
		Version:  "v1alpha1",
		Resource: "updates", // Pluralized resource name
	}

	// The loop attempts to create or update a CR until it succeeds or times out
	if err := wait.PollImmediate(2*time.Second, loopTimeout, func() (bool, error) {
		gvk := unstructuredObj.GroupVersionKind()
		dynamicResource := dynamicClient.Resource(gvk.GroupVersion().WithResource(resource.Resource)).Namespace(unstructuredObj.GetNamespace())

		//Attempts to get the specified Custom Resource from the Kubernetes API Server.
		obj, err := dynamicResource.Get(context.Background(), unstructuredObj.GetName(), metav1.GetOptions{})
		if err != nil {
			// Not found, create the resource
			_, err = dynamicResource.Create(context.Background(), &unstructuredObj, metav1.CreateOptions{})
			if err == nil {
				logrus.Infof("Custom Resource created successfully!")
				return true, nil
			}
		} else {
			// Found, update the resource
			unstructuredObj.SetResourceVersion(obj.GetResourceVersion())
			_, err = dynamicResource.Update(context.Background(), &unstructuredObj, metav1.UpdateOptions{})
			if err == nil {
				logrus.Infof("Custom Resource updated successfully!")
				return true, nil
			}
		}
		logrus.Errorf("Error creating or updating CR: %v\n", err)
		return false, nil
	}); err != nil {
		logrus.Errorf("Timeout while waiting for Custom Resource to be created or updated.")
	}

	return nil
}
