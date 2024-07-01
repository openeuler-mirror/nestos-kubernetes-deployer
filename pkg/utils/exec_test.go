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
package utils_test

import (
	"nestos-kubernetes-deployer/pkg/utils"
	"testing"
)

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name       string
		command    string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "Valid command",
			command:    "echo 'Hello, World!'",
			wantOutput: "Hello, World!\n",
			wantErr:    false,
		},
		{
			name:    "Invalid command",
			command: "invalidcommand",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, err := utils.RunCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutput != tt.wantOutput {
				t.Errorf("RunCommand() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
