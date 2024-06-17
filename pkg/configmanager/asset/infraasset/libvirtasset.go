/*
Copyright 2024 KylinSoft  Co., Ltd.

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

package infraasset

import (
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
)

type LibvirtAsset struct {
	URI     string
	OSPath  string `yaml:"osPath"`
	CIDR    string
	Gateway string
}

func (la *LibvirtAsset) InitAsset(libvirtMap map[string]interface{}, opts *opts.OptionsList, args ...interface{}) (InfraAsset, error) {
	updateFieldFromMap("uri", &la.URI, libvirtMap)
	asset.SetStringValue(&la.URI, opts.InfraPlatform.Libvirt.URI, "qemu:///system")

	updateFieldFromMap("osPath", &la.OSPath, libvirtMap)
	if err := asset.CheckStringValue(&la.OSPath, opts.InfraPlatform.Libvirt.OSPath, "libvirt-osPath"); err != nil {
		return nil, err
	}

	updateFieldFromMap("cidr", &la.CIDR, libvirtMap)
	asset.SetStringValue(&la.CIDR, opts.InfraPlatform.Libvirt.CIDR, "192.168.132.0/24")

	updateFieldFromMap("gateway", &la.Gateway, libvirtMap)
	asset.SetStringValue(&la.Gateway, opts.InfraPlatform.Libvirt.Gateway, "192.168.132.1")

	return la, nil
}
