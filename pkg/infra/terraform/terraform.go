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
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
)

func newTFExec(tfFilePath string, terraformBinary string) (*tfexec.Terraform, error) {
	tf, err := tfexec.NewTerraform(tfFilePath, terraformBinary)
	if err != nil {
		return nil, err
	}

	// TODO 日志输出
	// tf.SetStdout()
	// tf.SetStderr()
	// tf.SetLogger(newPrintfer())

	return tf, nil
}

func Init(tfFilePath string, terraformBinary string) error {
	tf, err := newTFExec(tfFilePath, terraformBinary)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec.")
	}

	// 如果想导入本地已有插件，需手动构建相应目录
	// 如openstack插件所需目录为plugins/registry.terraform.io/terraform-provider-openstack/openstack/1.51.1/linux_arm64/terraform-provider-openstack_v1.51.1
	return errors.Wrap(
		tf.Init(context.Background(), tfexec.PluginDir(filepath.Join(terraformBinary, "plugins"))),
		"failed doing terraform init.",
	)
}

func Apply(tfFilePath string, terraformBinary string, extraOpts ...tfexec.ApplyOption) error {
	tf, err := newTFExec(tfFilePath, terraformBinary)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec.")
	}

	// TODO 日志输出
	// tf.SetStdout()
	// tf.SetStderr()
	tf.SetLogger(newPrintfer())

	err = tf.Apply(context.Background(), extraOpts...)
	return errors.Wrap(err, "failed to apply Terraform.")
}

func Destroy(tfFilePath string, terraformBinary string, extraOpts ...tfexec.DestroyOption) error {
	tf, err := newTFExec(tfFilePath, terraformBinary)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec.")
	}

	// TODO 日志输出
	// tf.SetStdout()
	// tf.SetStderr()
	tf.SetLogger(newPrintfer())

	err = tf.Destroy(context.Background(), extraOpts...)
	return errors.Wrap(err, "failed to apply Terraform.")
}
