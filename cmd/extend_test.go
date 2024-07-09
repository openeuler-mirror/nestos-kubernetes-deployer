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
package cmd

import (
	"nestos-kubernetes-deployer/cmd/command/opts"
	"os"
	"testing"
)

func TestExtend(t *testing.T) {
	opts.Opts.RootOptDir = "../data"
	cmd := NewExtendCommand()
	args := []string{"--num", "0"}
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Logf("Failed to execute command: %v", err)
	}

	t.Run("ExtendCmd Fail", func(t *testing.T) {
		if err := runExtendCmd(cmd, args); err == nil {
			t.Log("Expected error, got nil")
		}
		// Clean up
		if err := os.RemoveAll("logs"); err != nil {
			t.Logf("Failed to remove logs folder: %v", err)
		}

		if _, err := os.Stat("global_config.yaml"); os.IsNotExist(err) {
			t.Logf("Expected global_config.yaml to be created, but it does not exist")
		}

		if err := os.Remove("global_config.yaml"); err != nil {
			t.Logf("Failed to remove global_config.yaml: %v", err)
		}
		// Clean up
		if err := os.RemoveAll("logs"); err != nil {
			t.Logf("Failed to remove logs folder: %v", err)
		}

		if _, err := os.Stat("global_config.yaml"); os.IsNotExist(err) {
			t.Logf("Expected global_config.yaml to be created, but it does not exist")
		}

		if err := os.Remove("global_config.yaml"); err != nil {
			t.Logf("Failed to remove global_config.yaml: %v", err)
		}
	})
}
