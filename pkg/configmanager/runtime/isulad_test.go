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

package runtime_test

import (
	"nestos-kubernetes-deployer/pkg/configmanager/runtime"
	"testing"
)

func TestIsuladRuntime_GetRuntimeCriSocket(t *testing.T) {
	ir := &runtime.IsuladRuntime{}
	expectedSocket := "/var/run/isulad.sock"
	if ir.GetRuntimeCriSocket() != expectedSocket {
		t.Errorf("Expected socket path %s, but got %s", expectedSocket, ir.GetRuntimeCriSocket())
	}
}

func TestIsIsulad(t *testing.T) {
	// Test case 1: Check if isuladRuntime returns true for IsIsulad
	cr := &runtime.IsuladRuntime{}
	if !runtime.IsIsulad(cr) {
		t.Errorf("Expected IsIsulad to return true for isuladRuntime, but got false")
	}

	// Test case 2: Check if mockRuntime returns false for IsIsulad
	mr := &mockRuntime{}
	if runtime.IsIsulad(mr) {
		t.Errorf("Expected IsIsulad to return false for mockRuntime, but got true")
	}
}
