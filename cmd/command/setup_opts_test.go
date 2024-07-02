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
package command

import (
	"testing"

	"github.com/spf13/cobra"
)

// Create a new test command
func setupCmd(cmdSetupFunc func(*cobra.Command)) *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmdSetupFunc(cmd)
	return cmd
}

func TestCommandSetups(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*cobra.Command)
	}{
		{name: "DeployCmdOpts", setupFunc: SetupDeployCmdOpts},
		{name: "DestroyCmdOpts", setupFunc: SetupDestroyCmdOpts},
		{name: "UpgradeCmdOpts", setupFunc: SetupUpgradeCmdOpts},
		{name: "ExtendCmdOpts", setupFunc: SetupExtendCmdOpts},
		{name: "TemplateCmdOpts", setupFunc: SetupTemplateCmdOpts},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := setupCmd(tt.setupFunc)
			if cmd.Use != "test" {
				t.Errorf("expected command use 'test', got '%s'", cmd.Use)
			}
		})
	}
}
