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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"io/ioutil"

	"github.com/pkg/errors"
)

// SetUserCA 读取用户提供的ca证书和密钥路径
func setUserCA(a *SelfSignedCertKey, certPath, keyPath string) error {
	cacert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return err
	}

	cakey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return err
	}

	// 存储到结构体中
	a.CertRaw = cacert
	a.KeyRaw = cakey

	return nil
}

//GenerateRootCA()用于生成/etc/kubernetes/pki/ca.crt
func GenerateRootCA() (*SelfSignedCertKey, error) {

	a := SelfSignedCertKey{}
	//接受用户输入的两个路径
	userCACertPath, userCAKeyPath := GetCustomRootCAPathFromConfig()

	// 如果用户提供了路径，则读取用户提供的证书和密钥
	if userCACertPath != "" && userCAKeyPath != "" {
		err := setUserCA(&a, userCACertPath, userCAKeyPath)
		if err != nil {
			return nil, err
		}
	} else {
		// 如果用户没有提供自定义CA证书路径，则继续生成
		cfg := &CertConfig{
			Subject:   pkix.Name{CommonName: "Kubernetes", OrganizationalUnit: []string{"NestOS"}},
			KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			Validity:  3650,
			IsCA:      true,
		}

		err := a.Generate(cfg)
		if err != nil {
			return nil, err
		}
	}

	return &a, nil
}

// GetCustomRootCAPathFromConfig 实现从配置文件中获取用户提供的自定义root CA路径和对应密钥路径的逻辑
func GetCustomRootCAPathFromConfig() (string, string) {
	// TODO: 从配置文件中获取用户提供的自定义 CA 证书路径和对应密钥路径
	// 如果用户没有提供路径，返回空字符串
	return "", ""
}

//GenerateEtcdCA()用于生成/etc/kubernetes/pki/etcd/ca.crt
func GenerateEtcdCA() (*SelfSignedCertKey, error) {

	a := SelfSignedCertKey{}
	//接受用户输入的两个路径
	userCACertPath, userCAKeyPath := GetCustomEtcdCAPathFromConfig()

	// 如果用户提供了路径，则读取用户提供的证书和密钥
	if userCACertPath != "" && userCAKeyPath != "" {
		err := setUserCA(&a, userCACertPath, userCAKeyPath)
		if err != nil {
			return nil, err
		}
	} else {
		// 如果用户没有提供自定义CA证书路径，则继续生成
		cfg := &CertConfig{
			Subject:   pkix.Name{CommonName: "etcd-ca", OrganizationalUnit: []string{"NestOS"}},
			KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			Validity:  3650,
			IsCA:      true,
		}

		err := a.Generate(cfg)
		if err != nil {
			return nil, err
		}
	}

	return &a, nil
}

// GetCustomEtcdCAPathFromConfig 实现从配置文件中获取用户提供的自定义ETCD CA路径和对应密钥路径的逻辑
func GetCustomEtcdCAPathFromConfig() (string, string) {
	// TODO: 从配置文件中获取用户提供的自定义 CA 证书路径和对应密钥路径
	// 如果用户没有提供路径，返回空字符串
	return "", ""
}

//GenerateFrontProxyCA()用于生成/etc/kubernetes/pki/front-proxy-ca.crt
func GenerateFrontProxyCA() (*SelfSignedCertKey, error) {

	a := SelfSignedCertKey{}
	//接受用户输入的两个路径
	userCACertPath, userCAKeyPath := GetCustomFrontProxyCAPathFromConfig()

	// 如果用户提供了路径，则读取用户提供的证书和密钥
	if userCACertPath != "" && userCAKeyPath != "" {
		err := setUserCA(&a, userCACertPath, userCAKeyPath)
		if err != nil {
			return nil, err
		}
	} else {
		// 如果用户没有提供自定义CA证书路径，则继续生成
		cfg := &CertConfig{
			Subject:   pkix.Name{CommonName: "front-proxy-ca", OrganizationalUnit: []string{"NestOS"}},
			KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			Validity:  3650,
			IsCA:      true,
		}

		err := a.Generate(cfg)
		if err != nil {
			return nil, err
		}
	}

	return &a, nil
}

// GetCustomFrontProxyCAPathFromConfig 实现从配置文件中获取用户提供的自定义front-proxy CA路径和对应密钥路径的逻辑
func GetCustomFrontProxyCAPathFromConfig() (string, string) {
	// TODO: 从配置文件中获取用户提供的自定义 CA 证书路径和对应密钥路径
	// 如果用户没有提供路径，返回空字符串
	return "", ""
}

/* GenerateKeyPair()用于生成/etc/kubernetes/pki/sa.pub和/etc/kubernetes/pki/sa.key ，
通常创建自定义证书时说生成四组 CA-Key与CA-Cert，其中一组就是指这个密钥对*/
func GenerateKeyPair() (*KeyPairPEM, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate RSA private key")
	}

	privateKeyPEM := PrivateKeyToPem(privateKey)
	publicKeyPEM, err := PublicKeyToPem(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return &KeyPairPEM{
		PrivateKeyPEM: privateKeyPEM,
		PublicKeyPEM:  publicKeyPEM,
	}, nil
}
