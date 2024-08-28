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
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestHTTPServer(t *testing.T) {
	hs := NewHTTPService("9876")

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

	t.Run("StartHTTPService", func(t *testing.T) {
		hs.Port = "8520"
		StartHTTPService(hs)
	})

	t.Run("Start", func(t *testing.T) {
		hs.Port = "3698"
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/dir" {
				// 根据您的预期响应设置响应头和响应体
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"message": "API endpoint success"}`))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer ts.Close()
		//t.Log(ts.URL + "/dir")
		//// 发送 HTTP 请求
		//client := &http.Client{}
		//req, err := http.NewRequest("GET", ts.URL+"/dir", nil)
		//if err != nil {
		//	t.Log(err)
		//}
		//
		//resp, err := client.Do(req)
		//if err != nil {
		//	t.Log(err)
		//}
		//defer resp.Body.Close()
		//return
		hs := &HTTPService{
			server: &http.Server{
				Addr:    ts.URL[7:],
				Handler: ts.Config.Handler,
			},
		}
		hs.HttpLastRequestTime = time.Now().Unix() - TimeOut + 10
		err := hs.Start()
		if err != nil {
			t.Log(err)
			return
		}
		_, err = http.Get("http://localhost:3698/testfile")
		if err != nil {
			t.Log("test fail", err)
			return
		}

		_, err = http.Get("http://localhost:3698/dir" + os.TempDir())
		if err != nil {
			t.Log("test fail", err)
			return
		}

		_, err = http.Get("http://localhost:3698" + constants.RpmPackageList)
		if err != nil {
			t.Log("test fail", err)
			return
		}

		t.Log("start success")
	})

	t.Run("stop", func(t *testing.T) {
		err := hs.Stop()
		if err != nil {
			t.Log(err)
			return
		}

		t.Log("stop1 success")
	})

	t.Run("stop_runing", func(t *testing.T) {
		hs.running = false
		err := hs.Stop()
		if err != nil {
			t.Log(err)
			return
		}
		t.Log("stop2 success")
	})

	t.Run("stop_server_empty", func(t *testing.T) {
		hs.running = true
		hs.server = nil
		err := hs.Stop()
		if err != nil {
			t.Log(err)
			return
		}
		t.Log("stop2 success")
	})

}
