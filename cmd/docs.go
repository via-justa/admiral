package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func init() {
	rootCmd.AddCommand(markdown)
}

var markdown = &cobra.Command{
	Use:   "docs",
	Short: "Generate commands markdown documentation",
	Long: `Generate commands markdown documentation. 
The output of the command will be written as a file per command in "./docs"`,
	Run:    markdownFunc,
	Hidden: true,
}

func markdownFunc(cmd *cobra.Command, args []string) {
	path := "./docs"

	// Make sure folder exist
	if err := os.MkdirAll(path, 0774); err != nil {
		log.Fatal(err)
	}

	err := doc.GenMarkdownTree(rootCmd, path)
	if err != nil {
		log.Fatal(err)
	}
}
