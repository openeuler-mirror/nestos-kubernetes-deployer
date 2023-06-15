package main

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:                   "nkd [options]",
		Long:                  "deploy k8s clusters",
		SilenceUsage:          true,
		SilenceErrors:         true,
		TraverseChildren:      true,
		//PersistentPreRunE:     persistentPreRunE,
		//RunE:                  validate.SubCommandExists,
		//PersistentPostRunE:    persistentPostRunE,
		//Version:               version.Version.String(),
		DisableFlagsInUseLine: true,
	}

)
