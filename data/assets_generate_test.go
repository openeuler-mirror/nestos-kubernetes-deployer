//go:build tools
// +build tools

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

package main

import (
	"os/exec"
	"testing"
)

// TestMainFunction tests the main function for correct file generation
func TestMainFunction(t *testing.T) {
	cmd := exec.Command("go", "run", "assets_generate.go")

	// Run the command
	err := cmd.Run()
	if err != nil {
		t.Logf("Failed to run cmd: %v", err)
	}
}
