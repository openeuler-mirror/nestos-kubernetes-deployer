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
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// HttpService encapsulates the properties of the HTTP file service
type HttpService struct {
	Port      string
	DirPath   string
	running   bool
	fileCache map[string][]byte
	server    *http.Server
	mutex     sync.RWMutex
}

// NewFileService creates a new instance of file service
func NewFileService(port string) *HttpService {
	return &HttpService{
		Port:      port,
		running:   false,
		fileCache: make(map[string][]byte),
	}
}

// NewDirService creates a new instance of dir service
func NewDirService(port string, dirPath string) *HttpService {
	return &HttpService{
		Port:    port,
		DirPath: dirPath,
		running: false,
	}
}

// AddFileToCache add file content to the file cache
func (hs *HttpService) AddFileToCache(fileName string, content []byte) error {
	if len(content) == 0 {
		return fmt.Errorf("failed to add file '%s' to cache: content is empty", fileName)
	}
	hs.mutex.Lock()
	defer hs.mutex.Unlock()

	hs.fileCache[fileName] = content

	return nil
}

func (hs *HttpService) Start() error {
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

	// 处理文件请求
	smux.HandleFunc("/file/", func(w http.ResponseWriter, r *http.Request) {
		rpath := r.URL.Path[len("/file/"):]

		fileContent, ok := hs.fileCache[rpath]
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

	go func() {
		logrus.Infof("HTTP server is listening on port %s...\n", hs.Port)
		hs.running = true
		if hs.server != nil {
			if err := hs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logrus.Errorf("ListenAndServe(): %v", err)
				hs.running = false
				return
			}
		} else {
			logrus.Error("Server is nil. Cannot start.")
		}
	}()
	return nil
}

func (hs *HttpService) Stop() error {
	if !hs.running {
		logrus.Warn("HTTP server is not running.")
		return nil
	}

	if hs.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := hs.server.Shutdown(ctx); err != nil {
			logrus.Errorf("Shut down the http server: %v", err)
			return err
		}
		hs.server = nil
	}

	hs.running = false
	logrus.Infof("HTTP server stopped.")

	return nil
}
