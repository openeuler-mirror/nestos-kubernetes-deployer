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
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	const platform = "--platform"
	platformTests := []struct {
		platform string
		args     []string
	}{
		{"openstack", []string{platform, "openstack"}},
		{"pxe", []string{platform, "pxe"}},
		{"ipxe", []string{platform, "ipxe"}},
		{"ipxesss", []string{platform, "ipxesss"}},
		{"libvirt", []string{platform, "libvirt"}},
	}

	cmd := NewTemplateCommand()

	t.Run("Template Fail", func(t *testing.T) {
		err := createTemplate(cmd, nil)
		if err == nil {
			t.Log("Expected error, got nil")
		}
	})

	for _, pt := range platformTests {
		t.Run("Template "+pt.platform+" Success", func(t *testing.T) {
			cmd.SetArgs(pt.args)
			if err := cmd.Execute(); err != nil {
				t.Logf("Failed to execute command: %v", err)
			}

			if err := createTemplate(cmd, pt.args); err != nil {
				t.Logf("createTemplate failed: %v", err)
			}

			// Clean up
			if _, err := os.Stat("template.yaml"); os.IsNotExist(err) {
				t.Logf("Expected template.yaml to be created, but it does not exist")
			}

			if err := os.Remove("template.yaml"); err != nil {
				t.Logf("Failed to remove template.yaml: %v", err)
			}
		})
	}
}
