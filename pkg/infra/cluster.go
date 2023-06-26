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
	"os"
	"path/filepath"

	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/assets"
	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/terraform"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Cluster struct {
	FileList []*assets.File
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
