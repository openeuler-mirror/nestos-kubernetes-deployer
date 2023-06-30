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
	platformStages "gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/terraform/stages/platform"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Cluster struct {
	FileList []*assets.File
}

func (c *Cluster) Create(platform string, node string, installDir string) (err error) {
	if installDir == "" {
		logrus.Fatal("InstallDir has not been set")
	}

	terraformVariables := &TerraformVariables{}
	stages, err := platformStages.StagesForPlatform(platform, node)

	terraformBinary := filepath.Join(installDir, "terraform")
	_, err = os.Stat(terraformBinary)
	if os.IsNotExist(err) {
		if err := os.Mkdir(terraformBinary, 0777); err != nil {
			return errors.Wrap(err, "could not create the terraform directory")
		}
	}

	// TODO 将所需二进制文件保存到terraformBinary目录

	terraformBinaryPath, err := filepath.Abs(terraformBinary)
	if err != nil {
		return errors.Wrap(err, "cannot get absolute path of terraform directory")
	}

	tfvarsFiles := make([]*assets.File, 0, len(terraformVariables.Files())+len(stages))
	tfvarsFiles = append(tfvarsFiles, terraformVariables.Files()...)

	for _, stage := range stages {
		outputs, err := c.applyStage(platform, stage, terraformBinaryPath, tfvarsFiles)
		if err != nil {
			return errors.Wrapf(err, "failure applying terraform for %q stage", stage.Name())
		}
		tfvarsFiles = append(tfvarsFiles, outputs)
		c.FileList = append(c.FileList, outputs)
	}

	return nil
}

func (c *Cluster) applyStage(platform string, stage terraform.Stage, terraformBinary string, tfvarsFiles []*assets.File) (*assets.File, error) {
	tmpDir := filepath.Join(terraformBinary, platform, stage.Name())
	if err := os.MkdirAll(tmpDir, 0777); err != nil {
		return nil, errors.Wrapf(err, "cannot create the directory for %s", stage.Name())
	}

	var applyOpts []tfexec.ApplyOption
	for _, file := range tfvarsFiles {
		if err := os.WriteFile(filepath.Join(tmpDir, file.Filename), file.Data, 0o600); err != nil {
			return nil, err
		}
		applyOpts = append(applyOpts, tfexec.VarFile(filepath.Join(tmpDir, file.Filename)))
	}

	return c.applyTerraform(tmpDir, platform, stage, terraformBinary, applyOpts...)
}

func (c *Cluster) applyTerraform(tmpDir string, platform string, stage terraform.Stage, terraformBinary string, applyOpts ...tfexec.ApplyOption) (*assets.File, error) {
	applyErr := terraform.TFApply(tmpDir, platform, stage, terraformBinary, applyOpts...)

	if data, err := os.ReadFile(filepath.Join(tmpDir, terraform.StateFilename)); err == nil {
		c.FileList = append(c.FileList, &assets.File{
			Filename: stage.StateFilename(),
			Data:     data,
		})
	} else if !os.IsNotExist(err) {
		logrus.Errorf("Failed to read tfstate: %v", err)
		return nil, errors.Wrap(err, "failed to read tfstate")
	}

	if applyErr != nil {
		return nil, errors.WithMessage(applyErr, "failed to apply Terraform.")
	}

	outputs, err := terraform.Outputs(tmpDir, terraformBinary)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get outputs from stage %q", stage.Name())
	}

	outputsFile := &assets.File{
		Filename: stage.OutputsFilename(),
		Data:     outputs,
	}

	return outputsFile, nil
}
