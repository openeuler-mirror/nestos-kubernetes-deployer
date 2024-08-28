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

package runtime

import (
	"testing"
)

// Define mockRuntime for testing purposes
type mockRuntime struct{}

func (mr *mockRuntime) GetRuntimeCriSocket() string {
	return "/var/run/mock.sock"
}

func TestContainerdRuntimeGetRuntimeCriSocket(t *testing.T) {
	cr := &containerdRuntime{}
	expectedSocket := "unix:///var/run/containerd/containerd.sock"
	if cr.GetRuntimeCriSocket() != expectedSocket {
		t.Errorf("Expected socket path %s, but got %s", expectedSocket, cr.GetRuntimeCriSocket())
	}
}

func TestIsContainerd(t *testing.T) {
	// Test case 1: Check if containerdRuntime returns true for IsContainerd
	cr := &containerdRuntime{}
	if !IsContainerd(cr) {
		t.Errorf("Expected IsContainerd to return true for containerdRuntime, but got false")
	}

	// Test case 2: Check if mockRuntime returns false for IsContainerd
	mr := &mockRuntime{}
	if IsContainerd(mr) {
		t.Errorf("Expected IsContainerd to return false for mockRuntime, but got true")
	}
}
