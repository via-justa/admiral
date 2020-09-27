package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/via-justa/admiral/cli"
)

func init() {
	rootCmd.AddCommand(genInventory)
}

var genInventory = &cobra.Command{
	Use:   "inventory",
	Short: "Output Ansible compatible inventory structure",
	Run:   genInventoryFunc,
}

func genInventoryFunc(cmd *cobra.Command, args []string) {
	inv, err := cli.GenInventory()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", inv)
}
