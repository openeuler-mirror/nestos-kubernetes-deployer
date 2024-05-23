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
	"nestos-kubernetes-deployer/pkg/httpserver"
	"os"
)

type IPXE struct {
	IPXEPort              string
	IPXEFilePath          string
	IPXEOSInstallTreePath string
	HTTPService           *httpserver.HTTPService
}

func (i *IPXE) deployHTTP(port string, dirPath string, filePath string) error {
	i.HTTPService.Port = port
	i.HTTPService.DirPath = dirPath

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := i.HTTPService.AddFileToCache(filePath, fileContent); err != nil {
		return err
	}

	if err := i.HTTPService.Start(); err != nil {
		return fmt.Errorf("error starting file service: %v", err)
	}

	return nil
}

func (i *IPXE) Deploy() error {
	if err := i.deployHTTP(i.IPXEPort, i.IPXEOSInstallTreePath, i.IPXEFilePath); err != nil {
		return err
	}

	return nil
}

func (i *IPXE) Extend() error {
	if err := i.deployHTTP(i.IPXEPort, i.IPXEOSInstallTreePath, i.IPXEFilePath); err != nil {
		return err
	}

	return nil
}

func (i *IPXE) Destroy() error {
	return nil
}
