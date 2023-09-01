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
	"nestos-kubernetes-deployer/app/cmd/phases/extend"

	"github.com/spf13/cobra"
)

func NewExtendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extend",
		Short: "Extend kubernetes cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(extend.NewExtendMasterCommand())
	cmd.AddCommand(extend.NewExtendWorkerCommand())

	return cmd
}
