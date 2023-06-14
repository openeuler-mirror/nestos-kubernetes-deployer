package main

import (
	"github.com/spf13/cobra"
)

var (
	deployCommand = &cobra.Command{
		Use:               "deploy [nodes]",
		Short:             "deploy masters ",
		Long:              "deploy cluster bootconfig for a new cluster",
		RunE:              deploy,
		Args:              cobra.MaximumNArgs(1),
	}
)


func deploy(command *cobra.Command, args []string) error{
	println("deploy")
	return nil
}
