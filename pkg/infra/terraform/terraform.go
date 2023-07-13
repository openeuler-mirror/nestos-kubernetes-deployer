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
	"runtime"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/openshift/installer/pkg/lineprinter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func newTFExec(dir string, terraformDir string) (*tfexec.Terraform, error) {
	tfPath := filepath.Join(terraformDir, "bin", runtime.GOOS+"_"+runtime.GOARCH, "terraform")
	tf, err := tfexec.NewTerraform(dir, tfPath)
	if err != nil {
		return nil, err
	}

	// If the log path is not set, terraform will not receive debug logs.
	if logPath, ok := os.LookupEnv("TERRAFORM_LOG_PATH"); ok {
		if err := tf.SetLog(os.Getenv("TERRAFORM_LOG")); err != nil {
			logrus.Infof("Skipping setting terraform log levels: %v", err)
		} else {
			tf.SetLogCore(os.Getenv("TERRAFORM_LOG_CORE"))         //nolint:errcheck
			tf.SetLogProvider(os.Getenv("TERRAFORM_LOG_PROVIDER")) //nolint:errcheck
			tf.SetLogPath(logPath)                                 //nolint:errcheck
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
func TFApply(dir string, terraformDir string, applyOpts ...tfexec.ApplyOption) error {
	if err := TFInit(dir, terraformDir); err != nil {
		return err
	}

	tf, err := newTFExec(dir, terraformDir)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec")
	}

	err = tf.Apply(context.Background(), applyOpts...)
	if err != nil {
		return errors.Wrap(err, "failed to apply Terraform")
	}

	return nil
}

// terraform destroy
func TFDestroy(dir string, terraformDir string, destroyOpts ...tfexec.DestroyOption) error {
	if err := TFInit(dir, terraformDir); err != nil {
		return err
	}

	tf, err := newTFExec(dir, terraformDir)
	if err != nil {
		return errors.Wrap(err, "failed to destroy a new tfexec")
	}

	err = tf.Destroy(context.Background(), destroyOpts...)
	if err != nil {
		return errors.Wrap(err, "failed to destroy terraform")
	}

	return nil
}
