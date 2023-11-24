package cert

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// GenerateAllKubeconfigs 生成所有 kubeconfig 文件
func GenerateAllKubeconfigs(caPath, baseDir string) error {
	// 生成 admin kubeconfig
	if err := generateKubeconfig(caPath, "admin", baseDir, "kubernetes-admin", "kubernetes-admin@kubernetes"); err != nil {
		return err
	}

	// 生成 kube-scheduler kubeconfig
	if err := generateKubeconfig(caPath, "kube-scheduler", baseDir, "system:kube-scheduler", "system:kube-scheduler@kubernetes"); err != nil {
		return err
	}

	// 生成 kubelet kubeconfig
	if err := generateKubeconfig(caPath, "kubelet", baseDir, "system:kubelet", "system:kubelet@kubernetes"); err != nil {
		return err
	}

	// 生成 controller-manager kubeconfig
	if err := generateKubeconfig(caPath, "controller-manager", baseDir, "system:kube-controller-manager", "system:kube-controller-manager@kubernetes"); err != nil {
		return err
	}

	return nil
}

// generateKubeconfig 生成指定角色的 kubeconfig 文件
func generateKubeconfig(caPath, role, baseDir, clientName, contextName string) error {
	kubeconfigPath := fmt.Sprintf("%s/%s.conf", baseDir, role) //返回一个格式化后的字符串,即ca证书路径

	// 创建 kubeconfig 结构体
	kubeconfig := NewKubeconfig()

	// 设置集群信息
	kubeconfig.Clusters["kubernetes"] = &clientcmdapi.Cluster{
		Server:                   "https://api-server-url", //todo后续从配置传入
		CertificateAuthority:     caPath,                   //传入ca证书路径
		CertificateAuthorityData: nil,                      // 如果已经有 CA 证书文件，则不需要设置这个字段
	}

	// 设置用户信息
	kubeconfig.AuthInfos[clientName] = &clientcmdapi.AuthInfo{
		ClientCertificate: fmt.Sprintf("%s/%s-client.crt", baseDir, role),
		ClientKey:         fmt.Sprintf("%s/%s-client.key", baseDir, role),
	}

	// 设置上下文信息
	kubeconfig.Contexts[contextName] = &clientcmdapi.Context{
		Cluster:  "kubernetes",
		AuthInfo: clientName, //  context里面的那个user，和下面的uesr name保持一致
	}

	// 设置当前上下文，与前面设置的上下文name保持一致
	kubeconfig.CurrentContext = contextName

	// 保存 kubeconfig 到文件
	err := SaveKubeconfig(kubeconfig, kubeconfigPath)
	if err != nil {
		return err
	}

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
