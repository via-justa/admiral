package cmd

import (
	"log"
	"os"

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
	Run: completionCmdFunc,
}

func completionCmdFunc(cmd *cobra.Command, args []string) {
	err := rootCmd.GenBashCompletion(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
