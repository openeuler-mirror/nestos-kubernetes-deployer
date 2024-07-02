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

func TestDockerRuntime_GetRuntimeCriSocket(t *testing.T) {
	dr := &dockerRuntime{}
	expectedSocket := "/var/run/dockershim.sock"
	if dr.GetRuntimeCriSocket() != expectedSocket {
		t.Errorf("Expected socket path %s, but got %s", expectedSocket, dr.GetRuntimeCriSocket())
	}
}

func TestIsDocker(t *testing.T) {
	// Test case 1: Check if dockerRuntime returns true for IsDocker
	cr := &dockerRuntime{}
	if !IsDocker(cr) {
		t.Errorf("Expected Isdocker to return true for dockerRuntime, but got false")
	}

	// Test case 2: Check if mockRuntime returns false for IsDocker
	mr := &mockRuntime{}
	if IsDocker(mr) {
		t.Errorf("Expected IsDocker to return false for mockRuntime, but got true")
	}
}
