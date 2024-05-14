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

type IPXEAsset struct {
	IPXEPort              string `yaml:"ipxe_port"`
	IPXEFilePath          string `yaml:"ipxe_file_path"`
	IPXEOSInstallTreePath string `yaml:"ipxe_os_install_tree_path"`
}

func (ia *IPXEAsset) InitAsset(ipxeMap map[string]interface{}, opts *opts.OptionsList, args ...interface{}) (InfraAsset, error) {
	updateFieldFromMap("ipxe_port", &ia.IPXEPort, ipxeMap)
	asset.SetStringValue(&ia.IPXEPort, opts.InfraPlatform.IPXE.IPXEPort, "9080")

	updateFieldFromMap("ipxe_file_path", &ia.IPXEFilePath, ipxeMap)
	if err := asset.CheckStringValue(&ia.IPXEFilePath, opts.InfraPlatform.IPXE.IPXEFilePath, "ipxe-file-path"); err != nil {
		return nil, err
	}

	updateFieldFromMap("ipxe_os_install_tree_path", &ia.IPXEOSInstallTreePath, ipxeMap)
	asset.SetStringValue(&ia.IPXEOSInstallTreePath, opts.InfraPlatform.IPXE.IPXEOSInstallTreePath, "/var/www/html/")

	return ia, nil
}
