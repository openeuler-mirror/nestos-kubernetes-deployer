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
	"crypto/x509"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/utils"
	"net"

	netutils "k8s.io/utils/net"

	"github.com/sirupsen/logrus"
)

type CertGenerator struct {
	ClusterID  string
	CaCertHash string
	Node       *asset.NodeAsset
}

func NewCertGenerator(clusterID string, node *asset.NodeAsset) *CertGenerator {
	return &CertGenerator{
		ClusterID: clusterID,
		Node:      node,
	}
}

// 生成所有证书文件和kubeconfig
func (cg *CertGenerator) GenerateAllFiles() error {

	var certs []utils.StorageContent
	clusterID := cg.ClusterID
	//读取配置
	clusterconfig, _ := configmanager.GetClusterConfig(clusterID)
	globalconfig, _ := configmanager.GetGlobalConfig()

	//获取node节点hostname和ip地址
	hostname := cg.Node.Hostname
	ipaddress := cg.Node.IP

	//用于后续kubeconfig生成
	apiserverEndpoint := "https://" + clusterconfig.Kubernetes.ApiServerEndpoint

	//读取用户自定义服务子网IP
	/*TODO: 1. 新增internalAPIServerVirtualIP 字段用于读取用户自定义内容；
	        2. 新增判断，默认值取用Network.Service_Subnet并进行以下解析，如用户填充internalAPIServerVirtualIP
			   则读取用户自定义内容
			3. 持续调研service clusterip相关内容，是否有统一入口进行相关配置。*/
	_, svcSubnet, err := net.ParseCIDR(clusterconfig.Network.ServiceSubnet)
	if err != nil {
		logrus.Errorf("unable to get internal Kubernetes Service IP from the given service CIDR: %v\n", err)
		return err
	}
	internalAPIServerVirtualIP, err := netutils.GetIndexedIP(svcSubnet, 1)
	if err != nil {
		logrus.Errorf("unable to get the first IP address from the given CIDR: %v\n", err)
		return err
	}

	/* **********生成root CA 证书和密钥********** */

	rootCACert, err := GenerateAllCA(clusterconfig.CertAsset.RootCaCertPath,
		clusterconfig.CertAsset.RootCaKeyPath, "kubernetes", []string{"kubernetes"})
	if err != nil {
		logrus.Errorf("Error generating root CA:%v", err)
		return err
	}

	/*如果用户没有提供自定义路径，则将ca保存在以下目录；
	  如果用户提供了自定义路径，也保存一份在以下路径，并反存到配置文件中*/
	clusterconfig.CertAsset.RootCaCertPath = globalconfig.PersistDir + "/" + clusterID + "/pki/ca.crt"
	clusterconfig.CertAsset.RootCaKeyPath = globalconfig.PersistDir + "/" + clusterID + "/pki/ca.key"

	//保存root CA证书和密钥到宿主机
	err = SaveFileToLocal(globalconfig.PersistDir+"/"+clusterID+"/pki/ca.crt", rootCACert.CertRaw)
	if err != nil {
		return err
	}

	err = SaveFileToLocal(globalconfig.PersistDir+"/"+clusterID+"/pki/ca.key", rootCACert.KeyRaw)
	if err != nil {
		return err
	}

	rootCACertContent := utils.StorageContent{
		Path:    utils.CaCrt,
		Mode:    int(utils.CertFileMode),
		Content: rootCACert.CertRaw,
	}

	rootCAKeyContent := utils.StorageContent{
		Path:    utils.CaKey,
		Mode:    int(utils.CertFileMode),
		Content: rootCACert.KeyRaw,
	}

	certs = append(certs, rootCACertContent, rootCAKeyContent)

	cg.CaCertHash, err = GenerateCACertHashes(rootCACert.CertRaw)
	if err != nil {
		logrus.Errorf("error to generate ca cert hash: %v", err)
	}

	/* **********生成etcd CA 证书和密钥********** */

	etcdCACert, err := GenerateAllCA(clusterconfig.CertAsset.EtcdCaCertPath,
		clusterconfig.CertAsset.EtcdCaKeyPath, "etcd-ca", []string{"etcd-ca"})
	if err != nil {
		logrus.Errorf("Error generating etcd CA:%v", err)
		return err
	}

	/*如果用户没有提供自定义路径，则将ca保存在以下目录；
	  如果用户提供了自定义路径，也保存一份在以下路径，并反存到配置文件中*/
	clusterconfig.CertAsset.EtcdCaCertPath = globalconfig.PersistDir + "/" + clusterID + "/pki/etcd/ca.crt"
	clusterconfig.CertAsset.EtcdCaKeyPath = globalconfig.PersistDir + "/" + clusterID + "/pki/etcd/ca.key"

	//保存etcd-ca和密钥到宿主机
	err = SaveFileToLocal(globalconfig.PersistDir+"/"+clusterID+"/pki/etcd/ca.crt", etcdCACert.CertRaw)
	if err != nil {
		return err
	}

	err = SaveFileToLocal(globalconfig.PersistDir+"/"+clusterID+"/pki/etcd/ca.key", etcdCACert.KeyRaw)
	if err != nil {
		return err
	}

	etcdCACertContent := utils.StorageContent{
		Path:    utils.EtcdCaCrt,
		Mode:    int(utils.CertFileMode),
		Content: etcdCACert.CertRaw,
	}

	etcdCAKeyContent := utils.StorageContent{
		Path:    utils.EtcdCaKey,
		Mode:    int(utils.CertFileMode),
		Content: etcdCACert.KeyRaw,
	}

	certs = append(certs, etcdCACertContent, etcdCAKeyContent)

	/* **********生成front-proxy CA 证书和密钥********** */

	frontProxyCACert, err := GenerateAllCA(clusterconfig.CertAsset.FrontProxyCaCertPath,
		clusterconfig.CertAsset.FrontProxyCaKeyPath, "front-proxy-ca", []string{"front-proxy-ca"})
	if err != nil {
		logrus.Errorf("Error generating front-proxy CA:%v", err)
		return err
	}

	/*如果用户没有提供自定义路径，则将ca保存在以下目录；
	  如果用户提供了自定义路径，也保存一份在以下路径，并反存到配置文件中*/
	clusterconfig.CertAsset.FrontProxyCaCertPath = globalconfig.PersistDir + "/" + clusterID + "/pki/front-proxy-ca.crt"
	clusterconfig.CertAsset.FrontProxyCaKeyPath = globalconfig.PersistDir + "/" + clusterID + "/pki/front-proxy-ca.key"

	//保存front-proxy-ca和密钥到宿主机
	err = SaveFileToLocal(globalconfig.PersistDir+"/"+clusterID+"/pki/front-proxy-ca.crt", frontProxyCACert.CertRaw)
	if err != nil {
		return err
	}

	err = SaveFileToLocal(globalconfig.PersistDir+"/"+clusterID+"/pki/front-proxy-ca.key", frontProxyCACert.KeyRaw)
	if err != nil {
		return err
	}

	frontProxyCACertContent := utils.StorageContent{
		Path:    utils.FrontProxyCaCrt,
		Mode:    int(utils.CertFileMode),
		Content: frontProxyCACert.CertRaw,
	}

	frontProxyCAKeyContent := utils.StorageContent{
		Path:    utils.FrontProxyCaKey,
		Mode:    int(utils.CertFileMode),
		Content: frontProxyCACert.KeyRaw,
	}

	certs = append(certs, frontProxyCACertContent, frontProxyCAKeyContent)

	/* **********生成 sa.pub和sa.key********** */

	sakeypair, err := GenerateKeyPair()
	if err != nil {
		logrus.Errorf("Error generating sa keypair:%v", err)
		return err
	}

	/*如果用户没有提供自定义路径，则将密钥对保存在以下目录；
	  如果用户提供了自定义路径，也保存一份在以下路径，并反存到配置文件中*/
	clusterconfig.CertAsset.SaKey = globalconfig.PersistDir + "/pki/sa.key"
	clusterconfig.CertAsset.SaPub = globalconfig.PersistDir + "/pki/sa.pub"

	//保存密钥对到宿主机
	err = SaveFileToLocal(globalconfig.PersistDir+"/"+clusterID+"/pki/sa.key", sakeypair.PrivateKeyPEM)
	if err != nil {
		return err
	}

	err = SaveFileToLocal(globalconfig.PersistDir+"/"+clusterID+"/pki/sa.pub", sakeypair.PublicKeyPEM)
	if err != nil {
		return err
	}

	saKeyContent := utils.StorageContent{
		Path:    utils.SaKey,
		Mode:    int(utils.CertFileMode),
		Content: sakeypair.PrivateKeyPEM,
	}

	saPubContent := utils.StorageContent{
		Path:    utils.SaPub,
		Mode:    int(utils.CertFileMode),
		Content: sakeypair.PublicKeyPEM,
	}

	certs = append(certs, saKeyContent, saPubContent)

	/* **********生成 /etcd/server.crt********** */

	commonName := hostname
	dnsNames := []string{hostname, "localhost"}
	extKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	ipAddresses := []net.IP{net.ParseIP(ipaddress), net.ParseIP("127.0.0.1")}

	servercrt, err := GenerateAllSignedCert(commonName,
		nil, dnsNames, extKeyUsage, ipAddresses, etcdCACert.CertRaw, etcdCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generating /etcd/server cert:%v", err)
		return err
	}

	serverCertContent := utils.StorageContent{
		Path:    utils.ServerCrt,
		Mode:    int(utils.CertFileMode),
		Content: servercrt.CertRaw,
	}

	serverKeyContent := utils.StorageContent{
		Path:    utils.ServerKey,
		Mode:    int(utils.CertFileMode),
		Content: servercrt.KeyRaw,
	}

	certs = append(certs, serverCertContent, serverKeyContent)

	/* **********生成 /etcd/peer.crt********** */

	commonName = hostname
	dnsNames = []string{hostname, "localhost"}
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	ipAddresses = []net.IP{net.ParseIP(ipaddress), net.ParseIP("127.0.0.1"), net.ParseIP("::1")}

	peercrt, err := GenerateAllSignedCert(commonName,
		nil, dnsNames, extKeyUsage, ipAddresses, etcdCACert.CertRaw, etcdCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generating /etcd/peer cert:%v", err)
		return err
	}

	peerCertContent := utils.StorageContent{
		Path:    utils.PeerCrt,
		Mode:    int(utils.CertFileMode),
		Content: peercrt.CertRaw,
	}

	peerKeyContent := utils.StorageContent{
		Path:    utils.PeerKey,
		Mode:    int(utils.CertFileMode),
		Content: peercrt.KeyRaw,
	}

	certs = append(certs, peerCertContent, peerKeyContent)

	/* **********生成 apiserver.crt********** */

	commonName = "kube-apiserver"
	dnsNames = []string{hostname, "kubernetes", "kubernetes.default",
		"kubernetes.default.svc", "kubernetes.default.svc.cluster", "kubernetes.default.svc.cluster.local"}
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	ipAddresses = []net.IP{net.ParseIP(ipaddress), net.ParseIP("127.0.0.1"), net.ParseIP(internalAPIServerVirtualIP.String())}

	apiservercrt, err := GenerateAllSignedCert(commonName,
		nil, dnsNames, extKeyUsage, ipAddresses, rootCACert.CertRaw, rootCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generating apiserver cert:%v", err)
		return err
	}

	apiserverCertContent := utils.StorageContent{
		Path:    utils.ApiserverCrt,
		Mode:    int(utils.CertFileMode),
		Content: apiservercrt.CertRaw,
	}

	apiserverKeyContent := utils.StorageContent{
		Path:    utils.ApiserverKey,
		Mode:    int(utils.CertFileMode),
		Content: apiservercrt.KeyRaw,
	}

	certs = append(certs, apiserverCertContent, apiserverKeyContent)

	/* **********生成 front-proxy-client.crt********** */

	commonName = "front-proxy-client"
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	frontProxyClientcrt, err := GenerateAllSignedCert(commonName,
		nil, nil, extKeyUsage, nil, frontProxyCACert.CertRaw, frontProxyCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generating front-proxy-client cert:%v", err)
		return err
	}

	frontProxyClientCertContent := utils.StorageContent{
		Path:    utils.FrontProxyClientCrt,
		Mode:    int(utils.CertFileMode),
		Content: frontProxyClientcrt.CertRaw,
	}

	frontProxyClientKeyContent := utils.StorageContent{
		Path:    utils.FrontProxyClientKey,
		Mode:    int(utils.CertFileMode),
		Content: frontProxyClientcrt.KeyRaw,
	}

	certs = append(certs, frontProxyClientCertContent, frontProxyClientKeyContent)

	/* **********生成 apiserver-kubelet-client.crt********** */

	commonName = "kube-apiserver-kubelet-client"
	organization := []string{"system:masters"}
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	apiserverKubeletClientcrt, err := GenerateAllSignedCert(commonName,
		organization, nil, extKeyUsage, nil, rootCACert.CertRaw, rootCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generating apiserver-kubelet-client cert:%v", err)
		return err
	}

	apiserverKubeletClientCertContent := utils.StorageContent{
		Path:    utils.ApiserverKubeletClientCrt,
		Mode:    int(utils.CertFileMode),
		Content: apiserverKubeletClientcrt.CertRaw,
	}

	apiserverKubeletClientKeyContent := utils.StorageContent{
		Path:    utils.ApiserverKubeletClientKey,
		Mode:    int(utils.CertFileMode),
		Content: apiserverKubeletClientcrt.KeyRaw,
	}

	certs = append(certs, apiserverKubeletClientCertContent, apiserverKubeletClientKeyContent)

	/* **********生成 apiserver-etcd-client.crt********** */

	commonName = "kube-apiserver-etcd-client"
	organization = []string{"system:masters"}
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	apiserverEtcdClient, err := GenerateAllSignedCert(commonName,
		organization, nil, extKeyUsage, nil, etcdCACert.CertRaw, etcdCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generating kube-apiserver-etcd-client cert:%v", err)
		return err
	}

	apiserverEtcdClientCertContent := utils.StorageContent{
		Path:    utils.ApiserverEtcdClientCrt,
		Mode:    int(utils.CertFileMode),
		Content: apiserverEtcdClient.CertRaw,
	}

	apiserverEtcdClientKeyContent := utils.StorageContent{
		Path:    utils.ApiserverEtcdClientKey,
		Mode:    int(utils.CertFileMode),
		Content: apiserverEtcdClient.KeyRaw,
	}

	certs = append(certs, apiserverEtcdClientCertContent, apiserverEtcdClientKeyContent)

	/* **********生成 healthcheck.crt********** */

	commonName = "kube-etcd-healthcheck-client"
	organization = []string{"system:masters"}
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	healthcheckcrt, err := GenerateAllSignedCert(commonName,
		organization, nil, extKeyUsage, nil, etcdCACert.CertRaw, etcdCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generating healthcheck cert:%v", err)
		return err
	}

	healthcheckCertContent := utils.StorageContent{
		Path:    utils.HealthcheckClientCrt,
		Mode:    int(utils.CertFileMode),
		Content: healthcheckcrt.CertRaw,
	}

	healthcheckKeyContent := utils.StorageContent{
		Path:    utils.HealthcheckClientKey,
		Mode:    int(utils.CertFileMode),
		Content: healthcheckcrt.KeyRaw,
	}

	certs = append(certs, healthcheckCertContent, healthcheckKeyContent)

	/* **********生成 admin.config********** */

	commonName = "kubernetes-admin"
	organization = []string{"system:masters"}
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	admincrt, err := GenerateAllSignedCert(commonName,
		organization, nil, extKeyUsage, nil, rootCACert.CertRaw, rootCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generate admin cert:%v", err)
		return err
	}

	adminKubeconfig, err := generateKubeconfig(rootCACert.CertRaw, admincrt.CertRaw, admincrt.KeyRaw,
		apiserverEndpoint, "kubernetes-admin", "kubernetes-admin@kubernetes")
	if err != nil {
		logrus.Errorf("Error generate admin.config:%v", err)
		return err
	}

	clusterconfig.Kubernetes.AdminKubeConfig = globalconfig.PersistDir + "/" + clusterID + "/admin.config"

	//将admin.config文件保存至宿主机
	err = SaveFileToLocal(globalconfig.PersistDir+"/"+clusterID+"/admin.config", adminKubeconfig)
	if err != nil {
		return err
	}

	adminKubeconfigContent := utils.StorageContent{
		Path:    utils.AdminConfig,
		Mode:    int(utils.CertFileMode),
		Content: adminKubeconfig,
	}

	certs = append(certs, adminKubeconfigContent)

	/* **********生成 controller-manager.config********** */

	commonName = "system:kube-controller-manager"
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	controllerManagercrt, err := GenerateAllSignedCert(commonName,
		nil, nil, extKeyUsage, nil, rootCACert.CertRaw, rootCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generate controller-manager cert:%v", err)
		return err
	}

	controllerManagerKubeconfig, err := generateKubeconfig(rootCACert.CertRaw, controllerManagercrt.CertRaw, controllerManagercrt.KeyRaw,
		apiserverEndpoint, "system:kube-controller-manager", "system:kube-controller-manager@kubernetes")
	if err != nil {
		logrus.Errorf("Error generate controller-manager.config:%v", err)
		return err
	}

	controllerManagerKubeconfigContent := utils.StorageContent{
		Path:    utils.ControllerManager,
		Mode:    int(utils.CertFileMode),
		Content: controllerManagerKubeconfig,
	}

	certs = append(certs, controllerManagerKubeconfigContent)

	/* **********生成 scheduler.config********** */

	commonName = "system:kube-scheduler"
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	schedulercrt, err := GenerateAllSignedCert(commonName,
		nil, nil, extKeyUsage, nil, rootCACert.CertRaw, rootCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generate scheduler cert:%v", err)
		return err
	}

	schedulerKubeconfig, err := generateKubeconfig(rootCACert.CertRaw, schedulercrt.CertRaw, schedulercrt.KeyRaw,
		apiserverEndpoint, "system:kube-scheduler", "system:kube-scheduler@kubernetes")
	if err != nil {
		logrus.Errorf("Error generate scheduler.config:%v", err)
		return err
	}

	schedulerKubeconfigContent := utils.StorageContent{
		Path:    utils.SchedulerConf,
		Mode:    int(utils.CertFileMode),
		Content: schedulerKubeconfig,
	}

	certs = append(certs, schedulerKubeconfigContent)

	/* **********生成 kubelet.config********** */

	commonName = "system:node:" + hostname
	organization = []string{"system:nodes"}
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	kubeletcrt, err := GenerateAllSignedCert(commonName,
		organization, nil, extKeyUsage, nil, rootCACert.CertRaw, rootCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generate kubelet cert:%v", err)
		return err
	}

	kubeletKubeconfig, err := generateKubeconfig(rootCACert.CertRaw, kubeletcrt.CertRaw, kubeletcrt.KeyRaw,
		apiserverEndpoint, "system:node:"+hostname, "system:node:"+hostname+"@kubernetes")
	if err != nil {
		logrus.Errorf("Error generate kubelet.config:%v", err)
		return err
	}

	kubeletKubeconfigContent := utils.StorageContent{
		Path:    utils.KubeletConfig,
		Mode:    int(utils.CertFileMode),
		Content: kubeletKubeconfig,
	}

	certs = append(certs, kubeletKubeconfigContent)

	cg.Node.Certs = certs

	return nil
}
