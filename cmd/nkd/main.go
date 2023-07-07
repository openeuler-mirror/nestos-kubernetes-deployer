package main

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   filepath.Base(os.Args[0]),
		Short: "Create kubernetes cluster",
		Long:  "",
	}
	return cmd
}

func main() {
	rootCmd := newRootCmd()

	for _, subCmd := range []*cobra.Command{
		/*todo:
		  newDeployCmd()
		  newGenerateCmd()
		  newInitCmd()
		  newJoinCommand()
		*/
		newUpgradeCmd(),
	} {
		rootCmd.AddCommand(subCmd)
	}
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("Error executing nkd: %v", err)
	}
}
