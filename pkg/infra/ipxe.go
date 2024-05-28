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

package infra

import (
	"fmt"
	"nestos-kubernetes-deployer/pkg/constants"
	"nestos-kubernetes-deployer/pkg/httpserver"
	"os"
	"time"
)

type IPXE struct {
	Port              string
	FilePath          string
	OSInstallTreePath string
	HTTPService       *httpserver.HTTPService
}

func (i *IPXE) deployHTTP(port string, dirPath string, filePath string) error {
	i.HTTPService.Port = port
	i.HTTPService.DirPath = dirPath
	i.HTTPService.HttpLastRequestTime = time.Now().Unix()

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := i.HTTPService.AddFileToCache(constants.IPXECfg, fileContent); err != nil {
		return err
	}

	if err := i.HTTPService.Start(); err != nil {
		return fmt.Errorf("error starting file service: %v", err)
	}

	return nil
}

func (i *IPXE) Deploy() error {
	return i.deployHTTP(i.Port, i.OSInstallTreePath, i.FilePath)
}

func (i *IPXE) Extend() error {
	return i.Deploy()
}

func (i *IPXE) Destroy() error {
	return nil
}
