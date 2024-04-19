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
package main

import (
	"io"
	"nestos-kubernetes-deployer/cmd"
	"nestos-kubernetes-deployer/cmd/command"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	terminal "golang.org/x/term"
)

func main() {
	rootCmd := newRootCmd()

	for _, subCmd := range []*cobra.Command{
		cmd.NewDeployCommand(),
		cmd.NewDestroyCommand(),
		cmd.NewUpgradeCommand(),
		cmd.NewExtendCommand(),
		cmd.NewVersionCommand(),
		cmd.NewTemplateCommand(),
	} {
		rootCmd.AddCommand(subCmd)
	}

	if err := rootCmd.Execute(); err != nil {
		return
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "nkd",
		Short:            "Creates Kubernetes Clusters",
		PersistentPreRun: runRootCmd,
	}

	cmd.PersistentFlags().StringVar(&opts.Opts.RootOptDir, "dir", "/etc/nkd", "Assets directory")
	cmd.PersistentFlags().StringVar(&opts.RootOpts.LogLevel, "log-level", "info", "log level (e.g. \"debug | info | warn | error\")")
	return cmd
}

func runRootCmd(cmd *cobra.Command, args []string) {
	logrus.SetOutput(io.Discard)

	level, err := logrus.ParseLevel(opts.RootOpts.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}

	fmt := &logrus.TextFormatter{
		ForceColors:            terminal.IsTerminal(int(os.Stderr.Fd())),
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		DisableQuote:           true,
	}
	logrus.AddHook(command.NewloggerHook(os.Stderr, level, fmt))
}
