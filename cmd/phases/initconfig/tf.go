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

package initconfig

import (
	"fmt"
	"io"
	"nestos-kubernetes-deployer/data"
	"nestos-kubernetes-deployercmd/phases/workflow"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

func NewGenerateTFCmd() workflow.Phase {
	return workflow.Phase{
		Name:  "tf",
		Short: "Run tf to generate .tf config files",
		Run:   runGenerateTFConfig,
	}
}

func runGenerateTFConfig(r workflow.RunData, node string) error {
	outputFile, err := os.Create(filepath.Join(node, fmt.Sprintf("%s.tf", node)))
	if err != nil {
		return errors.Wrap(err, "failed to create terraform config file")
	}
	defer outputFile.Close()

	tfFilePath := filepath.Join("terraform", fmt.Sprintf("%s.tf.template", node))
	tfFile, err := data.Assets.Open(tfFilePath)
	if err != nil {
		return err
	}
	defer tfFile.Close()
	data, err := io.ReadAll(tfFile)
	if err != nil {
		return err
	}
	tmpl, err := template.New("terraform").Parse(string(data))
	if err != nil {
		return errors.Wrap(err, "failed to create terraform config template")
	}

	// 将填充后的数据写入文件
	if node == "master" {
		quotedStrs := make([]string, len(r.(InitData).MasterCfg().System.MasterIps))
		for i, s := range r.(InitData).MasterCfg().System.MasterIps {
			quotedStrs[i] = fmt.Sprintf(`"%s"`, s)
		}
		joinedStr := strings.Join(quotedStrs, ",")

		r.(InitData).MasterCfg().System.MasterIps = []string{joinedStr}

		err = tmpl.Execute(outputFile, r.(InitData).MasterCfg())
	} else {
		err = tmpl.Execute(outputFile, r.(InitData).WorkerCfg())
	}

	if err != nil {
		return errors.Wrap(err, "failed to write terraform config")
	}

	return nil
}
