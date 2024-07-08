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
	"testing"
)

func TestExtend(t *testing.T) {
	opts.Opts.RootOptDir = "../data"
	cmd := NewExtendCommand()
	args := []string{"--num", "0"}
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("Failed to execute command: %v", err)
	}

	t.Run("ExtendCmd Fail", func(t *testing.T) {
		if err := runExtendCmd(cmd, args); err == nil {
			t.Error("Expected error, got nil")
		}
	})
}