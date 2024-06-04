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
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"nestos-kubernetes-deployer/pkg/utils"
	"os"
	"path/filepath"

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

// PublicKeyToPem 返回公钥的PEM格式字节切片
func PublicKeyToPem(key *rsa.PublicKey) ([]byte, error) {
	keyInBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal public key")
	}

	keyinPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: keyInBytes,
		},
	)
	return keyinPem, nil
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

// SaveFileToLocal 将文件保存到本地
func SaveFileToLocal(savepath string, file []byte) error {
	err := os.MkdirAll(filepath.Dir(savepath), 0755)
	if err != nil {
		logrus.Errorf("Failed to create directory: %v", err)
		return err
	}

	err = os.WriteFile(savepath, file, 0644)
	if err != nil {
		logrus.Errorf("Faile to save %s: %v", savepath, err)
		return err
	}

	// logrus.Infof("Successfully saved %s", savepath)

	return nil
}

// GenerateCACertHashes: specifies a set of public key pins to verify when token-based discovery is used.
func GenerateCACertHashes(certData []byte) (string, error) {
	block, _ := pem.Decode(certData)
	if block == nil {
		return "", fmt.Errorf("failed to decode certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse certificate: %v", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %v", err)
	}

	// Convert the hash to a hexadecimal string
	caCertHashes := fmt.Sprintf("%x", sha256.Sum256(publicKeyBytes))
	caCertHashes = "sha256:" + caCertHashes

	return caCertHashes, nil
}

func CertsToBytes(certs []utils.StorageContent) ([]byte, error) {
	jsonData, err := json.Marshal(certs)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
