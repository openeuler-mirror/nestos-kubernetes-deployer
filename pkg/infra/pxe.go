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
	"nestos-kubernetes-deployer/pkg/httpserver"
	"nestos-kubernetes-deployer/pkg/tftpserver"
)

type PXE struct {
	HTTPServerPort string
	HTTPRootDir    string
	TFTPServerIP   string
	TFTPServerPort string
	TFTPRootDir    string
	HTTPService    *httpserver.HTTPService
}

func (p *PXE) deployHTTP(port string, dirPath string) error {
	p.HTTPService.Port = port
	p.HTTPService.DirPath = dirPath

	if err := p.HTTPService.Start(); err != nil {
		return err
	}

	return nil
}

func (p *PXE) deployTFTP(ip string, port string, rootDir string) error {
	tftpService := &tftpserver.TFTPService{
		IP:      ip,
		Port:    port,
		RootDir: rootDir,
	}

	if err := tftpService.Start(); err != nil {
		return err
	}
	defer tftpService.Stop()

	return nil
}

func (p *PXE) Deploy() error {
	if err := p.deployHTTP(p.HTTPServerPort, p.HTTPRootDir); err != nil {
		return err
	}
	if err := p.deployTFTP(p.TFTPServerIP, p.TFTPServerPort, p.TFTPRootDir); err != nil {
		return err
	}

	return nil
}

func (p *PXE) Extend() error {
	if err := p.deployHTTP(p.HTTPServerPort, p.HTTPRootDir); err != nil {
		return err
	}
	if err := p.deployTFTP(p.TFTPServerIP, p.TFTPServerPort, p.TFTPRootDir); err != nil {
		return err
	}

	return nil
}

func (p *PXE) Destroy() error {
	return nil
}
