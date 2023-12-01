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
	"net"

	"github.com/pkg/errors"
)

// SetUserCA 读取用户提供的各类ca证书和密钥路径中的内容
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

/*
GenerateAllCA()用于生成
/etc/kubernetes/pki/ca.crt
/etc/kubernetes/pki/front-proxy-ca.crt
/etc/kubernetes/pki/etcd/ca.crt   以及所有对应key
*/
func GenerateAllCA(userCACertPath, userCAKeyPath, commonname string, dnsname []string) (*SelfSignedCertKey, error) {

	a := SelfSignedCertKey{}

	// 如果用户提供了路径，则读取用户提供的证书和密钥
	if userCACertPath != "" && userCAKeyPath != "" {
		err := setUserCA(&a, userCACertPath, userCAKeyPath)
		if err != nil {
			return nil, err
		}
	} else {
		// 如果用户没有提供自定义CA证书路径，则继续生成
		cfg := &CertConfig{
			Subject:   pkix.Name{CommonName: commonname},
			KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			Validity:  3650,
			IsCA:      true,
			DNSNames:  dnsname,
		}

		err := a.Generate(cfg)
		if err != nil {
			return nil, err
		}
	}

	return &a, nil
}

/* GenerateKeyPair()用于生成/etc/kubernetes/pki/sa.pub和/etc/kubernetes/pki/sa.key ，
通常创建自定义证书时说生成四组 CA-Key与CA-Cert其中一组就是指这个密钥对*/
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

//GenerateAllSignedCert()用于生成所有签发的证书
func GenerateAllSignedCert(commonname string, org, dnsname []string, extkeyusage []x509.ExtKeyUsage,
	ip []net.IP, cacert, cakey []byte) (*SignedCertKey, error) {
	a := SignedCertKey{}

	cfg := &CertConfig{
		Subject:      pkix.Name{CommonName: commonname, Organization: org},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: extkeyusage,
		Validity:     3650,
		IsCA:         false,
		DNSNames:     dnsname,
		IPAddresses:  ip,
	}
	err := a.Generate(cfg, cacert, cakey)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
