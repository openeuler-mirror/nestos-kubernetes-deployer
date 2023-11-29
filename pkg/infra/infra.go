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
	"nestos-kubernetes-deployer/pkg/infra/terraform"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Cluster struct {
	PersistDir string
	ClusterID  string
	Node       string
	Count      int
}

func (c *Cluster) Deploy() (err error) {
	tfFileDir := filepath.Join(c.PersistDir, c.ClusterID, c.Node)
	outputs, err := terraform.ExecuteApplyTerraform(tfFileDir, c.PersistDir)
	if err != nil {
		return errors.Wrap(err, "failed to execute terraform apply")
	}
	logrus.Println(string(outputs))

	return nil
}

func (c *Cluster) Extend() (err error) {
	tfFileDir := filepath.Join(c.PersistDir, c.ClusterID, c.Node)
	outputs, err := terraform.ExtendTerraform(tfFileDir, c.PersistDir, c.Count)
	if err != nil {
		return errors.Wrap(err, "failed to execute terraform apply")
	}
	logrus.Println(string(outputs))

	return nil
}

func (c *Cluster) Destroy() (err error) {
	// tf file directory.
	tfFileDir := filepath.Join(c.PersistDir, c.ClusterID, c.Node)
	err = terraform.ExecuteDestroyTerraform(tfFileDir, c.PersistDir)
	if err != nil {
		return errors.Wrap(err, "failed to execute terraform destroy")
	}

	return nil
}
