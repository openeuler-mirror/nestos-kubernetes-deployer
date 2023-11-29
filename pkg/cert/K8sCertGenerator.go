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

func GenerateApiServer(rootcacert, rootcakey []byte) error {
	a := SignedCertKey{}

	cfg := &CertConfig{
		Subject:      pkix.Name{CommonName: "apiserver server", Organization: []string{"NestOS"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     3650,
	}

	return a.Generate(cfg, rootcacert, rootcakey) //
}

//用于创建apiserver-kubelet-client.crt，是kube-apiserver 访问 kubelet 所需的客户端证书及私钥。

func GenerateApiServerToKubeletclient(rootcacert, rootcakey []byte) error {
	a := SignedCertKey{}

	cfg := &CertConfig{
		Subject:      pkix.Name{CommonName: "apiserver-kubelet-client", Organization: []string{"NestOS"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     3650,
	}

	return a.Generate(cfg, rootcacert, rootcakey)
}

//用于创建apiserver-etcd-client.crt，是kube-apiserver 访问 Etcd 所需的客户端证书及私钥
func GenerateApiServerToEtcdclient(rootcacert, rootcakey []byte) error {

	a := SignedCertKey{}

	cfg := &CertConfig{
		Subject:      pkix.Name{CommonName: "apiserver-etcd-client", Organization: []string{"NestOS"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     3650,
	}

	return a.Generate(cfg, rootcacert, rootcakey)
}
