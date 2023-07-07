package main

import (
	"github.com/spf13/cobra"
)

func newUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade your cluster to a newer version",
		Long:  "",
		RunE:  runUpgradeCmd,
	}
	return cmd
}

func runUpgradeCmd(command *cobra.Command, args []string) error {
	//todo: Upgrade k8s version function implementation
	return nil
}
