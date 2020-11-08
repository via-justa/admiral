package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ping)
}

var ping = &cobra.Command{
	Use:     "ping",
	Short:   "run `ansible -m ping` on requested host",
	Example: "admiral ping host1",
	Run:     pingFunc,
}

func pingFunc(cmd *cobra.Command, args []string) {
	host, err := viewHostByHostname(args[0])
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
