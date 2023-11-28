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
	"fmt"
	"nestos-kubernetes-deployer/pkg/infra/terraform"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Cluster struct {
	Node string // 节点类别
	Num  int    // 扩展节点个数
}

func (c *Cluster) Deploy() error {
	// 工作目录，包含terraform执行文件以及所需plugins
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// tf配置文件所在目录
	tfDir := filepath.Join(workDir, c.Node)

	outputs, err := terraform.ExecuteApplyTerraform(tfDir, workDir)
	if err != nil {
		return errors.Wrap(err, "failed to execute terraform apply")
	}
	fmt.Println(string(outputs))

	return nil
}

func (c *Cluster) Extend() error {
	// 工作目录，包含terraform执行文件以及所需plugins
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// tf配置文件所在目录
	tfDir := filepath.Join(workDir, c.Node)

	outputs, err := terraform.ExtendTerraform(tfDir, workDir, c.Num)
	if err != nil {
		return errors.Wrap(err, "failed to execute terraform apply")
	}
	fmt.Println(string(outputs))

	return nil
}

func (c *Cluster) Destroy() error {
	// 工作目录，包含terraform执行文件以及所需plugins
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// tf配置文件所在目录
	tfDir := filepath.Join(workDir, c.Node)

	err = terraform.ExecuteDestroyTerraform(tfDir, workDir)
	if err != nil {
		return errors.Wrap(err, "failed to execute terraform destroy")
	}

	return nil
}
