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

import "nestos-kubernetes-deployer/pkg/utils"

type NodeAsset struct {
	Hostname string
	IP       string
	HardwareInfo
	Certs []utils.StorageContent `json:"-" yaml:"-"` // Certificates content (not printed in JSON and YAML)
}

type HardwareInfo struct {
	CPU  uint
	RAM  uint
	Disk uint
}
