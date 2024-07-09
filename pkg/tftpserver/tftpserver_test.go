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
	"testing"
)

type RFWT struct {
	RootDir string
}

func (rf *RFWT) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}

func (rf *RFWT) WriteTo(w io.Writer) (n int64, err error) {
	return 0, nil
}

func TestTFTPHandler(t *testing.T) {
	tftpHandler := TFTPHandler{
		RootDir: "",
	}
	sf := &RFWT{
		RootDir: "./",
	}

	t.Run("ReadHandler_file_fail", func(t *testing.T) {
		err := tftpHandler.ReadHandler("tftpserver.ssss", sf)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("success")
	})
	t.Run("ReadHandler", func(t *testing.T) {
		err := tftpHandler.ReadHandler("tftpserver.go", sf)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("success")
	})
	p := "tftpservers"
	t.Run("WriteHandler", func(t *testing.T) {
		err := tftpHandler.WriteHandler(p, sf)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("success")
		err = os.RemoveAll(p)
		if err != nil {
			t.Errorf("Error removing directory: %s\n", err)
		} else {
			t.Log("Directory removed successfully")
		}
	})

}

func TestTFTPServer(t *testing.T) {
	service := NewTFTPService("127.0.0.1", "69", "testDir")
	t.Run("start", func(t *testing.T) {
		service.Start()
		t.Log("sssssss")
	})
	t.Run("stop", func(t *testing.T) {
		service.server = nil
		err := service.Stop()
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("stop success")
	})

}
