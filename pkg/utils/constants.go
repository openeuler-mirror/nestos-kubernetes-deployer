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
package utils

import "os"

const (
	// certificates relate constants
	EtcdCaKey                 = "/etc/kubernetes/pki/etcd/ca.key"
	EtcdCaCrt                 = "/etc/kubernetes/pki/etcd/ca.crt"
	ApiserverEtcdClientKey    = "/etc/kubernetes/pki/apiserver-etcd-client.key"
	ApiserverEtcdClientCrt    = "/etc/kubernetes/pki/apiserver-etcd-client.crt"
	CaKey                     = "/etc/kubernetes/pki/ca.key"
	CaCrt                     = "/etc/kubernetes/pki/ca.crt"
	ApiserverKey              = "/etc/kubernetes/pki/apiserver.key"
	ApiserverCrt              = "/etc/kubernetes/pki/apiserver.crt"
	ApiserverKubeletClientKey = "/etc/kubernetes/pki/apiserver-kubelet-client.key"
	ApiserverKubeletClientCrt = "/etc/kubernetes/pki/apiserver-kubelet-client.crt"
	FrontProxyCaKey           = "/etc/kubernetes/pki/front-proxy-ca.key"
	FrontProxyCaCrt           = "/etc/kubernetes/pki/front-proxy-ca.crt"
	FrontProxyClientKey       = "/etc/kubernetes/pki/front-proxy-client.key"
	FrontProxyClientCrt       = "/etc/kubernetes/pki/front-proxy-client.crt"
	ServerKey                 = "/etc/kubernetes/pki/etcd/server.key"
	ServerCrt                 = "/etc/kubernetes/pki/etcd/server.crt"
	PeerKey                   = "/etc/kubernetes/pki/etcd/peer.key"
	PeerCrt                   = "/etc/kubernetes/pki/etcd/peer.crt"
	HealthcheckClientKey      = "/etc/kubernetes/pki/etcd/healthcheck-client.key"
	HealthcheckClientCrt      = "/etc/kubernetes/pki/etcd/healthcheck-client.crt"
	SaKey                     = "/etc/kubernetes/pki/sa.key"
	SaPub                     = "/etc/kubernetes/pki/sa.pub"

	AdminConfig       = "/etc/kubernetes/admin.conf"
	KubeletConfig     = "/etc/kubernetes/kubelet.conf"
	ControllerManager = "/etc/kubernetes/controller-manager.conf"
	schedulerConf     = "/etc/kubernetes/scheduler.conf"

	CertFileMode os.FileMode = 0644
)
