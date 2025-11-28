/*
Copyright 2024 KylinSoft  Co., Ltd.

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
	"github.com/agiledragon/gomonkey/v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"testing"
)

var (
	kubeconfigPath = "./kubeconfig"
	yamlContent    = ""
	//namespace      = "test_"
)

func TestCreateClient(t *testing.T) {
	t.Run("CreateClient_fail", func(t *testing.T) {
		clientset, err := CreateClient(kubeconfigPath)

		if err != nil {
			t.Logf("CreateClient returned  error: %v", err)
			return
		}
		if clientset == nil {
			t.Logf("CreateClient is empty")
			return
		}

		log.Println("TestCreateClient  success")
	})

	t.Run("CreateDynamicClient_fail", func(t *testing.T) {
		CreateDynamicClient, err := CreateDynamicClient(kubeconfigPath)
		if err != nil {
			t.Logf("CreateDynamicClient returned  error: %v", err)
			return
		}
		if CreateDynamicClient == nil {
			t.Logf("CreateDynamicClient is empty")
			return
		}

		log.Println("TestCreateDynamicClient success")
	})

	p := gomonkey.ApplyFunc(clientcmd.BuildConfigFromFlags, func(string, string) (*rest.Config, error) {
		return &rest.Config{}, nil
	})

	defer p.Reset()

	t.Run("CreateClient_dyn_fail", func(t *testing.T) {
		clientset, err := CreateClient(kubeconfigPath)

		if err != nil {
			t.Logf("CreateClient returned  error: %v", err)
			return
		}
		if clientset == nil {
			t.Logf("CreateClient is empty")
			return
		}

		log.Println("TestCreateClient  success")
	})

	t.Run("CreateDynamicClient_dyn_fail", func(t *testing.T) {
		CreateDynamicClient, err := CreateDynamicClient(kubeconfigPath)
		if err != nil {
			t.Logf("CreateDynamicClient returned  error: %v", err)
			return
		}
		if CreateDynamicClient == nil {
			t.Logf("CreateDynamicClient is empty")
			return
		}

		log.Println("TestCreateDynamicClient success")
	})

	kp := gomonkey.ApplyFunc(kubernetes.NewForConfig, func(config *rest.Config) (*kubernetes.Clientset, error) {
		return &kubernetes.Clientset{}, nil
	})

	defer kp.Reset()

	t.Run("CreateClient", func(t *testing.T) {
		clientset, err := CreateClient(kubeconfigPath)

		if err != nil {
			t.Logf("CreateClient returned  error: %v", err)
			return
		}
		if clientset == nil {
			t.Logf("CreateClient is empty")
			return
		}

		log.Println("TestCreateClient  success")
	})

	t.Run("CreateDynamicClient", func(t *testing.T) {
		CreateDynamicClient, err := CreateDynamicClient(kubeconfigPath)
		if err != nil {
			t.Logf("CreateDynamicClient returned  error: %v", err)
			return
		}
		if CreateDynamicClient == nil {
			t.Logf("CreateDynamicClient is empty")
			return
		}

		log.Println("TestCreateDynamicClient success")
	})

}

