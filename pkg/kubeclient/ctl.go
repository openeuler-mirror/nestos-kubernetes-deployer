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
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
		logrus.Errorf("error loading kubeconfig: %v", err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Errorf("failed to create a Kubernetes client: %v", err)
		return nil, err
	}

	return clientset, nil
}

// Apply a Kubernetes resource of the specified type using the provided content.
// Parameters:
//   - clientset: Kubernetes clientset for cluster interaction.
//     Input: *kubernetes.Clientset - configured Kubernetes client.
//   - resourceType: Type of Kubernetes resource (e.g., "pods", "services").
//     Input: string - type of the Kubernetes resource.
//   - content: YAML or JSON content for creating/updating the resource.
//     Input: string - content of the Kubernetes resource.
func ApplyResource(clientset *kubernetes.Clientset, resourceType, content string) error {
	_, err := clientset.RESTClient().
		Post().
		Resource(resourceType).
		Body([]byte(content)).
		Do(context.TODO()).
		Get()

	if err != nil {
		logrus.Errorf("Error applying content: %v", err)
		return err
	}
	return nil
}
