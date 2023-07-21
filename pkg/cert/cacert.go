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
)

type RootCA struct {
	SelfSignedCertKey
}

func (c *RootCA) Generate() error {
	cfg := &CertConfig{
		Subject:   pkix.Name{CommonName: "rootca", OrganizationalUnit: []string{"NestOS"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  3650,
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(cfg, "rootca")
}
