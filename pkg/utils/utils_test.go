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
package utils

import (
	"fmt"
	"net"
	"path/filepath"
	"testing"
)

func TestGetKubernetesApiVersion(t *testing.T) {
	tests := []struct {
		versionNumber uint
		expected      string
		expectError   bool
	}{
		{1, "v1beta1", false},
		{2, "v1beta2", false},
		{3, "v1beta3", false},
		{0, "", false},
		{4, "", true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("versionNumber=%d", tt.versionNumber), func(t *testing.T) {
			version, err := GetKubernetesApiVersion(tt.versionNumber)
			if tt.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if version != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, version)
			}
		})
	}
}

func TestGetDefaultPubKeyPath(t *testing.T) {
	expected := filepath.Join(getSysHome(), ".ssh", "id_rsa.pub")
	if result := GetDefaultPubKeyPath(); result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestGetApiServerEndpoint(t *testing.T) {
	ip := "127.0.0.1"
	expected := fmt.Sprintf("%s:%s", ip, "6443")
	if result := GetApiServerEndpoint(ip); result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestGetLocalIP(t *testing.T) {
	ip, err := GetLocalIP()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if net.ParseIP(ip) == nil {
		t.Errorf("expected a valid IP address, got %s", ip)
	}
}

func TestIsPortOpen(t *testing.T) {
	port := "8080"
	if !IsPortOpen(port) {
		t.Errorf("expected port %s to be open", port)
	}
}

func TestConstructURL(t *testing.T) {
	host := "localhost"
	role := "worker"
	expected := "http://localhost/worker"
	if result := ConstructURL(host, role); result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestGenerateWWN(t *testing.T) {
	wwn, err := GenerateWWN()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(wwn) != 16 {
		t.Errorf("expected WWN length 16, got %d", len(wwn))
	}
}
