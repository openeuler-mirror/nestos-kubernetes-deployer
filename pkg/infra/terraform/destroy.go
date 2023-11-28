/*
Copyright 2023 KylinSoft  Co., Ltd.

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
	"github.com/hashicorp/terraform-exec/tfexec"
)

func ExecuteDestroyTerraform(tfDir string, terraformDir string) error {
	var destroyOpts []tfexec.DestroyOption
	return destroyTerraform(tfDir, terraformDir, destroyOpts...)
}

func destroyTerraform(tfDir string, terraformDir string, destroyOpts ...tfexec.DestroyOption) error {
	destroyErr := TFDestroy(tfDir, terraformDir, destroyOpts...)
	if destroyErr != nil {
		return destroyErr
	}

	return nil
}
