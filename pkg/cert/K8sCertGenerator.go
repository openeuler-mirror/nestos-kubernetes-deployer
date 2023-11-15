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

// 用于创建apiserver.crt,是kube-apiserver 对外提供服务的服务器证书及私钥

func GenerateApiServer() error {
	a := SignedCertKey{}

	ca := &RootCA{} //仍需在资产管理模块完善，未来可以直接调用

	cfg := &CertConfig{
		Subject:      pkix.Name{CommonName: "apiserver server", Organization: []string{"NestOS"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     3650,
	}

	return a.Generate(cfg, ca, "apiserver.crt") //这里的ca会报错，因为类型不符合原先定义的接口，需搭配资产管理修改
}

//用于创建apiserver-kubelet-client.crt，是kube-apiserver 访问 kubelet 所需的客户端证书及私钥。

func GenerateApiServerToKubeletclient() error {
	a := SignedCertKey{}

	ca := &RootCA{} //仍需在资产管理模块完善，未来可以直接调用

	cfg := &CertConfig{
		Subject:      pkix.Name{CommonName: "apiserver-kubelet-client", Organization: []string{"NestOS"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     3650,
	}

	return a.Generate(cfg, ca, "apiserver-kubelet-client.crt")
}

//用于创建apiserver-Etcd-client.crt，是kube-apiserver 访问 Etcd 所需的客户端证书及私钥
func GenerateApiServerToEtcdclient() error {

	a := SignedCertKey{}

	ca := &RootCA{} //仍需在资产管理模块完善，未来可以直接调用

	cfg := &CertConfig{
		Subject:      pkix.Name{CommonName: "apiserver-etcd-client", Organization: []string{"NestOS"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     3650,
	}

	return a.Generate(cfg, ca, "apiserver-etcd-client.crt")
}
