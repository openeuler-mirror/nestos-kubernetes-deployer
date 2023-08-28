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
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"

	wait "k8s.io/apimachinery/pkg/util/wait"
)

var (
	ArgsData struct {
		osVersion      string
		kubeVersion    string
		evictPodForce  bool
		maxUnavailable int
		kubeconfig     string
		osImageURL     string
	}
)

func NewUpgradeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade your cluster to a newer version",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpgradeCmd()
		},
	}

	cmd.PersistentFlags().StringVarP(&ArgsData.osVersion, "nestos-version", "", "", "Choose a specific nestos version for upgrade or rollback")
	cmd.PersistentFlags().StringVarP(&ArgsData.kubeVersion, "kubernetes-version", "", "", "Choose a specific kubernetes version for upgrading")
	cmd.PersistentFlags().BoolVarP(&ArgsData.evictPodForce, "force", "f", false, "Force evict pod")
	cmd.PersistentFlags().IntVarP(&ArgsData.maxUnavailable, "maxunavailable", "n", 2, "Number of nodes that are upgraded at the same time")
	cmd.PersistentFlags().StringVarP(&ArgsData.kubeconfig, "kubeconfig", "", "./auth/config", "kubeconfig file access path")
	cmd.PersistentFlags().StringVarP(&ArgsData.osImageURL, "imageurl", "", "", "The address of the container image to use for upgrading")
	return cmd
}

func runUpgradeCmd() error {
	loopTimeout := 2 * time.Minute
	// Get the kubeconfig configuration
	config, err := clientcmd.BuildConfigFromFlags("", ArgsData.kubeconfig)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			logrus.Errorf("Error getting Kubernetes client config: %v\n", err)
			return err
		}
	}

	// Create dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		logrus.Errorf("Error creating Dynamic client: %v\n", err)
		return err
	}

	// Define the YAML data for the Custom Resource (CR)
	yamlData := fmt.Sprintf(`
apiVersion: housekeeper.io/v1alpha1
kind: Update
metadata:
name: housekeeper-upgrade
namespace: housekeeper-system
spec:
osVersion: %s
osImageURL: %s
kubeVersion: %s
evictPodForce: %t
maxUnavailable: %d
`, ArgsData.osVersion, ArgsData.osImageURL, ArgsData.kubeVersion, ArgsData.evictPodForce, ArgsData.maxUnavailable)

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
