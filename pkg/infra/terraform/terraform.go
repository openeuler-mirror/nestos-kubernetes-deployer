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
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/*
	tfFileDir: tf file directory
	persistDir: nkd config directory, default /etc/nkd
*/

const execPath string = "/usr/bin/tofu"

func newTFExec(tfFileDir string) (*tfexec.Terraform, error) {
	tf, err := tfexec.NewTerraform(tfFileDir, execPath)
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

	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)

	tf.SetLogger(newPrintfer())

	return tf, nil
}

// terraform init
func TFInit(tfFileDir string, persistDir string) (err error) {
	tf, err := newTFExec(tfFileDir)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec")
	}

	// Try to perform initialization using the local plug-in directory.
	err = tf.Init(context.Background(), tfexec.PluginDir(filepath.Join(persistDir, "providers")))
	if err == nil {
		return nil
	}

	fmt.Print("Failed to initialize Terraform with existed plugin directory\nStart downloading plugins...\n")
	// Set the path for downloading plug-ins.
	os.Setenv("TF_DATA_DIR", persistDir)
	err = tf.Init(context.Background(), tfexec.Upgrade(false))
	if err != nil {
		return errors.Wrap(err, "failed to init terraform")
	}

	return nil
}

// terraform apply
func TFApply(tfFileDir string, persistDir string, applyOpts ...tfexec.ApplyOption) error {
	if err := TFInit(tfFileDir, persistDir); err != nil {
		return errors.Wrap(err, "failed to init terraform")
	}

	tf, err := newTFExec(tfFileDir)
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
func TFDestroy(tfFileDir string, persistDir string, destroyOpts ...tfexec.DestroyOption) error {
	if err := TFInit(tfFileDir, persistDir); err != nil {
		return errors.Wrap(err, "failed to init terraform")
	}

	tf, err := newTFExec(tfFileDir)
	if err != nil {
		return errors.Wrap(err, "failed to destroy a new tfexec")
	}

	err = tf.Destroy(context.Background(), destroyOpts...)
	if err != nil {
		return errors.Wrap(err, "failed to destroy terraform")
	}

	return nil
}
