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
	IP                string
	Port              string
	FilePath          string `yaml:"filePath"`
	OSInstallTreePath string `yaml:"osInstallTreePath"`
}

func (ia *IPXEAsset) InitAsset(ipxeMap map[string]interface{}, opts *opts.OptionsList, args ...interface{}) (InfraAsset, error) {
	updateFieldFromMap("port", &ia.Port, ipxeMap)
	asset.SetStringValue(&ia.Port, opts.InfraPlatform.IPXE.Port, "9080")

	updateFieldFromMap("filePath", &ia.FilePath, ipxeMap)
	if err := asset.CheckStringValue(&ia.FilePath, opts.InfraPlatform.IPXE.FilePath, "ipxe-osInstallTreePath"); err != nil {
		return nil, err
	}

	updateFieldFromMap("osInstallTreePath", &ia.OSInstallTreePath, ipxeMap)
	asset.SetStringValue(&ia.OSInstallTreePath, opts.InfraPlatform.IPXE.OSInstallTreePath, "/var/www/html/")

	return ia, nil
}
