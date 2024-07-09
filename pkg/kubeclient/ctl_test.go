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

	t.Run("DeployCRD", func(t *testing.T) {
		DeployCRD(yamlContent, kubeconfigPath)
	})
	t.Run("DeployNamespace", func(t *testing.T) {
		DeployNamespace(yamlContent, kubeconfigPath)
	})

	t.Run("DeployClusterRoleBinding", func(t *testing.T) {
		DeployClusterRoleBinding(yamlContent, kubeconfigPath)
	})

	//t.Run("DeployDeployment", func(t *testing.T) {
	//	//patches := gomonkey.ApplyMethod(reflect.TypeOf(&appsv1.Deployment{}), "Create", func(_ *appsv1.Deployment, ctx context.Context, deployment *appsv1.Deployment, opts metav1.CreateOptions) (*appsv1.Deployment, error) {
	//	//	return &appsv1.Deployment{
	//	//		ObjectMeta: metav1.ObjectMeta{
	//	//			Name: deployment.Name,
	//	//		},
	//	//	}, nil
	//	//})
	//	//defer patches.Reset()
	//	d := appsv1.Deployment{}
	//	t.Log(d)
	//	DeployDeployment(yamlContent, kubeconfigPath, namespace)
	//})
	//
	//t.Run("DeployDaemonSet", func(t *testing.T) {
	//	DeployDaemonSet(yamlContent, kubeconfigPath, namespace)
	//})
	t.Run("ApplyHousekeeperCR", func(t *testing.T) {
		ApplyHousekeeperCR(yamlContent, kubeconfigPath)
	})
	t.Run("RunKubectlApplyWithYaml", func(t *testing.T) {
		RunKubectlApplyWithYaml(yamlContent)
	})

}

func TestIsKubectlInstalled(t *testing.T) {
	b := IsKubectlInstalled()
	if !b {
		t.Log("no install")
		return
	}
	t.Log("IsKubectlInstalled success")
}
