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
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
)

func newTFExec(workingDir string, terraformBinary string) (*tfexec.Terraform, error) {
	execPath := filepath.Join(terraformBinary, "bin", "terraform")
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		return nil, err
	}

	// TODO 日志打印
	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)
	tf.SetLogger(newPrintfer())

	return tf, nil
}

// terraform apply
func TFApply(workingDir string, platform string, stage Stage, terraformBinary string, applyOpts ...tfexec.ApplyOption) error {
	if err := tfInit(workingDir, platform, stage.Name(), terraformBinary, stage.Providers()); err != nil {
		return err
	}

	tf, err := newTFExec(workingDir, terraformBinary)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec.")
	}

	// TODO 日志打印
	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)
	tf.SetLogger(newPrintfer())

	err = tf.Apply(context.Background(), applyOpts...)
	return errors.Wrap(err, "failed to apply Terraform.")
}

// terraform destroy
func TFDestroy(workingDir string, platform string, stage Stage, terraformBinary string, destroyOpts ...tfexec.DestroyOption) error {
	if err := tfInit(workingDir, platform, stage.Name(), terraformBinary, stage.Providers()); err != nil {
		return err
	}

	tf, err := newTFExec(workingDir, terraformBinary)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec.")
	}

	// TODO 日志打印
	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)
	tf.SetLogger(newPrintfer())

	return errors.Wrap(
		tf.Destroy(context.Background(), destroyOpts...),
		"failed doing terraform destroy.",
	)
}
