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
	"context"
	"errors"
	"fmt"
	"nestos-kubernetes-deployer/pkg/constants"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const TimeOut = 10 * 60

// HTTPService encapsulates the properties of the HTTP file service
type HTTPService struct {
	Port                string
	DirPath             string
	PackageDir          string
	FileCache           map[string][]byte
	running             bool
	server              *http.Server
	mutex               sync.RWMutex
	HttpLastRequestTime int64 `json:"http_last_request_time"`
	Ch                  chan struct{}
}

func NewHTTPService(port string) *HTTPService {
	return &HTTPService{
		Port:                port,
		FileCache:           make(map[string][]byte),
		HttpLastRequestTime: time.Now().Unix(),
		Ch:                  make(chan struct{}, 1),
	}
}

// AddFileToCache add file content to the file cache
func (hs *HTTPService) AddFileToCache(fileName string, content []byte) error {
	if len(content) == 0 {
		return fmt.Errorf("failed to add file '%s' to cache: content is empty", fileName)
	}
	hs.mutex.Lock()
	defer hs.mutex.Unlock()

	if !strings.HasPrefix(fileName, "/") {
		fileName = "/" + fileName
	}
	hs.FileCache[fileName] = content

	return nil
}

func (hs *HTTPService) Start() error {
	// Check if the http server is already running
	if hs.running {
		return errors.New("HTTP server is already running")
	}

	var dirPath string
	if hs.DirPath != "" {
		var err error
		dirPath, err = filepath.Abs(hs.DirPath)
		if err != nil {
			return err
		}
	}

	hs.mutex.Lock()
	defer hs.mutex.Unlock()

	smux := http.NewServeMux()
	// 处理目录请求
	smux.HandleFunc("/dir/", func(w http.ResponseWriter, r *http.Request) {
		hs.HttpLastRequestTime = time.Now().Unix()
		rpath := filepath.Join(dirPath, r.URL.Path[len("/dir/"):])
		_, err := os.Stat(rpath)
		if err != nil {
			// 如果请求对应目录，返回目录下的文件列表
			http.FileServer(http.Dir(rpath)).ServeHTTP(w, r)
			return
		}

		// 如果请求是文件，返回文件内容
		http.ServeFile(w, r, rpath)
	})

	// 处理rpm软件包请求
	smux.HandleFunc(constants.RpmPackageList, func(w http.ResponseWriter, r *http.Request) {
		hs.HttpLastRequestTime = time.Now().Unix()
		rpath := filepath.Join(hs.PackageDir, r.URL.Path[len(constants.RpmPackageList):])
		_, err := os.Stat(rpath)
		if err != nil {
			http.FileServer(http.Dir(rpath)).ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, rpath)
	})

	// 处理文件请求
	smux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hs.HttpLastRequestTime = time.Now().Unix()
		rpath := r.URL.Path
		fileContent, ok := hs.FileCache[rpath]
		if !ok {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "%s", fileContent)
	})

	hs.server = &http.Server{
		Addr:    ":" + hs.Port,
		Handler: smux,
	}

	logrus.Infof("HTTP server is listening on port %s...\n", hs.Port)
	go func() {
		for {
			val := time.Now().Unix() - hs.HttpLastRequestTime
			if val < TimeOut {
				time.Sleep(30 * time.Second)
				continue
			}
			hs.server.Close()
			hs.running = false
			hs.server = nil
			hs.Ch <- struct{}{}
			return
		}
	}()
	hs.running = true
	if err := hs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Errorf("ListenAndServe(): %v", err)
		hs.running = false
		return err
	}
	return nil
}

func (hs *HTTPService) Stop() error {
	if !hs.running {
		logrus.Warn("HTTP server is not running.")
		return nil
	}

	if hs.server == nil {
		hs.running = false
		logrus.Infof("HTTP server stopped.")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := hs.server.Shutdown(ctx); err != nil {
		logrus.Errorf("Shut down the http server: %v", err)
		return err
	}
	hs.server = nil
	return nil
}
