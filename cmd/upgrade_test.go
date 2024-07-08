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
	"testing"
)

func TestUpgrade(t *testing.T) {
	cmd := NewUpgradeCommand()
	args := []string{"--cluster-id", "cluster", "--imageurl", "http://example.com/image"}
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("Failed to execute command: %v", err)
	}

	t.Run("UpgradeCmd Fail", func(t *testing.T) {
		if err := runUpgradeCmd(cmd, args); err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
