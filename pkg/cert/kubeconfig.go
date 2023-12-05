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

package cert

import (
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// generateKubeconfig 生成指定角色的 kubeconfig 文件
func generateKubeconfig(rootcaContent, certContent, keyContent []byte, apiserverEndpoint, clientName, contextName string) error {

	// 创建 kubeconfig 结构体
	kubeconfig := NewKubeconfig()

	// 设置集群信息
	kubeconfig.Clusters["kubernetes"] = &clientcmdapi.Cluster{
		Server:                   apiserverEndpoint, //todo后续从配置传入
		CertificateAuthorityData: rootcaContent,     // 如果已经有 CA 证书文件，则不需要设置这个字段
	}

	// 设置用户信息
	kubeconfig.AuthInfos[clientName] = &clientcmdapi.AuthInfo{
		ClientCertificateData: certContent,
		ClientKeyData:         keyContent,
	}

	// 设置上下文信息
	kubeconfig.Contexts[contextName] = &clientcmdapi.Context{
		Cluster:  "kubernetes",
		AuthInfo: clientName, //  context里面的那个user，和下面的uesr name保持一致
	}

	// 设置当前上下文，与前面设置的上下文name保持一致
	kubeconfig.CurrentContext = contextName

	return nil
}

// SaveKubeconfig 将 kubeconfig 结构体保存到文件
func SaveKubeconfig(config *clientcmdapi.Config, filePath string) error {
	err := clientcmd.WriteToFile(*config, filePath)
	if err != nil {
		return err
	}

	return nil
}

// NewKubeconfig 返回一个初始化好的 kubeconfig 结构体实例
func NewKubeconfig() *clientcmdapi.Config {
	return &clientcmdapi.Config{
		APIVersion:     "v1",
		Kind:           "Config",
		Clusters:       make(map[string]*clientcmdapi.Cluster),
		Contexts:       make(map[string]*clientcmdapi.Context),
		CurrentContext: "", // 这里根据需要设置默认的当前上下文
		AuthInfos:      make(map[string]*clientcmdapi.AuthInfo),
	}
}
