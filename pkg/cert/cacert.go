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

func (cm *CertificateManager) GenerateCACertificate() error {

	// 生成CA的私钥和公钥
	var err error
	cm.CAKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	caPublicKey := &cm.CACert.PublicKey

	// 生成一个介于 0 和 2^128 - 1 之间的随机序列号，并将结果存储在 serialNumber 变量中
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}
	now := time.Now()

	// 设置生成证书的参数，构建CA证书模板
	caTemplate := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"NKD"},
			CommonName:   "CA",
		}, //这里还可以加很多参数信息
		NotBefore:             now,
		NotAfter:              now.AddDate(10, 0, 0),                                                      // 有效期为10年
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,               // openssl 中的 keyUsage 字段
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}, // openssl 中的 extendedKeyUsage = clientAuth, serverAuth 字段
		BasicConstraintsValid: true,
		IsCA:                  true, //表示用于CA
	}

	//caCertBytes是生成证书的中间步骤，它用于将证书的二进制表示存储在内存中，以便后续操作可以使用它
	caCertBytes, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, caPublicKey, cm.CAKey)
	if err != nil {
		return err
	}
	//cm.CACert表示已经生成的CA证书，用于存储CA证书的详细信息，例如证书序列号、主题、有效期等
	cm.CACert, err = x509.ParseCertificate(caCertBytes)
	if err != nil {
		return err
	}

	return nil

}
