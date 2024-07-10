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
	"bytes"
	"fmt"
	"os"
	"runtime"
	"testing"
)

func TestNewVersionCommand(t *testing.T) {
	cmd := NewVersionCommand()

	// 采集命令输出数据
	var output bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	w.Close()
	os.Stdout = stdout
	output.ReadFrom(r)

	// 获取期望数据
	expectedOutput := fmt.Sprintf("Version:    0.3.0\nOS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	// 对比命令输出数据和期望数据
	if output.String() != expectedOutput {
		t.Logf("Expected output %q but got %q", expectedOutput, output.String())
	}
}
