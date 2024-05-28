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

	"github.com/sirupsen/logrus"
)

type PXE struct {
	IP             string
	HTTPServerPort string
	HTTPRootDir    string
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
	go func() {
		select {
		case <-p.HTTPService.Ch:
			logrus.Info("tftp server stop")
			tftpService.Stop()
			return
		}
	}()

	if err := tftpService.Start(); err != nil {
		return err
	}
	defer tftpService.Stop()

	return nil
}

func (p *PXE) Deploy() error {
	go func() {
		err := p.deployHTTP(p.HTTPServerPort, p.HTTPRootDir)
		if err != nil {
			logrus.Errorf("PXE deploy http server err: %v", err)
			return
		}
	}()

	if err := p.deployTFTP(p.IP, p.TFTPServerPort, p.TFTPRootDir); err != nil {
		return err
	}

	return nil
}

func (p *PXE) Extend() error {
	return p.Deploy()
}

func (p *PXE) Destroy() error {
	return nil
}
