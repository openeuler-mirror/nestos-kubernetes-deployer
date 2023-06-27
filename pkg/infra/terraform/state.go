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
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

const StateFilename = "terraform.tfstate"

// Reads the terraform state file.
func Outputs(workingDir string, terraformBinary string) ([]byte, error) {
	tf, err := newTFExec(workingDir, terraformBinary)
	if err != nil {
		return nil, err
	}

	tfoutput, err := tf.Output(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to read terraform state file")
	}

	outputs := make(map[string]interface{}, len(tfoutput))
	for key, value := range tfoutput {
		outputs[key] = value.Value
	}

	data, err := json.Marshal(outputs)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal outputs")
	}

	return data, nil
}
