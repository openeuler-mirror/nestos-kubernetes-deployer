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

package phases

import (
	"fmt"
	"nestos-kubernetes-deployer/app/cmd/phases/workflow"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

func NewGenerateTFCmd() workflow.Phase {
	return workflow.Phase{
		Name:  "tf",
		Short: "Run tf to generate certs",
		Run:   runGenerateTFConfig,
	}
}

func runGenerateTFConfig(r workflow.RunData, node string) error {
	outputFile, err := os.Create(filepath.Join("./", fmt.Sprintf("%s.tf", node)))
	if err != nil {
		return errors.Wrap(err, "failed to create terraform config file")
	}
	defer outputFile.Close()

	// 从文件中读取 terraformConfig
	terraformConfig, err := os.ReadFile(filepath.Join("resource/templates/terraform", fmt.Sprintf("%s.tf.tpl", node)))
	if err != nil {
		return errors.Wrap(err, "failed to read terraform config template")
	}

	// 使用模板填充数据
	tmpl, err := template.New("terraform").Parse(string(terraformConfig))
	if err != nil {
		return errors.Wrap(err, "failed to create terraform config template")
	}

	quotedStrs := make([]string, len(r.(InitData).MasterCfg().System.Ips))
	for i, s := range r.(InitData).MasterCfg().System.Ips {
		quotedStrs[i] = fmt.Sprintf(`"%s"`, s)
	}
	joinedStr := strings.Join(quotedStrs, ",")

	r.(InitData).MasterCfg().System.Ips = []string{joinedStr}

	// 将填充后的数据写入文件
	if node == "master" {
		err = tmpl.Execute(outputFile, r.(InitData).MasterCfg())
	} else {
		err = tmpl.Execute(outputFile, r.(InitData).WorkerCfg())
	}
	if err != nil {
		return errors.Wrap(err, "failed to write terraform config")
	}

	return nil
}
