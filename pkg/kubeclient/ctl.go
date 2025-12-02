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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	apiyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"

	"nestos-kubernetes-deployer/pkg/utils"
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

func ApplyYAML(yamlContent []byte) error {
	var config *rest.Config

	kubeconfig := utils.GetKubeConfigPath()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logrus.Errorf("Error get kubernetes config: %v", err)
		return err
	}

	yamlReader := apiyaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(yamlContent)))
	for {
		section, err := yamlReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Errorf("failed to read YAML document: %v", err)
			return err
		}

		unstructuredobj, err := parseYAMLToUnstructured(string(section))
		if err != nil {
			return err
		}

		if unstructuredobj.Object == nil || len(unstructuredobj.Object) == 0 {
			continue
		}

		if err := applySingleResources(kubeconfig, config, unstructuredobj); err != nil {
			gvk := unstructuredobj.GetObjectKind().GroupVersionKind()
			logrus.Errorf("failed to apply %s %s/%s: %v", gvk.Kind, unstructuredobj.GetNamespace(), unstructuredobj.GetName(), err)
			return err
		}
	}
	return nil
}

func applySingleResources(kubeconfig string, config *rest.Config, unstructuredobj *unstructured.Unstructured) error {
	var ns string

	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return err
	}
	groupResources, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		return err
	}
	mapper := restmapper.NewDiscoveryRESTMapper(groupResources)

	gvk := unstructuredobj.GetObjectKind().GroupVersionKind()
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("unable to get mapping for %v: %w", gvk, err)
	}

	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		ns = unstructuredobj.GetNamespace()
		if ns == "" {
			ns = "default"
		}
	}

	dyn, err := CreateDynamicClient(kubeconfig)
	if err != nil {
		return err
	}

	resource := dyn.Resource(mapping.Resource).Namespace(ns)
	applyConfig := &metav1.PatchOptions{
		Force:        pointer.Bool(false),
		FieldManager: "nkd-controller",
	}

	patchData, err := unstructuredobj.MarshalJSON()
	if err != nil {
		logrus.Errorf("failed to marshal object: %v", err)
		return err
	}

	_, err = resource.Patch(context.TODO(), unstructuredobj.GetName(), types.ApplyPatchType, patchData, *applyConfig)
	return err
}
