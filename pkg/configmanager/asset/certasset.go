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

package asset

//接受用户自定义的各类ca证书路径
type CertAsset struct {
	RootCaCertPath       string
	RootCaKeyPath        string
	EtcdCaCertPath       string
	EtcdCaKeyPath        string
	FrontProxyCaCertPath string
	FrontProxyCaKeyPath  string
}
