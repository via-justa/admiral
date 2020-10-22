package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/via-justa/admiral/cli"
)

func init() {
	rootCmd.AddCommand(ping)

	ping.Flags().StringVarP(&name, "hostname", "n", "", "base hostname")
}

var ping = &cobra.Command{
	Use:   "ping",
	Short: "run `ansible -m ping` on requested host",
	Run:   pingFunc,
}

func pingFunc(cmd *cobra.Command, args []string) {
	host, err := cli.ViewHostByHostname(name)
	if err != nil {
		log.Fatal(err)
	}

	pingArgs := []string{"-i", host.Hostname + "." + host.Domain + ",", host.Hostname + "." + host.Domain, "-m", "ping"}
	pingCMD := exec.Command("ansible", pingArgs...)
	pingCMD.Stdout = os.Stdout
	pingCMD.Stderr = os.Stderr

	err = pingCMD.Run()
	if err != nil {
		log.Fatal(err)
	}
}
