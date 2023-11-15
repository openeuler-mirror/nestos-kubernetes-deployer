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
	"math/big"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// SelfSignedCertificate 只负责创建自签名的证书，这里只传入cfg和privatekey
func SelfSignedCertificate(cfg *CertConfig, key *rsa.PrivateKey) (*x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logrus.Errorf("Failed to generate serialNumber: %v", err)
		return nil, err
	}
	certTemplate := x509.Certificate{
		BasicConstraintsValid: true,
		IsCA:                  cfg.IsCA,
		KeyUsage:              cfg.KeyUsages,
		NotAfter:              time.Now().Add(cfg.Validity),
		NotBefore:             time.Now(),
		SerialNumber:          serialNumber,
		Subject:               cfg.Subject,
	}
	// 判断subject字段中CommonName和OrganizationalUnit是否为空
	if len(cfg.Subject.CommonName) == 0 || len(cfg.Subject.OrganizationalUnit) == 0 {
		return nil, errors.Errorf("certification's subject is not set, or invalid")
	}

	//certBytes是生成证书的中间步骤，它用于将证书的二进制表示存储在内存中，以便后续操作可以使用它
	certBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, key.Public(), key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create certificate")
	}

	return x509.ParseCertificate(certBytes)
}

/*GenerateSelfSignedCertificate负责根据cfg生成私钥和证书
  在这一步会调用PrivateKey生成私钥并调用SelfSignedCertificate生成证书，并将结果返回*/
func GenerateSelfSignedCertificate(cfg *CertConfig) (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := PrivateKey()
	if err != nil {
		logrus.Debugf("Failed to generate private key: %s", err)
		return nil, nil, errors.Wrap(err, "Failed to generate private key")
	}

	//这里的crt是parse之后的，表示已经生成的CA证书，用于存储CA证书的详细信息，例如证书序列号、主题、有效期等
	crt, err := SelfSignedCertificate(cfg, key)
	if err != nil {
		logrus.Debugf("Failed to create self-signed certificate: %s", err)
		return nil, nil, errors.Wrap(err, "failed to create self-signed certificate")
	}
	return key, crt, nil
}

type SelfSignedCertKey struct {
	CertKey
}

//自签名证书生成器，封装后该方法用于所有自签名的证书,并将证书和私钥转换格式后保存
func (c *SelfSignedCertKey) Generate(cfg *CertConfig, filename string) error {

	c.CertKey.SavePath = "/tmp"
	key, crt, err := GenerateSelfSignedCertificate(cfg)
	if err != nil {
		return errors.Wrap(err, "Failed to generate self-signed cert/key pair")
	}

	c.KeyRaw = PrivateKeyToPem(key)
	c.CertRaw = CertToPem(crt)

	err = c.SaveCertificateToFile(filename)
	if err != nil {
		logrus.Errorf("Faile to save %s: %v", filename, err)
	}

	return nil

}
