package main

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

func newUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade your cluster to a newer version",
		Long:  "",
		RunE:  runUpgradeCmd,
	}
	return cmd
}

func runUpgradeCmd(command *cobra.Command, args []string) error {
	var (
		osVersion      = ""
		osImageURL     = ""
		kubeVersion    = ""
		evictPodForce  = false
		maxUnavailable = 2
		loopTimeout    = 2 * time.Minute
		kubeconfig     = ""
	)
	// Get the kubeconfig configuration
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
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
  name: update-sample
  namespace: housekeeper-system
spec:
  osVersion: %s
  osImageURL: %s
  kubeVersion: %s
  evictPodForce: %t
  maxUnavailable: %d
`, osVersion, osImageURL, kubeVersion, evictPodForce, maxUnavailable)

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
