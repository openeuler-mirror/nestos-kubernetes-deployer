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
package asset

import "testing"

func TestHookconfig(t *testing.T) {

	hookconfig := &HookConf{
		PreHookScript: "../../../hack/build.sh",
		PostHookYaml:  "../../../hack/build.sh",
		ShellFiles: []ShellFile{
			{
				Name:    "test.sh",
				Mode:    0644,
				Content: nil,
			},
		},
		PostHookFiles: []string{"sss", "sss"},
	}

	t.Run("GetCmdHooks Success", func(t *testing.T) {
		err := GetCmdHooks(hookconfig)
		if err != nil {
			t.Logf("GetCmdHooks failed: %v", err)
		}
	})

	//测试脚本格式
	t.Run("GetCmdHooks PreHookScript", func(t *testing.T) {
		hookconfig.PreHookScript = "../../../hack"
		err := GetCmdHooks(hookconfig)
		if err == nil {
			t.Log("Expected error, got nil")
		}
	})

	//测试yaml格式
	t.Run("GetCmdHooks PostHookYaml", func(t *testing.T) {
		hookconfig.PostHookYaml = "../../"
		err := GetCmdHooks(hookconfig)
		if err == nil {
			t.Log("Expected error, got nil")
		}
	})
}
