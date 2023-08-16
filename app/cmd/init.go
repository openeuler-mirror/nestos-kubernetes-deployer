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
	"fmt"

	"nestos-kubernetes-deployer/app/apis/nkd"
	phases "nestos-kubernetes-deployer/app/cmd/phases/init"
	"nestos-kubernetes-deployer/app/cmd/phases/workflow"
	"nestos-kubernetes-deployer/app/util/config"

	"github.com/spf13/cobra"
)

type initData struct {
	cfg *nkd.Nkd
}

func NewInitCommand() *cobra.Command {
	initRunner := workflow.NewRunner()
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Use this command to init insert or config",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := initRunner.InitData(args)
			if err != nil {
				return err
			}
			data := c.(*initData)
			fmt.Println(data.cfg)
			return initRunner.Run()
		},
	}

	phases.NewGenerateCertsCmd()
	initRunner.AppendPhase(phases.NewGenerateCertsCmd())
	initRunner.AppendPhase(phases.NewGenerateIgnCmd())
	initRunner.SetDataInitializer(func(cmd *cobra.Command, args []string) (workflow.RunData, error) {
		data, err := newInitData(cmd, args)
		if err != nil {
			return nil, err
		}
		return data, nil
	})
	return cmd
}

func (i *initData) Cfg() *nkd.Nkd {
	return i.cfg
}
func newInitData(cmd *cobra.Command, args []string) (*initData, error) {

	var newNkd *nkd.Nkd
	cfg, err := config.LoadOrDefaultInitConfiguration("path", newNkd)
	if err != nil {
		return nil, err
	}
	return &initData{
		cfg: cfg,
	}, nil
}
