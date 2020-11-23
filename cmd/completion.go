package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

//nolint:errcheck
func init() {
	rootCmd.AddCommand(completionCmd)
}

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates completion scripts (Bash for Linux and MacOS; PowerShell for windows)",
	Long: `To load completion run
. <(admiral completion)
To configure your bash shell to load completions for each session add to your bashrc
# ~/.bashrc or ~/.profile
. <(admiral completion)
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := rootCmd.GenBashCompletion(os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func hostsArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string

	hosts, _ := DB.GetHosts()

	for _, host := range hosts {
		if strings.HasPrefix(host.Hostname, toComplete) {
			completions = append(completions, host.Hostname)
		}
	}

	return completions, cobra.ShellCompDirectiveDefault
}

func groupsArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string

	groups, _ := DB.GetGroups()

	for _, group := range groups {
		if strings.HasPrefix(group.Name, toComplete) {
			completions = append(completions, group.Name)
		}
	}

	return completions, cobra.ShellCompDirectiveDefault
}
