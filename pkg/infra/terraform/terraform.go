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
	"github.com/openshift/installer/pkg/lineprinter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func newTFExec(workingDir string, terraformBinary string) (*tfexec.Terraform, error) {
	execPath := filepath.Join(terraformBinary, "bin", "terraform")
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		return nil, err
	}

	// If the log path is not set, terraform will not receive debug logs.
	if path, ok := os.LookupEnv("TEREFORM_LOG_PATH"); ok {
		if err := tf.SetLog(os.Getenv("TEREFORM_LOG")); err != nil {
			logrus.Infof("Skipping setting terraform log levels: %v", err)
		} else {
			tf.SetLogCore(os.Getenv("TEREFORM_LOG_CORE"))         //nolint:errcheck
			tf.SetLogProvider(os.Getenv("TEREFORM_LOG_PROVIDER")) //nolint:errcheck
			tf.SetLogPath(path)                                   //nolint:errcheck
		}
	}

	// Add terraform info logs to the installer log
	lpPrint := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Print}).Print}
	lpError := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Error}).Print}
	defer lpPrint.Close()
	defer lpError.Close()

	tf.SetStdout(lpPrint)
	tf.SetStderr(lpError)
	tf.SetLogger(newPrintfer())

	return tf, nil
}

// terraform apply
func TFApply(workingDir string, platform string, stage Stage, terraformBinary string, applyOpts ...tfexec.ApplyOption) error {
	if err := TFInit(workingDir, platform, stage.Name(), terraformBinary, stage.Providers()); err != nil {
		return err
	}

	tf, err := newTFExec(workingDir, terraformBinary)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec.")
	}

	err = tf.Apply(context.Background(), applyOpts...)
	return errors.Wrap(err, "failed to apply Terraform.")
}

// terraform destroy
func TFDestroy(workingDir string, platform string, stage Stage, terraformBinary string, destroyOpts ...tfexec.DestroyOption) error {
	if err := TFInit(workingDir, platform, stage.Name(), terraformBinary, stage.Providers()); err != nil {
		return err
	}

	tf, err := newTFExec(workingDir, terraformBinary)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec.")
	}

	return errors.Wrap(
		tf.Destroy(context.Background(), destroyOpts...),
		"failed doing terraform destroy.",
	)
}
