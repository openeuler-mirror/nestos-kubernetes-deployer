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
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"

	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/assets"
	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/tfvars/openstack"
	"github.com/pkg/errors"
)

type TerraformVariables struct {
	FileList []*assets.File
}

func (t *TerraformVariables) Files() []*assets.File {
	return t.FileList
}

func Generate() error {
	// 从文件中读取 jsonData
	jsonData, err := os.ReadFile("data/json/openstack/main.tf.json")
	if err != nil {
		return errors.Wrap(err, "error reading json data")
	}

	// 解析 JSON 数据
	var terraformData openstack.TerraformData
	err = json.Unmarshal(jsonData, &terraformData)
	if err != nil {
		return errors.Wrap(err, "error parsing json data")
	}

	// 从文件中读取 terraformConfig
	terraformConfig, err := os.ReadFile("data/templates/openstack/main.tf.template")
	if err != nil {
		return errors.Wrap(err, "error reading terraform config template")
	}

	// 使用模板填充数据
	tmpl, err := template.New("terraform").Parse(string(terraformConfig))
	if err != nil {
		return errors.Wrap(err, "error creating terraform config template")
	}

	// 创建一个新的文件用于写入填充后的数据
	tfDir := filepath.Join("/root", "terraform")
	if err := os.MkdirAll(tfDir, os.ModePerm); err != nil {
		return errors.Wrap(err, "could not create the terraform directory")
	}

	outputFile, err := os.Create(filepath.Join(tfDir, "main.tf"))
	if err != nil {
		return errors.Wrap(err, "error creating terraform config")
	}
	defer outputFile.Close()

	// 将填充后的数据写入文件
	err = tmpl.Execute(outputFile, terraformData)
	if err != nil {
		return errors.Wrap(err, "error executing terraform config")
	}

	return nil
}
