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

package cluster

import (
	"os"
	"path/filepath"

	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/assets"
	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/terraform"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/*
	installDir：自定义工作目录，在该路径下自动创建terraform文件夹
	dir：该路径存放tf配置文件
	terraformDir：该路径包含bin、plugins目录，存放terraform执行文件以及所需plugins
*/

type InfraProvider interface {
	Create()
	Destroy()
}

type Cluster struct {
	InstallDir string
	Platform   string
	Name       string
}

// TODO: 配置文件准备阶段

func (c *Cluster) Create() error {
	terraformDir := filepath.Join(c.InstallDir, "terraform")
	dir := filepath.Join(terraformDir, c.Platform, c.Name)

	terraformVariables := &TerraformVariables{}
	tfvarsFiles := make([]*assets.File, 0, len(terraformVariables.Files())+len(c.Platform)+len(c.Name))
	tfvarsFiles = append(tfvarsFiles, terraformVariables.Files()...)

	logrus.Infof("start to create %s in %s", c.Name, c.Platform)

	outputs, err := c.applyStage(dir, terraformDir, tfvarsFiles)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s in %s", c.Name, c.Platform)
	}

	logrus.Info(string(outputs.Data))
	logrus.Infof("succeed in creating %s in %s", c.Name, c.Platform)

	return nil
}

func (c *Cluster) applyStage(dir string, terraformDir string, tfvarsFiles []*assets.File) (*assets.File, error) {
	var applyOpts []tfexec.ApplyOption
	for _, file := range tfvarsFiles {
		if err := os.WriteFile(filepath.Join(dir, file.Filename), file.Data, 0o600); err != nil {
			return nil, err
		}
		applyOpts = append(applyOpts, tfexec.VarFile(filepath.Join(dir, file.Filename)))
	}

	return c.applyTerraform(dir, terraformDir, applyOpts...)
}

func (c *Cluster) applyTerraform(dir string, terraformDir string, applyOpts ...tfexec.ApplyOption) (*assets.File, error) {
	applyErr := terraform.TFApply(dir, terraformDir, applyOpts...)

	_, err := os.Stat(filepath.Join(dir, "terraform.tfstate"))
	if os.IsNotExist(err) {
		logrus.Errorf("Failed to read tfstate: %v", err)
		return nil, errors.Wrap(err, "failed to read tfstate")
	}

	if applyErr != nil {
		return nil, errors.WithMessage(applyErr, "failed to apply Terraform")
	}

	outputs, err := terraform.Outputs(dir, terraformDir)
	if err != nil {
		return nil, errors.Wrap(err, "could not get outputs file")
	}

	outputsFile := &assets.File{
		Filename: "outputs",
		Data:     outputs,
	}

	return outputsFile, nil
}

func (c *Cluster) Destroy() error {
	terraformDir := filepath.Join(c.InstallDir, "terraform")
	dir := filepath.Join(terraformDir, c.Platform, c.Name)

	logrus.Infof("start to destroy %s in %s", c.Name, c.Platform)

	// Question: Destroy的tfvarsFiles的获取

	terraformVariables := &TerraformVariables{}
	tfvarsFiles := make([]*assets.File, 0, len(terraformVariables.Files())+len(c.Platform)+len(c.Name))
	tfvarsFiles = append(tfvarsFiles, terraformVariables.Files()...)

	err := c.destroyStage(dir, terraformDir, tfvarsFiles)
	if err != nil {
		return errors.Wrapf(err, "failed to destroy %s in %s", c.Name, c.Platform)
	}
	os.Remove(dir)

	logrus.Infof("succeed in destroying %s in %s", c.Name, c.Platform)

	return nil
}

func (c *Cluster) destroyStage(dir string, terraformDir string, tfvarsFiles []*assets.File) error {
	var destroyOpts []tfexec.DestroyOption
	for _, file := range tfvarsFiles {
		if err := os.WriteFile(filepath.Join(dir, file.Filename), file.Data, 0o600); err != nil {
			return err
		}
		destroyOpts = append(destroyOpts, tfexec.VarFile(filepath.Join(dir, file.Filename)))
	}

	return destroyTerraform(dir, terraformDir, destroyOpts...)
}

func destroyTerraform(dir string, terraformDir string, destroyOpts ...tfexec.DestroyOption) error {
	destroyErr := terraform.TFDestroy(dir, terraformDir, destroyOpts...)
	if destroyErr != nil {
		return errors.WithMessage(destroyErr, "failed to destroy Terraform")
	}

	return nil
}
