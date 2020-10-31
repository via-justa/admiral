package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/via-justa/admiral/release"
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
	if msg := release.CheckForUpdates(AppVersion); msg != "" {
		fmt.Println(release.CheckForUpdates(AppVersion))
	} else {
		fmt.Printf("Admiral version: %v\n", AppVersion)
	}
}
