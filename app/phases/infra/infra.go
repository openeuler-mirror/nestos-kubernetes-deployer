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
	"nestos-kubernetes-deployer/app/apis/nkd"
	"nestos-kubernetes-deployer/app/phases/infra/terraform"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Cluster struct {
	Dir  string // 由用户自定义，会在该路径下创建工作目录 /nkd
	Role string // 节点类别
}

type File struct {
	Filename string
	Data     []byte
}

func (c *Cluster) Create() error {
	// 若未指定c.Dir，则默认为当前路径
	if c.Dir == "" {
		c.Dir = "./"
	}

	if c.Role == "" {
		return errors.Errorf("please provide the role of the node, master or worker")
	}

	// 工作目录，包含terraform执行文件以及所需plugins
	workDir := filepath.Join(c.Dir, "nkd")
	// 从文件中读取 infraData
	infraData, err := os.ReadFile(filepath.Join(workDir, "config.yaml"))
	if err != nil {
		return errors.WithMessage(err, "failed to read yaml file")
	}

	// 解析 yaml 数据
	var configData nkd.Nkd
	err = yaml.Unmarshal(infraData, &configData)
	if err != nil {
		return errors.WithMessage(err, "failed to parse yaml data")
	}

	// 获取当前选择的平台
	platform := configData.Infra.Platform

	// tf配置文件所在目录
	tfDir := filepath.Join(workDir, platform, c.Role)

	// ToDo: 支持terraform apply参数
	applyOptsFile := &File{
		Filename: ".applyOptsFile",
		Data:     []byte{},
	}

	logrus.Infof("start to create %s in %s", c.Role, platform)

	outputs, err := executeApplyTerraform(tfDir, workDir, applyOptsFile)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s in %s", c.Role, platform)
	}

	logrus.Info(string(outputs.Data))
	logrus.Infof("succeed in creating %s in %s", c.Role, platform)

	return nil
}

func executeApplyTerraform(tfDir string, terraformDir string, tfvarsFile *File) (*File, error) {
	var applyOpts []tfexec.ApplyOption

	if err := os.WriteFile(filepath.Join(tfDir, tfvarsFile.Filename), tfvarsFile.Data, 0o600); err != nil {
		return nil, err
	}
	applyOpts = append(applyOpts, tfexec.VarFile(filepath.Join(tfDir, tfvarsFile.Filename)))

	return applyTerraform(tfDir, terraformDir, applyOpts...)
}

func applyTerraform(tfDir string, terraformDir string, applyOpts ...tfexec.ApplyOption) (*File, error) {
	applyErr := terraform.TFApply(tfDir, terraformDir, applyOpts...)

	_, err := os.Stat(filepath.Join(tfDir, "terraform.tfstate"))
	if os.IsNotExist(err) {
		return nil, errors.Wrap(err, "failed to read tfstate")
	}

	if applyErr != nil {
		return nil, errors.Wrap(applyErr, "failed to apply Terraform")
	}

	outputs, err := terraform.Outputs(tfDir, terraformDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get outputs file")
	}

	outputsFile := &File{
		Filename: "tfOutputs",
		Data:     outputs,
	}

	return outputsFile, nil
}

func (c *Cluster) Destroy() error {
	// 若未指定c.Dir，则默认为当前路径
	if c.Dir == "" {
		c.Dir = "./"
	}

	if c.Role == "" {
		return errors.Errorf("please provide the role of the node, master or worker")
	}

	// 从文件中读取 infraData
	workDir := filepath.Join(c.Dir, "nkd")
	infraData, err := os.ReadFile(filepath.Join(workDir, "config.yaml"))
	if err != nil {
		return errors.Wrap(err, "failed to read yaml data")
	}

	// 解析 yaml 数据
	var configData nkd.Nkd
	err = yaml.Unmarshal(infraData, &configData)
	if err != nil {
		return errors.Wrap(err, "failed to parse yaml data")
	}

	// 获取当前选择的平台
	platform := configData.Infra.Platform
	tfDir := filepath.Join(workDir, platform, c.Role)

	logrus.Infof("start to destroy %s in %s", c.Role, platform)

	// ToDo: 支持terraform destroy参数
	destroyOptsFile := &File{
		Filename: ".destroyOptsFile",
		Data:     []byte{},
	}

	err = executeDestroyTerraform(tfDir, workDir, destroyOptsFile)
	if err != nil {
		return errors.Wrapf(err, "failed to destroy %s in %s", c.Role, platform)
	}
	os.RemoveAll(tfDir)

	logrus.Infof("succeed in destroying %s in %s", c.Role, platform)

	return nil
}

func executeDestroyTerraform(tfDir string, terraformDir string, tfvarsFile *File) error {
	var destroyOpts []tfexec.DestroyOption
	destroyOpts = append(destroyOpts, tfexec.VarFile(filepath.Join(tfDir, tfvarsFile.Filename)))

	return destroyTerraform(tfDir, terraformDir, destroyOpts...)
}

func destroyTerraform(tfDir string, terraformDir string, destroyOpts ...tfexec.DestroyOption) error {
	destroyErr := terraform.TFDestroy(tfDir, terraformDir, destroyOpts...)
	if destroyErr != nil {
		return errors.Wrap(destroyErr, "failed to destroy Terraform")
	}

	return nil
}
