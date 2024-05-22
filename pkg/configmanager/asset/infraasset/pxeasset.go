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

type PXEAsset struct {
	HTTPServerPort string `yaml:"httpServerPort"`
	HTTPRootDir    string `yaml:"httpRootDir"`
	TFTPServerIP   string `yaml:"tftpServerIP"`
	TFTPServerPort string `yaml:"tftpServerPort"`
	TFTPRootDir    string `yaml:"tftpRootDir"`
}

func (pa *PXEAsset) InitAsset(pxeMap map[string]interface{}, opts *opts.OptionsList, args ...interface{}) (InfraAsset, error) {
	updateFieldFromMap("httpServerPort", &pa.HTTPServerPort, pxeMap)
	asset.SetStringValue(&pa.HTTPServerPort, opts.InfraPlatform.PXE.HTTPServerPort, "9080")

	updateFieldFromMap("httpRootDir", &pa.HTTPRootDir, pxeMap)
	asset.SetStringValue(&pa.HTTPRootDir, opts.InfraPlatform.PXE.HTTPRootDir, "/var/www/html/")

	updateFieldFromMap("tftpServerIP", &pa.TFTPServerIP, pxeMap)
	if err := asset.CheckStringValue(&pa.TFTPServerIP, opts.InfraPlatform.PXE.TFTPServerIP, "tftp-server-ip"); err != nil {
		return nil, err
	}

	updateFieldFromMap("tftpServerPort", &pa.TFTPServerPort, pxeMap)
	asset.SetStringValue(&pa.TFTPServerPort, opts.InfraPlatform.PXE.TFTPServerPort, "69")

	updateFieldFromMap("tftpRootDir", &pa.TFTPRootDir, pxeMap)
	asset.SetStringValue(&pa.TFTPRootDir, opts.InfraPlatform.PXE.TFTPRootDir, "/var/lib/tftpboot/")

	return pa, nil
}
