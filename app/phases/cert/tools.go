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
	"encoding/pem"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// PrivateKey负责生成密钥
func PrivateKey() (*rsa.PrivateKey, error) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate RSA private key")
	}

	return rsaKey, nil
}

// PrivateKeyToPem 返回私钥的PEM格式字节切片
func PrivateKeyToPem(key *rsa.PrivateKey) []byte {
	keyInBytes := x509.MarshalPKCS1PrivateKey(key)
	keyinPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyInBytes,
		},
	)
	return keyinPem
}

func PemToPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.Errorf("could not find a PEM block in the private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// CACertPEM 返回证书的PEM格式字节切片
func CertToPem(cert *x509.Certificate) []byte {
	certInPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		},
	)
	return certInPem
}

func PemToCertificate(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.Errorf("could not find a PEM block in the certificate")
	}
	return x509.ParseCertificate(block.Bytes)
}

// SaveCertificateToFile 将证书保存到文件
func (c *CertKey) SaveCertificateToFile(filename string) error {
	err := os.WriteFile(c.SavePath+"/"+filename, c.CertRaw, 0644)
	if err != nil {
		logrus.Errorf("Faile to save %s: %v", filename, err)
		return err
	}

	logrus.Infof("Successfully saved %s", filename)

	return nil
}

// SavePrivateKeyToFile 将私钥保存到文件
func (c *CertKey) SavePrivateKeyToFile(filename string) error {
	err := os.WriteFile(c.SavePath+"/"+filename, c.KeyRaw, 0600)
	if err != nil {
		logrus.Errorf("Faile to save %s: %v", filename, err)
		return err
	}

	logrus.Infof("Successfully saved %s", filename)

	return nil
}
