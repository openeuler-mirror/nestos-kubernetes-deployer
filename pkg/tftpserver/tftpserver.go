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

package tftpserver

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pin/tftp"
	"github.com/sirupsen/logrus"
)

type TFTPService struct {
	IP      string
	Port    string
	RootDir string
	server  *tftp.Server
}

func (t *TFTPService) Start() error {
	tftpHandler := TFTPHandler{
		RootDir: t.RootDir,
	}
	tftpServer := tftp.NewServer(tftpHandler.ReadHandler, tftpHandler.WriteHandler)

	tftpServerAddr := t.IP + ":" + t.Port
	logrus.Printf("TFTP server is listening on %s\n", tftpServerAddr)
	err := tftpServer.ListenAndServe(tftpServerAddr)
	if err != nil {
		logrus.Println(err)
		return err
	}

	return nil
}

type TFTPHandler struct {
	RootDir string
}

// ReadHandler handles TFTP read requests
func (h *TFTPHandler) ReadHandler(filename string, rf io.ReaderFrom) error {
	filePath := filepath.Join(h.RootDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = rf.ReadFrom(file)
	if err != nil {
		return err
	}

	return nil
}

// WriteHandler handles TFTP write requests
func (h *TFTPHandler) WriteHandler(filename string, wt io.WriterTo) error {
	filePath := filepath.Join(h.RootDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = wt.WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}

func (t *TFTPService) Stop() error {
	if t.server != nil {
		t.server.Shutdown()
	}
	t.server = nil

	return nil
}
