package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

//nolint:errcheck
func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print client version",
	Run:   versionCmdFunc,
}

func versionCmdFunc(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %v\n", AppVersion)
}
