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
	OSImage string
	CIDR    string
	Gateway string
}

func (la *LibvirtAsset) InitAsset(libvirtMap map[string]interface{}, opts *opts.OptionsList, args ...interface{}) (InfraAsset, error) {
	updateFieldFromMap("uri", &la.URI, libvirtMap)
	asset.SetStringValue(&la.URI, opts.InfraPlatform.Libvirt.URI, "qemu:///system")

	updateFieldFromMap("osimage", &la.OSImage, libvirtMap)
	osImage := "https://nestos.org.cn/nestos20230928/nestos-for-container/x86_64/NestOS-For-Container-22.03-LTS-SP2.20230928.0-qemu.x86_64.qcow2"
	if len(args) > 0 {
		if arch, ok := args[0].(string); ok {
			if arch == "arm64" || arch == "aarch64" {
				osImage = "https://nestos.org.cn/nestos20230928/nestos-for-container/aarch64/NestOS-For-Container-22.03-LTS-SP2.20230928.0-qemu.aarch64.qcow2"
			}
		}
	}
	asset.SetStringValue(&la.OSImage, opts.InfraPlatform.Libvirt.OSImage, osImage)

	updateFieldFromMap("cidr", &la.CIDR, libvirtMap)
	asset.SetStringValue(&la.CIDR, opts.InfraPlatform.Libvirt.CIDR, "192.168.132.0/24")

	updateFieldFromMap("gateway", &la.Gateway, libvirtMap)
	asset.SetStringValue(&la.Gateway, opts.InfraPlatform.Libvirt.Gateway, "192.168.132.1")

	return la, nil
}
