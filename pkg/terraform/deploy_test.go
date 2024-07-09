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
package terraform

import (
	"os"
	"testing"
)

func TestExecuteApplyTerraform(t *testing.T) {
	err := os.Chdir("../../data")
	if err != nil {
		t.Fatal(err)
	}
	tfFileDir := ""
	persistDir := ""
	t.Run("ExecuteApplyTerraform_empty_dir", func(t *testing.T) {
		terraform, err := ExecuteApplyTerraform(tfFileDir, persistDir)
		if err != nil {
			t.Error(err)
		}
		t.Log(string(terraform))
	})
	///media/weihuanhuan/本地磁盘/weihuanhuan/goWorkspace/src/nestos-kubernetes-deployer/data/data/terraform/x86_64/nestos/libvirt/master.tf.template
	tfFileDir = "data/terraform/x86_64/nestos/libvirt"
	persistDir = "pkg/terraform"
	t.Run("ExecuteApplyTerraform", func(t *testing.T) {
		terraform, err := ExecuteApplyTerraform(tfFileDir, persistDir)
		if err != nil {
			t.Error(err)
		}
		t.Log(string(terraform))
	})

	t.Run("ExecuteDestroyTerraform", func(t *testing.T) {
		err := ExecuteDestroyTerraform(tfFileDir, persistDir)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("ExecuteDestroyTerraform,success")

		tmpFile := tfFileDir + "/terraform.tfstate"
		if err := os.Remove(tmpFile); err != nil {
			t.Errorf("Failed to remove terraform.tfstate: %v", err)
		}

		tmpDir := tfFileDir + "/.terraform"
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove %s: %v", tmpDir, err)
		}
	})

	t.Run("Outputs", func(t *testing.T) {
		by, err := Outputs(tfFileDir)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("Outputs,success", string(by))
	})
}
