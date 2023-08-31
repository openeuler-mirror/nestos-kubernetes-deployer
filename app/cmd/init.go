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
	"nestos-kubernetes-deployer/app/apis/nkd"
	phases "nestos-kubernetes-deployer/app/cmd/phases/init"
	"nestos-kubernetes-deployer/app/cmd/phases/workflow"
	"nestos-kubernetes-deployer/app/util/config"

	"github.com/spf13/cobra"
)

type initData struct {
	mastercfg *nkd.Master
	workercfg *nkd.Worker
	// cfg *nkd.Master
	// cfg *interface{}
}

type Config struct {
	config string
}

func NewInitDefaultNkdConfigCommand() *cobra.Command {
	initRunner := workflow.NewRunner()
	var config string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Use this command to init ign, cert config",
		RunE: func(cmd *cobra.Command, args []string) error {
			initRunner.SetDataInitializer(func(cmd *cobra.Command, args []string) (workflow.RunData, string, error) {
				data, nodetype, err := newInitData(cmd, args, config)
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
	phases.NewGenerateCertsCmd()
	initRunner.AppendPhase(phases.NewGenerateCertsCmd())
	initRunner.AppendPhase(phases.NewGenerateIgnCmd())
	initRunner.AppendPhase(phases.NewGenerateTFCmd())
	cmd.PersistentFlags().StringVarP(&config, "config", "c", "", "config for init")
	return cmd
}

func (i *initData) MasterCfg() *nkd.Master {
	return i.mastercfg
}

func (i *initData) WorkerCfg() *nkd.Worker {
	return i.workercfg
}
func newInitData(cmd *cobra.Command, args []string, cfgPath string) (*initData, string, error) {
	// var newNkd *nkd.Master
	cfg, nodetype, err := config.LoadOrDefaultInitConfiguration(cfgPath)
	if err != nil {
		return nil, "", err
	}
	_, ok := cfg.(*nkd.Master)
	if ok == true {
		return &initData{
			mastercfg: cfg.(*nkd.Master),
		}, nodetype, nil
	} else {
		return &initData{
			workercfg: cfg.(*nkd.Worker),
		}, nodetype, nil
	}
}
