package main

import (
	"github.com/spf13/cobra"
)

var (
	upgradeCommand = &cobra.Command{
		Use:               "deploy [nodes]",
		Short:             "deploy masters ",
		Long:              "deploy cluster bootconfig for a new cluster",
		RunE:              upgrade,
		Args:              cobra.MaximumNArgs(1),
	}
)


func upgrade(command *cobra.Command, args []string) error{
	println("deploy")
	return nil
}
