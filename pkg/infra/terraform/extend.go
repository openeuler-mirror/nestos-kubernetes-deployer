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
	"os"
	"path/filepath"
)

func ExtendTerraform(tfDir string, terraformDir string, num int) ([]byte, error) {
	applyErr := TFExtend(tfDir, terraformDir, num)
	if applyErr != nil {
		return nil, applyErr
	}

	_, err := os.Stat(filepath.Join(tfDir, "terraform.tfstate"))
	if os.IsNotExist(err) {
		return nil, err
	}

	outputs, err := Outputs(tfDir, terraformDir)
	if err != nil {
		return nil, err
	}

	return outputs, nil
}
