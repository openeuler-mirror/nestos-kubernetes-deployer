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
	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/terraform"
	"github.com/hashicorp/terraform-exec/tfexec"
)

type OpenstackDeployer struct {
	tfFilePath      string
	terraformBinary string
}

func (t OpenstackDeployer) Create(extraOpts ...tfexec.ApplyOption) error {
	if err := terraform.Init(t.tfFilePath, t.terraformBinary); err != nil {
		return err
	}

	applyErr := terraform.Apply(t.tfFilePath, t.terraformBinary, extraOpts...)
	if applyErr != nil {
		return applyErr
	}

	return nil
}

func (t OpenstackDeployer) Destroy(extraOpts ...tfexec.DestroyOption) error {
	if err := terraform.Init(t.tfFilePath, t.terraformBinary); err != nil {
		return err
	}

	destroyErr := terraform.Destroy(t.tfFilePath, t.terraformBinary, extraOpts)
	if destroyErr != nil {
		return destroyErr
	}

	return nil
}
