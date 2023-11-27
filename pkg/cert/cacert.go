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
	"io/ioutil"
)

// SetUserCA 读取用户提供的证书和密钥路径
func SetUserCA(a *SelfSignedCertKey, certPath, keyPath string) error {
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

func GenerateRootCA() (*SelfSignedCertKey, error) {

	a := SelfSignedCertKey{}
	//接受用户输入的两个路径
	userCACertPath, userCAKeyPath := GetCustomCAPathFromConfig()

	// 如果用户提供了路径，则设置证书和密钥
	if userCACertPath != "" && userCAKeyPath != "" {
		err := SetUserCA(&a, userCACertPath, userCAKeyPath)
		if err != nil {
			return nil, err
		}
	} else {
		// 如果用户没有提供路径，则继续生成
		cfg := &CertConfig{
			Subject:   pkix.Name{CommonName: "Kubernetes", OrganizationalUnit: []string{"NestOS"}},
			KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			Validity:  3650,
			IsCA:      true,
		}

		err := a.Generate(cfg)
		if err != nil {
			return nil, err
		}
	}

	return &a, nil
}

// GetCustomCAPathFromConfig 实现从配置文件中获取用户提供的自定义 CA和对应密钥路径 证书路径的逻辑
func GetCustomCAPathFromConfig() (string, string) {
	// TODO: 从配置文件中获取用户提供的自定义 CA 证书路径和对应密钥路径
	// 如果用户没有提供路径，返回空字符串
	return "", ""
}
