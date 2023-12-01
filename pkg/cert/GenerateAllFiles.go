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
	"nestos-kubernetes-deployer/pkg/utils"
	"net"

	"github.com/sirupsen/logrus"
)

// GenerateAllFiles 生成所有证书、密钥、kubeconfig
func GenerateAllFiles(clusterID string) ([]utils.StorageContent, error) {

	var certs []utils.StorageContent

	//读取配置
	clusterconfig, _ := configmanager.GetClusterConfig(clusterID)
	globalconfig, _ := configmanager.GetGlobalConfig()

	//获取master0节点hostname和ip地址
	master0name := clusterconfig.Master.NodeAsset[0].Hostname
	master0ip := clusterconfig.Master.NodeAsset[0].IP

	/* **********生成root CA 证书和密钥********** */

	rootCACert, err := GenerateAllCA(clusterconfig.CertAsset.RootCaCertPath,
		clusterconfig.CertAsset.RootCaKeyPath, "kubernetes", []string{"kubernetes"})
	if err != nil {
		logrus.Errorf("Error generating root CA:%v", err)
		return nil, err
	}

	/*如果用户没有提供自定义路径，则将ca保存在以下目录；
	  如果用户提供了自定义路径，也保存一份在以下路径，并反存到配置文件中*/
	clusterconfig.CertAsset.RootCaCertPath = globalconfig.PersistDir + "/pki/ca.crt"
	clusterconfig.CertAsset.RootCaKeyPath = globalconfig.PersistDir + "/pki/ca.key"

	//保存root CA证书和密钥到宿主机
	err = SaveFileToLocal(globalconfig.PersistDir+"/pki/ca.crt", rootCACert.CertRaw)
	if err != nil {
		return nil, err
	}

	err = SaveFileToLocal(globalconfig.PersistDir+"/pki/ca.key", rootCACert.CertRaw)
	if err != nil {
		return nil, err
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

	/* **********生成etcd CA 证书和密钥********** */

	etcdCACert, err := GenerateAllCA(clusterconfig.CertAsset.EtcdCaCertPath,
		clusterconfig.CertAsset.EtcdCaKeyPath, "etcd-ca", []string{"etcd-ca"})
	if err != nil {
		logrus.Errorf("Error generating etcd CA:%v", err)
		return nil, err
	}

	/*如果用户没有提供自定义路径，则将ca保存在以下目录；
	  如果用户提供了自定义路径，也保存一份在以下路径，并反存到配置文件中*/
	clusterconfig.CertAsset.EtcdCaCertPath = globalconfig.PersistDir + "/pki/etcd/ca.crt"
	clusterconfig.CertAsset.EtcdCaKeyPath = globalconfig.PersistDir + "/pki/etcd/ca.key"

	//保存etcd-ca和密钥到宿主机
	err = SaveFileToLocal(globalconfig.PersistDir+"/pki/etcd/ca.crt", etcdCACert.CertRaw)
	if err != nil {
		return nil, err
	}

	err = SaveFileToLocal(globalconfig.PersistDir+"/pki/etcd/ca.key", etcdCACert.CertRaw)
	if err != nil {
		return nil, err
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

	frontProxyCACert, err := GenerateAllCA(clusterconfig.CertAsset.EtcdCaCertPath,
		clusterconfig.CertAsset.EtcdCaKeyPath, "front-proxy-ca", []string{"front-proxy-ca"})
	if err != nil {
		logrus.Errorf("Error generating front-proxy CA:%v", err)
		return nil, err
	}

	/*如果用户没有提供自定义路径，则将ca保存在以下目录；
	  如果用户提供了自定义路径，也保存一份在以下路径，并反存到配置文件中*/
	clusterconfig.CertAsset.FrontProxyCaCertPath = globalconfig.PersistDir + "/pki/front-proxy-ca.crt"
	clusterconfig.CertAsset.FrontProxyCaKeyPath = globalconfig.PersistDir + "/pki/front-proxy-ca.key"

	//保存front-proxy-ca和密钥到宿主机
	err = SaveFileToLocal(globalconfig.PersistDir+"/pki/front-proxy-ca.crt", frontProxyCACert.CertRaw)
	if err != nil {
		return nil, err
	}

	err = SaveFileToLocal(globalconfig.PersistDir+"/pki/front-proxy-ca.key", frontProxyCACert.CertRaw)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	/*如果用户没有提供自定义路径，则将密钥对保存在以下目录；
	  如果用户提供了自定义路径，也保存一份在以下路径，并反存到配置文件中*/
	clusterconfig.CertAsset.SaKey = globalconfig.PersistDir + "/pki/sa.key"
	clusterconfig.CertAsset.SaPub = globalconfig.PersistDir + "/pki/sa.pub"

	//保存密钥对到宿主机
	err = SaveFileToLocal(globalconfig.PersistDir+"/pki/sa.key", sakeypair.PrivateKeyPEM)
	if err != nil {
		return nil, err
	}

	err = SaveFileToLocal(globalconfig.PersistDir+"/pki/sa.pub", sakeypair.PublicKeyPEM)
	if err != nil {
		return nil, err
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

	/* **********生成 apiserver-etcd-client.crt********** */

	commonName := "kube-apiserver-etcd-client"
	organization := []string{"system:masters"}
	extKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	apiserverEtcdClient, err := GenerateAllSignedCert(commonName,
		organization, nil, extKeyUsage, nil, etcdCACert.CertRaw, etcdCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generating kube-apiserver-etcd-client cert:%v", err)
		return nil, err
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

	/* **********生成 /etcd/server.crt********** */

	commonName = master0name
	dnsNames := []string{master0name, "localhost"}
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	ipAddresses := []net.IP{net.ParseIP(master0ip), net.ParseIP("127.0.0.1")}

	servercrt, err := GenerateAllSignedCert(commonName,
		nil, dnsNames, extKeyUsage, ipAddresses, etcdCACert.CertRaw, etcdCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generating /etcd/server cert:%v", err)
		return nil, err
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

	commonName = master0name
	dnsNames = []string{master0name, "localhost"}
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	ipAddresses = []net.IP{net.ParseIP(master0ip), net.ParseIP("127.0.0.1"), net.ParseIP("::1")}

	peercrt, err := GenerateAllSignedCert(commonName,
		nil, dnsNames, extKeyUsage, ipAddresses, etcdCACert.CertRaw, etcdCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generating /etcd/peer cert:%v", err)
		return nil, err
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

	/* **********生成 healthcheck.crt********** */

	commonName = "kube-etcd-healthcheck-client"
	organization = []string{"system:masters"}
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	healthcheckcrt, err := GenerateAllSignedCert(commonName,
		organization, nil, extKeyUsage, nil, etcdCACert.CertRaw, etcdCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generating healthcheck cert:%v", err)
		return nil, err
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

	/* **********生成 front-proxy-client.crt********** */

	commonName = "front-proxy-client"
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	frontProxyClientcrt, err := GenerateAllSignedCert(commonName,
		nil, nil, extKeyUsage, nil, frontProxyCACert.CertRaw, frontProxyCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generating front-proxy-client cert:%v", err)
		return nil, err
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

	/* **********生成 apiserver.crt********** */

	commonName = "kube-apiserver"
	dnsNames = []string{master0name, "kubernetes", "kubernetes.default",
		"kubernetes.default.svc", "kubernetes.default.svc.cluster", "kubernetes.default.svc.cluster.local"}
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	ipAddresses = []net.IP{net.ParseIP(master0ip), net.ParseIP("127.0.0.1"), net.ParseIP("0.0.0.0")}

	apiservercrt, err := GenerateAllSignedCert(commonName,
		nil, dnsNames, extKeyUsage, ipAddresses, rootCACert.CertRaw, rootCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generating apiserver cert:%v", err)
		return nil, err
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

	/* **********生成 apiserver-kubelet-client.crt********** */

	commonName = "kube-apiserver-kubelet-client"
	organization = []string{"system:masters"}
	extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	apiserverKubeletClientcrt, err := GenerateAllSignedCert(commonName,
		organization, nil, extKeyUsage, nil, rootCACert.CertRaw, rootCACert.KeyRaw)
	if err != nil {
		logrus.Errorf("Error generating apiserver-kubelet-client cert:%v", err)
		return nil, err
	}

	apiserverKubeletClientCertContent := utils.StorageContent{
		Path:    utils.ApiserverKubeletClientCrt,
		Mode:    int(utils.CertFileMode),
		Content: apiserverKubeletClientcrt.CertRaw,
	}

	apiserverKubeletClientKeyContent := utils.StorageContent{
		Path:    utils.ApiserverKubeletClientKey,
		Mode:    int(utils.CertFileMode),
		Content: frontProxyClientcrt.KeyRaw,
	}

	certs = append(certs, apiserverKubeletClientCertContent, apiserverKubeletClientKeyContent)

	return certs, nil
}
