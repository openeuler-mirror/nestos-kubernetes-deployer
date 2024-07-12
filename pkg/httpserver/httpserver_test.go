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

package httpserver

import (
	"nestos-kubernetes-deployer/pkg/constants"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestHTTPServer(t *testing.T) {
	hs := NewHTTPService("1234")

	t.Run("TestAddFileToCache", func(t *testing.T) {
		var content = []byte("test")
		if err := hs.AddFileToCache("test", content); err != nil {
			t.Log("test fail", err)
			return
		}
		if cachedContent, ok := hs.FileCache["/test"]; !ok || string(cachedContent) != "test" {
			t.Log("test fail: cached content mismatch")
			return
		}
	})

	t.Run("TestStartHTTPService", func(t *testing.T) {
		go func() {
			StartHTTPService(hs)
		}()
		time.Sleep(1 * time.Second)
		if err := hs.Stop(); err != nil {
			t.Log("test fail", err)
			return
		}
	})

	t.Run("TestStartStop", func(t *testing.T) {
		hs.DirPath = "tmp"
		go func() {
			if err := hs.Stop(); err != nil {
				t.Log("test fail", err)
				return
			}
			if err := hs.Start(); err != nil {
				t.Log("test fail", err)
				return
			}
		}()
		time.Sleep(1 * time.Second)
		if err := hs.Stop(); err != nil {
			t.Log("test fail", err)
			return
		}
	})

	t.Run("TestServer", func(t *testing.T) {
		content := []byte("test content")
		hs.AddFileToCache("/testfile", content)

		go func() {
			if err := hs.Stop(); err != nil {
				t.Log("test fail", err)
				return
			}
			if err := hs.Start(); err != nil {
				t.Log("test fail", err)
				return
			}
		}()
		time.Sleep(1 * time.Second)

		_, err := http.Get("http://localhost:1234/testfile")
		if err != nil {
			t.Log("test fail", err)
			return
		}

		_, err = http.Get("http://localhost:1234/dir" + os.TempDir())
		if err != nil {
			t.Log("test fail", err)
			return
		}

		_, err = http.Get("http://localhost:1234" + constants.RpmPackageList)
		if err != nil {
			t.Log("test fail", err)
			return
		}

		if err := hs.Stop(); err != nil {
			t.Log("test fail", err)
			return
		}
	})
}
