package main

import "github.com/spf13/cobra"

func addCommands(command *cobra.Command){
	command.AddCommand(generateCommand)
	command.AddCommand(deployCommand)
	//command.AddCommand(initCommand)
	command.AddCommand(joinCommand)
	command.AddCommand(upgradeCommand)
}


func main() {
	addCommands(rootCmd)
	_ = rootCmd.Execute()
}
