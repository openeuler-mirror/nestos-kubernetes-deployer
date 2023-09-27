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
	"crypto/x509/pkix"
	"net"
	"time"
)

type CertInterface interface {
	// Cert returns the certificate.
	Cert() []byte
}

// CertKeyInterface contains a private key and the associated cert.
type CertKeyInterface interface {
	CertInterface
	// Key returns the private key.
	Key() []byte
}

type CertificateGenerator interface {
	GenerateCACertificate() error
	GenerateSignedCertificate(commonName string) error
}

//  CertKey 包含证书和私钥
type CertKey struct {
	CertRaw  []byte
	KeyRaw   []byte
	SavePath string
}

type CertConfig struct {
	DNSNames     []string
	ExtKeyUsages []x509.ExtKeyUsage
	IPAddresses  []net.IP
	KeyUsages    x509.KeyUsage
	Subject      pkix.Name
	Validity     time.Duration
	IsCA         bool
}
