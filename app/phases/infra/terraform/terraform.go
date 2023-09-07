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
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/openshift/installer/pkg/lineprinter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/*
	tfDir: tf配置文件所在目录
	terraformDir: terraform执行文件所在目录
*/

func newTFExec(tfDir string, terraformDir string) (*tfexec.Terraform, error) {
	tfPath := filepath.Join(terraformDir, "terraform")
	tf, err := tfexec.NewTerraform(tfDir, tfPath)
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
	lpDebug := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Debug}).Print}
	lpError := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Error}).Print}
	defer lpDebug.Close()
	defer lpError.Close()

	tf.SetStdout(lpDebug)
	tf.SetStderr(lpError)
	tf.SetLogger(newPrintfer())

	return tf, nil
}

// terraform init
func TFInit(tfDir string, terraformDir string) (err error) {
	tf, err := newTFExec(tfDir, terraformDir)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec")
	}

	// 使用本地terraform插件
	err = tf.Init(context.Background(), tfexec.PluginDir(filepath.Join(terraformDir, ".terraform/providers")))
	if err != nil {
		return errors.Wrap(err, "failed to init terraform")
	}

	return nil
}

// terraform apply
func TFApply(tfDir string, terraformDir string, applyOpts ...tfexec.ApplyOption) error {
	if err := TFInit(tfDir, terraformDir); err != nil {
		return errors.Wrap(err, "failed to init terraform")
	}

	tf, err := newTFExec(tfDir, terraformDir)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec")
	}

	err = tf.Apply(context.Background(), applyOpts...)
	if err != nil {
		return errors.Wrap(err, "failed to apply Terraform")
	}

	return nil
}

// terraform apply for nkd extend
func TFExtend(tfDir string, terraformDir string, num int) error {
	if err := TFInit(tfDir, terraformDir); err != nil {
		return errors.Wrap(err, "failed to init terraform")
	}

	tf, err := newTFExec(tfDir, terraformDir)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec")
	}

	err = tf.Apply(context.Background(), tfexec.Var(fmt.Sprintf("instance_count=%d", num)))
	if err != nil {
		return errors.Wrap(err, "failed to extend Terraform")
	}

	return nil
}

// terraform destroy
func TFDestroy(tfDir string, terraformDir string, destroyOpts ...tfexec.DestroyOption) error {
	if err := TFInit(tfDir, terraformDir); err != nil {
		return errors.Wrap(err, "failed to init terraform")
	}

	tf, err := newTFExec(tfDir, terraformDir)
	if err != nil {
		return errors.Wrap(err, "failed to destroy a new tfexec")
	}

	err = tf.Destroy(context.Background(), destroyOpts...)
	if err != nil {
		return errors.Wrap(err, "failed to destroy terraform")
	}

	return nil
}
