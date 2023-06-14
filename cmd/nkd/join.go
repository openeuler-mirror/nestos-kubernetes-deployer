package main

import (
	"github.com/spf13/cobra"
)

var (
	joinCommand = &cobra.Command{
		Use:               "join [nodes]",
		Short:             "join nodes ",
		Long:              "join nodes to cluster",
		RunE:              join,
		Args:              cobra.MaximumNArgs(1),
	}
)


func join(command *cobra.Command, args []string) error{
	println("join")
	return nil
}
