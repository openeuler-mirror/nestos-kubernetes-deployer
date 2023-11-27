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

func SignedCertificate(
	cfg *CertConfig,
	csr *x509.CertificateRequest,
	key *rsa.PrivateKey,
	caCert *x509.Certificate,
	caKey *rsa.PrivateKey,
) (*x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logrus.Errorf("Failed to generate serialNumber: %v", err)
		return nil, err
	}
	certTemplate := x509.Certificate{
		BasicConstraintsValid: true,
		IsCA:                  cfg.IsCA,
		DNSNames:              csr.DNSNames,
		ExtKeyUsage:           cfg.ExtKeyUsages,
		IPAddresses:           csr.IPAddresses,
		KeyUsage:              cfg.KeyUsages,
		NotAfter:              time.Now().Add(cfg.Validity),
		NotBefore:             caCert.NotBefore,
		SerialNumber:          serialNumber,
		Subject:               csr.Subject,
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, caCert, key.Public(), caKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create x509 certificate")
	}
	return x509.ParseCertificate(certBytes)
}

func GenerateSignedCertificate(caKey *rsa.PrivateKey, caCert *x509.Certificate,
	cfg *CertConfig) (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := PrivateKey()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate private key")
	}
	// create a CSR
	csrTmpl := x509.CertificateRequest{Subject: cfg.Subject, DNSNames: cfg.DNSNames, IPAddresses: cfg.IPAddresses}
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &csrTmpl, key)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create certificate request")
	}
	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		logrus.Debugf("Failed to parse x509 certificate request: %s", err)
		return nil, nil, errors.Wrap(err, "error parsing x509 certificate request")
	}

	// create a cert
	cert, err := SignedCertificate(cfg, csr, key, caCert, caKey)
	if err != nil {
		logrus.Debugf("Failed to create a signed certificate: %s", err)
		return nil, nil, errors.Wrap(err, "failed to create a signed certificate")
	}
	return key, cert, nil

}

type SignedCertKey struct {
	CertKey
}

func (c *SignedCertKey) Generate(
	cfg *CertConfig,
	parentCA CertKeyInterface,
) error {
	var key *rsa.PrivateKey
	var crt *x509.Certificate
	var err error

	caKey, err := PemToPrivateKey(parentCA.Key())
	if err != nil {
		logrus.Debugf("Failed to parse RSA private key: %s", err)
		return errors.Wrap(err, "failed to parse rsa private key")
	}

	caCert, err := PemToCertificate(parentCA.Cert())
	if err != nil {
		logrus.Debugf("Failed to parse x509 certificate: %s", err)
		return errors.Wrap(err, "failed to parse x509 certificate")
	}

	key, crt, err = GenerateSignedCertificate(caKey, caCert, cfg)
	if err != nil {
		logrus.Debugf("Failed to generate signed cert/key pair: %s", err)
		return errors.Wrap(err, "failed to generate signed cert/key pair")
	}

	c.KeyRaw = PrivateKeyToPem(key)
	c.CertRaw = CertToPem(crt)

	return nil
}
