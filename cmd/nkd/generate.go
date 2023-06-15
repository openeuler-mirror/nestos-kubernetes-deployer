package main

import "github.com/spf13/cobra"

var (
	generateCommand = &cobra.Command{
		Use:               "generate [options] [REGISTRY]",
		Short:             "generate cluster bootconfig",
		Long:              "generate cluster bootconfig for a new cluster",
		RunE:              generate,
		Args:              cobra.MaximumNArgs(1),
	}
)


func generate(command *cobra.Command, args []string) error{
	println("generate")
	return nil
}

