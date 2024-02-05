/*
Copyright 2023 KylinSoft  Co., Ltd.

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

// service.go
package httpserver

import (
	"errors"
	"net/http"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

// HttpFileService encapsulates the properties of the HTTP file service
type HttpFileService struct {
	Port      string
	server    *http.Server
	running   bool
	fileCache map[string][]byte
	mutex     sync.RWMutex
}

// NewFileService creates a new instance of file service
func NewFileService(port string) *HttpFileService {
	return &HttpFileService{
		Port:      port,
		running:   false,
		fileCache: make(map[string][]byte),
	}
}

// AddFileToCache add file content to the file cache
func (fs *HttpFileService) AddFileToCache(fileName string, content []byte) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	fileName = "/" + fileName
	fs.fileCache[fileName] = content
}

// RemoveFileFromCache removes file content from the file cache
func (fs *HttpFileService) RemoveFileFromCache(fileName string) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	delete(fs.fileCache, fileName)
}

func (fs *HttpFileService) Start() error {
	// Set up HTTP route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := fs.handleFileRequest(w, r)
		if err != nil {
			logrus.Errorf("Error handling file request: %v", err)

			var statusCode int
			var errorMessage string

			if os.IsNotExist(err) {
				statusCode = http.StatusNotFound
				errorMessage = "File Not Found"
			} else {
				statusCode = http.StatusInternalServerError
				errorMessage = "Internal Server Error"
			}

			http.Error(w, errorMessage, statusCode)
			return
		}
	})

	fs.server = &http.Server{
		Addr: ":" + fs.Port,
	}

	go func() {
		logrus.Infof("HTTP server listening on port %s...\n", fs.Port)
		fs.running = true
		if fs.server != nil {
			if err := fs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logrus.Errorf("ListenAndServe(): %v", err)
				fs.running = false
				return
			}
		} else {
			logrus.Error("Server is nil. Cannot start.")
		}
	}()

	return nil
}

// handleFileRequest handles file requests
func (fs *HttpFileService) handleFileRequest(w http.ResponseWriter, r *http.Request) error {
	// Get the requested file path
	filePath := r.URL.Path
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	// Check if the file exists in the cache
	fileContent, ok := fs.fileCache[filePath]
	if !ok || len(fileContent) == 0 {
		return os.ErrNotExist
	}

	// Set the content type of the file
	contentType := http.DetectContentType(fileContent)
	w.Header().Set("Content-Type", contentType)

	// Write file content directly into the response
	_, err := w.Write(fileContent)
	if err != nil {
		errMsg := "unable to write file to response: " + err.Error()
		return errors.New(errMsg)
	}

	return nil
}

// Stop method stops the file service
func (fs *HttpFileService) Stop() error {
	if !fs.running || fs.server == nil {
		logrus.Warn("Server is not running.")
		return nil
	}

	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	logrus.Info("Stopping http server...")
	if err := fs.server.Close(); err != nil {
		logrus.Errorf("Error closing server: %v", err)
		return errors.New("error closing server: " + err.Error())
	}

	// Clear the file cache
	for fileName := range fs.fileCache {
		delete(fs.fileCache, fileName)
	}

	fs.running = false
	return nil
}
