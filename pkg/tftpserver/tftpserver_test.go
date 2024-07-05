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
	"testing"
)

type RFWT struct {
	RootDir string
}

func (rf *RFWT) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}

func (wt *RFWT) WriteTo(w io.Writer) (n int64, err error) {
	return 0, nil
}

func TestTFTPServer(t *testing.T) {
	service := NewTFTPService("127.0.0.1", "69", "testDir")

	t.Run("TestStartStop", func(t *testing.T) {
		go func() {
			if err := service.Start(); err != nil {
				t.Error("test fail", err)
				return
			}
		}()
		if err := service.Stop(); err != nil {
			t.Error("test fail", err)
			return
		}
	})

	th := &TFTPHandler{
		RootDir: "testDir",
	}
	t.Run("TestReadHandler", func(t *testing.T) {
		rf := &RFWT{}

		if err := th.ReadHandler("test", rf); err != nil {
			t.Error("test fail", err)
			return
		}
	})
	t.Run("TestWriteHandler", func(t *testing.T) {
		wt := &RFWT{}

		if err := th.WriteHandler("test", wt); err != nil {
			t.Error("test fail", err)
			return
		}
	})
}
