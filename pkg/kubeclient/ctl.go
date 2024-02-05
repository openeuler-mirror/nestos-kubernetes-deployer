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
	"os/exec"

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

const (
	// Constants for CRD API groups, versions, and resources
	CRDAPIGroup   = "apiextensions.k8s.io"
	CRDAPIVersion = "v1"
	CRDResource   = "customresourcedefinitions"

	// custom resource
	HousekeeperAPIGroup   = "housekeeper.io"
	HousekeeperAPIVersion = "v1alpha1"
	HousekeeperResource   = "updates"

	// NAMESPACE
	NSResource   = "namespaces"
	NSAPIVersion = "v1"

	// RBAC
	RBACAPIGroup                = "rbac.authorization.k8s.io"
	RBACAPIVersion              = "v1"
	ClusterRolesResource        = "clusterroles"
	ClusterRoleBindingsResource = "clusterrolebindings"
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

// parseYAMLToUnstructured parses YAML into Unstructured object
func parseYAMLToUnstructured(yamlContent string) (*unstructured.Unstructured, error) {
	unstructuredObj := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(yamlContent), unstructuredObj); err != nil {
		logrus.Errorf("Error parsing YAML as Unstructured: %v", err)
		return nil, err
	}
	return unstructuredObj, nil
}

// deployResource deploys a resource using dynamic client
func deployResource(yamlContent, kubeconfig string, apiGroup, apiVersion, resource string) error {
	client, err := CreateDynamicClient(kubeconfig)
	if err != nil {
		return err
	}

	unstructuredObj, err := parseYAMLToUnstructured(yamlContent)
	if err != nil {
		return err
	}

	_, err = client.Resource(schema.GroupVersionResource{
		Group:    apiGroup,
		Version:  apiVersion,
		Resource: resource,
	}).Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("Error creating resource %s: %v", resource, err)
		return err
	}

	return nil
}

// DeployCRD deploys a CustomResourceDefinition.
func DeployCRD(yamlContent string, kubeconfig string) error {
	return deployResource(yamlContent, kubeconfig, CRDAPIGroup, CRDAPIVersion, CRDResource)
}

// DeployNamespace deploys a Namespace.
func DeployNamespace(yamlContent string, kubeconfig string) error {
	return deployResource(yamlContent, kubeconfig, "", NSAPIVersion, NSResource)
}

// DeployClusterRole deploys a ClusterRole.
func DeployClusterRole(yamlContent string, kubeconfig string) error {
	return deployResource(yamlContent, kubeconfig, RBACAPIGroup, RBACAPIVersion, ClusterRolesResource)
}

// DeployClusterRoleBinding deploys a ClusterRoleBinding.
func DeployClusterRoleBinding(yamlContent string, kubeconfig string) error {
	return deployResource(yamlContent, kubeconfig, RBACAPIGroup, RBACAPIVersion, ClusterRoleBindingsResource)
}

// DeployDeployment deploys a Deployment.
func DeployDeployment(yamlContent string, kubeconfig string, namespace string) error {
	clientset, err := CreateClient(kubeconfig)
	if err != nil {
		return err
	}

	unstructuredObj, err := parseYAMLToUnstructured(yamlContent)
	if err != nil {
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

// DeployDaemonSet deploys a DaemonSet.
func DeployDaemonSet(yamlContent string, kubeconfig string, namespace string) error {
	clientset, err := CreateClient(kubeconfig)
	if err != nil {
		return err
	}

	unstructuredObj, err := parseYAMLToUnstructured(yamlContent)
	if err != nil {
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

func ApplyHousekeeperCR(yamlContent string, kubeconfig string) error {
	// Create a dynamic client for interacting with the Kubernetes API server
	client, err := CreateDynamicClient(kubeconfig)
	if err != nil {
		return err
	}

	// Parse the YAML content into an Unstructured object
	unstructuredObj, err := parseYAMLToUnstructured(yamlContent)
	if err != nil {
		return err
	}

	// Try to get the existing custom resource
	existingObj, err := client.
		Resource(schema.GroupVersionResource{
			Group:    HousekeeperAPIGroup,
			Version:  HousekeeperAPIVersion,
			Resource: HousekeeperResource,
		}).
		Namespace(unstructuredObj.GetNamespace()).
		Get(context.TODO(), unstructuredObj.GetName(), metav1.GetOptions{})

	if err != nil {
		if errors.IsNotFound(err) {
			// Custom resource doesn't exist, create it
			_, err = client.
				Resource(schema.GroupVersionResource{
					Group:    HousekeeperAPIGroup,
					Version:  HousekeeperAPIVersion,
					Resource: HousekeeperResource,
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
			Group:    HousekeeperAPIGroup,
			Version:  HousekeeperAPIVersion,
			Resource: HousekeeperResource,
		}).
		Namespace(unstructuredObj.GetNamespace()).
		Update(context.TODO(), unstructuredObj, metav1.UpdateOptions{})
	if err != nil {
		logrus.Errorf("Error updating custom resource: %v", err)
		return err
	}

	return nil
}

func RunKubectlApplyWithYaml(yamlFilePath string) error {
	kubectlArgs := []string{"apply", "-f", yamlFilePath}
	cmd := exec.Command("kubectl", kubectlArgs...)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// run kubectl apply
	err := cmd.Run()
	if err != nil {
		logrus.Errorf("Error executing kubectl apply: %v", err)
		return err
	}

	return nil
}

// isKubectlInstalled checks if kubectl is installed on the system.
func IsKubectlInstalled() bool {
	_, err := exec.LookPath("kubectl")
	return err == nil
}
