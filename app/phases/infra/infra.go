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

package infra

import (
	"nestos-kubernetes-deployer/app/phases/infra/terraform"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Cluster struct {
	Dir  string // 由用户自定义，会在该路径下创建工作目录 .nkd
	Node string // 节点类别
	Num  string // 扩展节点个数
}

type File struct {
	Filename string
	Data     []byte
}

func (c *Cluster) Create() error {
	// 工作目录，包含terraform执行文件以及所需plugins
	workDir := "./"

	// tf配置文件所在目录
	tfDir := filepath.Join(workDir, c.Node)

	outputs, err := executeApplyTerraform(tfDir, workDir)
	if err != nil {
		return errors.Wrap(err, "failed to execute terraform apply")
	}
	logrus.Info(string(outputs))

	return nil
}

func executeApplyTerraform(tfDir string, terraformDir string) ([]byte, error) {
	var applyOpts []tfexec.ApplyOption
	return applyTerraform(tfDir, terraformDir, applyOpts...)
}

func applyTerraform(tfDir string, terraformDir string, applyOpts ...tfexec.ApplyOption) ([]byte, error) {
	applyErr := terraform.TFApply(tfDir, terraformDir, applyOpts...)
	if applyErr != nil {
		return nil, applyErr
	}

	_, err := os.Stat(filepath.Join(tfDir, "terraform.tfstate"))
	if os.IsNotExist(err) {
		return nil, err
	}

	outputs, err := terraform.Outputs(tfDir, terraformDir)
	if err != nil {
		return nil, err
	}

	return outputs, nil
}

func (c *Cluster) Destroy() error {
	// 工作目录，包含terraform执行文件以及所需plugins
	workDir := "./"

	// tf配置文件所在目录
	tfDir := filepath.Join(workDir, c.Node)

	err := executeDestroyTerraform(tfDir, workDir)
	if err != nil {
		return errors.Wrap(err, "failed to execute terraform destroy")
	}

	return nil
}

func executeDestroyTerraform(tfDir string, terraformDir string) error {
	var destroyOpts []tfexec.DestroyOption
	return destroyTerraform(tfDir, terraformDir, destroyOpts...)
}

func destroyTerraform(tfDir string, terraformDir string, destroyOpts ...tfexec.DestroyOption) error {
	destroyErr := terraform.TFDestroy(tfDir, terraformDir, destroyOpts...)
	if destroyErr != nil {
		return destroyErr
	}

	return nil
}

func (c *Cluster) Extend() error {
	return nil
}
