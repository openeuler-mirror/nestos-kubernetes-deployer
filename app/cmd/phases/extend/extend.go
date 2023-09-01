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

package extend

import (
	"nestos-kubernetes-deployer/app/phases/infra"

	"github.com/spf13/cobra"
)

func NewExtendMasterCommand() *cobra.Command {
	var num string
	cmd := &cobra.Command{
		Use:   "master",
		Short: "Extend kubernetes master node",
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster := &infra.Cluster{
				Dir:  "./",
				Node: "master",
				Num:  num,
			}

			if num == "" {
				return cmd.Help()
			}
			return cluster.Extend()
		},
	}
	cmd.PersistentFlags().StringVarP(&num, "num", "n", "", "number of the extend nodes")

	return cmd
}

func NewExtendWorkerCommand() *cobra.Command {
	var num string
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Extend kubernetes worker node",
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster := &infra.Cluster{
				Dir:  "./",
				Node: "worker",
				Num:  num,
			}

			if num == "" {
				return cmd.Help()
			}
			return cluster.Extend()
		},
	}
	cmd.PersistentFlags().StringVarP(&num, "num", "n", "", "number of the extend nodes")

	return cmd
}
