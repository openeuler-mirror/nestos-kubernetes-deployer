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
	"nestos-kubernetes-deployer/app/cmd/phases/initconfig"
	"nestos-kubernetes-deployer/app/cmd/phases/workflow"
	"nestos-kubernetes-deployer/app/util/config"

	"github.com/spf13/cobra"
)

type initData struct {
	mastercfg *nkd.Master
	workercfg *nkd.Worker
}

func NewInitDefaultNkdConfigCommand() *cobra.Command {
	initRunner := workflow.NewRunner()
	var config string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Use this command to init cert, ign and tf config",
		RunE: func(cmd *cobra.Command, args []string) error {
			initRunner.SetDataInitializer(func(cmd *cobra.Command, args []string) (workflow.RunData, string, error) {
				data, nodetype, err := newInitData(config)
				if err != nil {
					return nil, "", err
				}
				return data, nodetype, nil
			})
			_, _, err := initRunner.InitData(args)
			if err != nil {
				return err
			}
			return initRunner.Run()
		},
	}
	// initRunner.AppendPhase(initconfig.NewGenerateCertsCmd())
	initRunner.AppendPhase(initconfig.NewGenerateIgnCmd())
	initRunner.AppendPhase(initconfig.NewGenerateTFCmd())

	cmd.PersistentFlags().StringVarP(&config, "config", "c", "", "config for init")
	return cmd
}

func (i *initData) MasterCfg() *nkd.Master {
	return i.mastercfg
}

func (i *initData) WorkerCfg() *nkd.Worker {
	return i.workercfg
}

func newInitData(cfgPath string) (*initData, string, error) {
	cfg, nodetype, err := config.LoadOrDefaultInitConfiguration(cfgPath)
	if err != nil {
		return nil, "", err
	}

	initData := &initData{}
	switch cfg := cfg.(type) {
	case *nkd.Master:
		initData.mastercfg = cfg
	case *nkd.Worker:
		initData.workercfg = cfg
	default:
		return nil, "", fmt.Errorf("please provide the path of the cluster node config, master or worker")
	}

	return initData, nodetype, nil
}
