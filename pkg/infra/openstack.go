/*
Copyright 2024 KylinSoft  Co., Ltd.

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
	"nestos-kubernetes-deployer/pkg/terraform"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type OpenStack struct {
	PersistDir string
	ClusterID  string
	Node       string
	Count      uint
}

func (o *OpenStack) Deploy() error {
	tfFileDir := filepath.Join(o.PersistDir, o.ClusterID, o.Node)
	outputs, err := terraform.ExecuteApplyTerraform(tfFileDir, o.PersistDir)
	if err != nil {
		return errors.Wrap(err, "failed to execute terraform apply")
	}
	logrus.Println(string(outputs))

	return nil
}

func (o *OpenStack) Extend() error {
	tfFileDir := filepath.Join(o.PersistDir, o.ClusterID, o.Node)
	outputs, err := terraform.ExecuteApplyTerraform(tfFileDir, o.PersistDir)
	if err != nil {
		return errors.Wrap(err, "failed to execute terraform apply")
	}
	logrus.Println(string(outputs))

	return nil
}

func (o *OpenStack) Destroy() error {
	tfFileDir := filepath.Join(o.PersistDir, o.ClusterID, o.Node)
	err := terraform.ExecuteDestroyTerraform(tfFileDir, o.PersistDir)
	if err != nil {
		return errors.Wrap(err, "failed to execute terraform destroy")
	}

	return nil
}
