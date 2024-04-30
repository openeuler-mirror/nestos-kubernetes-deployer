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
	"nestos-kubernetes-deployer/pkg/constants"
	"testing"
)

func TestGetRuntime(t *testing.T) {
	for _, runtimeName := range []string{constants.Isulad, constants.Docker, constants.Crio} {
		rt, err := runtime.GetRuntime(runtimeName)
		if err != nil {
			t.Errorf("Unexpected error for supported runtime %s: %v", runtimeName, err)
		}
		if rt == nil {
			t.Errorf("Unexpected nil runtime for supported runtime %s", runtimeName)
		}
	}

	unknownRuntime := "unknown"
	rt, err := runtime.GetRuntime(unknownRuntime)
	if err == nil {
		t.Errorf("Expected error for unsupported runtime %s", unknownRuntime)
	}
	if rt != nil {
		t.Errorf("Expected nil runtime for unsupported runtime %s", unknownRuntime)
	}
	expectedErrorMessage := "unsupported runtime"
	if err != nil && err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s' for unsupported runtime, got '%s'", expectedErrorMessage, err.Error())
	}
}
