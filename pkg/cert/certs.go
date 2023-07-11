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
	"math/big"
	"time"
)

// 使用CA证书和私钥生成组件证书和私钥
func (cm *CertificateManager) GenerateComponentCertificate(componentName string) error {

	// 创建组件的公钥和私钥
	var err error
	cm.ComponentKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	componentPublicKey := &cm.ComponentKey.PublicKey

	// 生成一个介于 0 和 2^128 - 1 之间的随机序列号，并将结果存储在 serialNumber 变量中
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}
	now := time.Now()

	// 组件证书模板
	componentTemplate := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{Organization: []string{"NKD"}, CommonName: componentName},
		NotBefore:    now,
		NotAfter:     time.Now().AddDate(0, 0, cm.ValidDays), // 有效期为1年
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	// 使用CA证书和私钥生成组件证书
	componentCertBytes, err := x509.CreateCertificate(rand.Reader, componentTemplate, cm.CACert, componentPublicKey, cm.CAKey)
	if err != nil {
		return err
	}

	cm.ComponentCert, err = x509.ParseCertificate(componentCertBytes)
	if err != nil {
		return err
	}

	return nil
}
