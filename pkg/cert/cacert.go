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
	"fmt"
	"os"
)

type RootCA struct {
	SelfSignedCertKey
}

func (c *RootCA) Generate() error {
	// 检查用户是否提供了自定义的 CA 证书路径
	userCAPath := "/tmp/ca.crt" // 默认路径
	if userProvidedCAPath := GetCustomCAPathFromConfig(); userProvidedCAPath != "" {
		userCAPath = userProvidedCAPath
	}

	// 检查 CA 证书文件是否已存在
	if _, err := os.Stat(userCAPath); err == nil {
		fmt.Printf("CA 证书已存在于路径：%s。跳过生成过程。\n", userCAPath)
		return nil
	}

	// 如果 CA 证书不存在，则继续生成
	cfg := &CertConfig{
		Subject:   pkix.Name{CommonName: "rootca", OrganizationalUnit: []string{"NestOS"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  3650,
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(cfg, userCAPath)
}

// GetCustomCAPathFromConfig 实现从配置文件中获取用户提供的自定义 CA 证书路径的逻辑
func GetCustomCAPathFromConfig() string {
	// TODO: 从配置文件中获取用户提供的自定义 CA 证书路径
	// 如果用户没有提供路径，返回空字符串
	return ""
}
