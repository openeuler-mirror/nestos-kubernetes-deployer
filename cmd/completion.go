package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewCompletionCommand() *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "generate an auto-completion script",
		Long: `usage:
Bash:
  $ source <(nkd completion bash)
  
# Permanent effect (Linux):
  $ nkd completion bash > /etc/bash_completion.d/nkd

Zsh:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  $ nkd completion zsh > "${fpath[1]}/_nkd"

Fish:
  $ nkd completion fish | source
  $ nkd completion fish > ~/.config/fish/completions/nkd.fish
`,
		Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: runCompletionCmd,
	}

	return completionCmd
}

func runCompletionCmd(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "bash":
		return cmd.Root().GenBashCompletionV2(os.Stdout, true)
	case "zsh":
		return cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		return cmd.Root().GenFishCompletion(os.Stdout, true)
	default:
		return fmt.Errorf("unsupported Shell types: %s", args[0])
	}
}
