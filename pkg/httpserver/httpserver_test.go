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

package httpserver_test

import (
	"io"
	"nestos-kubernetes-deployer/pkg/httpserver"
	"net/http"
	"testing"
)

func TestHttpFileService(t *testing.T) {
	// Create a new file service instance
	fileService := httpserver.NewFileService("9080")

	// Start the file service
	if err := fileService.Start(); err != nil {
		t.Fatalf("Error starting file service: %v", err)
	}
	defer fileService.Stop()

	// Add test file to the file service
	testContent := []byte("Hello, world!")
	fileService.AddFileToCache("test.txt", testContent)

	// Make an HTTP request to retrieve the test file content
	resp, err := http.Get("http://localhost:9080/file/test.txt")
	if err != nil {
		t.Fatalf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body content
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	// Check if the response body content matches the expected content
	if string(respBody) != string(testContent) {
		t.Errorf("Expected response body %s, got %s", string(testContent), string(respBody))
	}
}
