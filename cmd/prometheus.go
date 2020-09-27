package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/via-justa/admiral/cli"
)

func init() {
	rootCmd.AddCommand(genPromSDFile)
}

var genPromSDFile = &cobra.Command{
	Use:   "prometheus",
	Short: "Output prometheus compatible SD file structure",
	Run:   genPromSDFileFunc,
}

func genPromSDFileFunc(cmd *cobra.Command, args []string) {
	prom, err := cli.GenPrometheusSDFile()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", prom)
}
