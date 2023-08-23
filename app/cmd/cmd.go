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
	"io"

	"github.com/spf13/cobra"
)

func NewNkdCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "nkd",
		Short: "nkd: easily bootstrap a secure Kubernetes cluster",
	}

	cmds.ResetFlags()
	// TODO: 修改名称
	// TODO: 添加ExtendCommand
	cmds.AddCommand(NewConfigCommand())
	cmds.AddCommand(NewInitCommand())
	cmds.AddCommand(NewDeployCommand())
	cmds.AddCommand(NewDestroyCommand())

	return cmds
}
