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

package cmd

import (
	"io/ioutil"
	"nestos-kubernetes-deployer/cmd/command"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/utils"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func NewTemplateCommand() *cobra.Command {
	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Create a default template of nkd config",
		RunE:  createTemplate,
	}

	command.SetupTemplateCmdOpts(templateCmd)

	return templateCmd
}

func createTemplate(cmd *cobra.Command, args []string) error {
	arch := runtime.GOARCH
	if opts.Opts.Arch != "" {
		arch = opts.Opts.Arch
	}
	conf, err := asset.GetDefaultClusterConfig(arch)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(conf)
	if err != nil {
		logrus.Errorf("Faild to marshal template config: %v", err)
		return err
	}
	file, err := cmd.Flags().GetString("file")
	if err != nil {
		logrus.Errorf("Failed to create template file: %v", err)
		return err
	}
	if file == "" {
		file = "./template.yaml"
	}
	if err := ioutil.WriteFile(file, data, utils.DeployConfigFileMode); err != nil {
		logrus.Errorf("Faild to write template config file: %v", err)
		return err
	}

	return nil
}
