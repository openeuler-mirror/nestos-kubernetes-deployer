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
	"encoding/pem"
	"io/ioutil"
)

const certpath = "/etc/kubernetes/pki/"

// CACertPEM 返回CA证书的PEM格式字节切片
func (cm *CertificateManager) CACertPEM() []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cm.CACert.Raw,
	})
}

// CAKeyPEM 返回CA私钥的PEM格式字节切片
func (cm *CertificateManager) CAKeyPEM() []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(cm.CAKey),
	})
}

// ComponentCertPEM 返回组件证书的PEM格式字节切片
func (cm *CertificateManager) ComponentCertPEM() []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cm.ComponentCert.Raw,
	})
}

// ComponentKeyPEM 返回组件私钥的PEM格式字节切片
func (cm *CertificateManager) ComponentKeyPEM() []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(cm.ComponentKey),
	})
}

// SaveCertificateToFile 将证书保存到文件
func SaveCertificateToFile(filename string, certPEM []byte) error {
	return ioutil.WriteFile(certpath+filename, certPEM, 0644)
}

// SavePrivateKeyToFile 将私钥保存到文件
func SavePrivateKeyToFile(filename string, keyPEM []byte) error {
	return ioutil.WriteFile(certpath+filename, keyPEM, 0600)
}
