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
package kubeclient

import (
	"context"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

// CreateClient creates a Kubernetes clientset.
// Parameters:
// - kubeconfig: Path to the kubeconfig file.
//               Input: string - kubeconfig file path.
// Returns:
//     Output: *kubernetes.Clientset - Kubernetes client.
//   - error: Error

func CreateClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logrus.Errorf("Error loading kubeconfig: %v", err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Errorf("Failed to create a Kubernetes client: %v", err)
		return nil, err
	}

	return clientset, nil
}

// CreateDynamicClient creates a dynamic client.
func CreateDynamicClient(kubeconfig string) (dynamic.Interface, error) {
	// Get the kubeconfig configuration
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			logrus.Errorf("Error getting Kubernetes client config: %v\n", err)
			return nil, err
		}
	}

	// Create dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		logrus.Errorf("Error creating Dynamic client: %v\n", err)
		return nil, err
	}

	return dynamicClient, nil
}

func DeployCRD(yamlContent string, kubeconfig string) error {
	client, err := CreateDynamicClient(kubeconfig)
	if err != nil {
		return err
	}

	// Parse YAML into CustomResourceDefinition
	unstructuredObj := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(yamlContent), unstructuredObj); err != nil {
		logrus.Errorf("Error parsing YAML as Unstructured: %v", err)
		return err
	}

	// Specify the API group, version, and resource for CustomResourceDefinitions
	apiGroup := "apiextensions.k8s.io"
	apiVersion := "v1"
	resource := "customresourcedefinitions"

	// Create the CRD using the dynamic client
	_, err = client.Resource(schema.GroupVersionResource{
		Group:    apiGroup,
		Version:  apiVersion,
		Resource: resource,
	}).Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("error creating CRD: %v", err)
		return err
	}

	return nil
}

func DeployNamespace(yamlContent string, kubeconfig string) error {
	client, err := CreateDynamicClient(kubeconfig)
	if err != nil {
		logrus.Errorf("error creating dynamic client: %v", err)
		return err
	}

	// Parse YAML content into Unstructured object
	unstructuredObj := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(yamlContent), unstructuredObj); err != nil {
		logrus.Errorf("error converting YAML to Unstructured: %v", err)
		return err
	}

	// Specify the API group, version, and resource for Namespaces
	apiGroup, apiVersion, resource := "", "v1", "namespaces"

	// Create the Namespace using the dynamic client
	_, err = client.Resource(schema.GroupVersionResource{
		Group:    apiGroup,
		Version:  apiVersion,
		Resource: resource,
	}).Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("error creating Namespace: %v", err)
		return err
	}

	return nil
}

func DeployClusterRole(yamlContent string, kubeconfig string) error {
	client, err := CreateDynamicClient(kubeconfig)
	if err != nil {
		return err
	}

	unstructuredObj := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(yamlContent), unstructuredObj); err != nil {
		logrus.Errorf("Error parsing YAML as Unstructured: %v", err)
		return err
	}

	apiGroup := "rbac.authorization.k8s.io"
	apiVersion := "v1"
	resource := "clusterroles"

	_, err = client.Resource(schema.GroupVersionResource{
		Group:    apiGroup,
		Version:  apiVersion,
		Resource: resource,
	}).Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("error creating CRD: %v", err)
		return err
	}
	return nil
}

func DeployClusterRoleBinding(yamlContent string, kubeconfig string) error {
	client, err := CreateDynamicClient(kubeconfig)
	if err != nil {
		return err
	}

	unstructuredObj := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(yamlContent), unstructuredObj); err != nil {
		logrus.Errorf("Error parsing YAML as Unstructured: %v", err)
		return err
	}

	apiGroup := "rbac.authorization.k8s.io"
	apiVersion := "v1"
	resource := "clusterrolebindings"

	_, err = client.Resource(schema.GroupVersionResource{
		Group:    apiGroup,
		Version:  apiVersion,
		Resource: resource,
	}).Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("error creating CRD: %v", err)
		return err
	}

	return nil
}

func DeployDeployment(yamlContent string, kubeconfig string, namespace string) error {
	clientset, err := CreateClient(kubeconfig)
	if err != nil {
		return err
	}

	// Parse YAML content into Unstructured object
	unstructuredObj := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(yamlContent), unstructuredObj); err != nil {
		logrus.Errorf("error converting YAML to Unstructured: %v", err)
		return err
	}

	deployment := &appsv1.Deployment{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.Object, deployment)
	if err != nil {
		logrus.Errorf("error converting Unstructured to deployment: %v", err)
		return err
	}

	// Create the Deployment using the Kubernetes clientset
	_, err = clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("error creating Deployment: %v", err)
		return err
	}

	return nil
}

func DeployDaemonSet(yamlContent string, kubeconfig string, namespace string) error {
	clientset, err := CreateClient(kubeconfig)
	if err != nil {
		return err
	}

	// Parse YAML content into Unstructured object
	unstructuredObj := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(yamlContent), unstructuredObj); err != nil {
		logrus.Errorf("error converting YAML to Unstructured: %v", err)
		return err
	}

	daemonSet := &appsv1.DaemonSet{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.Object, daemonSet)
	if err != nil {
		logrus.Errorf("error converting Unstructured to daemonset: %v", err)
		return err
	}

	// Create the DaemonSet using the Kubernetes clientset
	_, err = clientset.AppsV1().DaemonSets(namespace).Create(context.TODO(), daemonSet, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("error creating DaemonSet: %v", err)
		return err
	}

	return nil
}

func DeployCR(yamlContent string, kubeconfig string) error {
	// Create a dynamic client for interacting with the Kubernetes API server
	client, err := CreateDynamicClient(kubeconfig)
	if err != nil {
		return err
	}

	// Parse the YAML content into an Unstructured object
	unstructuredObj := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(yamlContent), unstructuredObj); err != nil {
		logrus.Errorf("Error parsing YAML as Unstructured: %v", err)
		return err
	}

	// Specify the API group, version, and resource for the custom resource
	apiGroup := "housekeeper.io"
	apiVersion := "v1alpha1"
	resource := "updates"

	// Try to get the existing custom resource
	existingObj, err := client.
		Resource(schema.GroupVersionResource{
			Group:    apiGroup,
			Version:  apiVersion,
			Resource: resource,
		}).
		Namespace(unstructuredObj.GetNamespace()).
		Get(context.TODO(), unstructuredObj.GetName(), metav1.GetOptions{})

	if err != nil {
		if errors.IsNotFound(err) {
			// Custom resource doesn't exist, create it
			_, err = client.
				Resource(schema.GroupVersionResource{
					Group:    apiGroup,
					Version:  apiVersion,
					Resource: resource,
				}).
				Namespace(unstructuredObj.GetNamespace()).
				Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
			if err != nil {
				logrus.Errorf("Error creating custom resource: %v", err)
				return err
			}
			return nil
		}

		// Error other than "not found" occurred
		logrus.Errorf("Error checking custom resource existence: %v", err)
		return err
	}

	// Custom resource already exists, update it with new configuration
	unstructuredObj.SetResourceVersion(existingObj.GetResourceVersion())
	_, err = client.
		Resource(schema.GroupVersionResource{
			Group:    apiGroup,
			Version:  apiVersion,
			Resource: resource,
		}).
		Namespace(unstructuredObj.GetNamespace()).
		Update(context.TODO(), unstructuredObj, metav1.UpdateOptions{})
	if err != nil {
		logrus.Errorf("Error updating custom resource: %v", err)
		return err
	}

	return nil
}
