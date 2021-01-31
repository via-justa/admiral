package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/via-justa/admiral/datastructs"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sshVar)
}

var sshVar = &cobra.Command{
	Use:   "ssh {hostname}",
	Short: "ssh to inventory host",
	Long: "use host auto-complete for ssh, if `proxy` is set to true," +
		" proxy the connection / command through the server configured on `ssh-proxy`" +
		" The domain will be appended to the hostname automatically",
	Example:            "admiral ssh host\nadmiral ssh host ls -l",
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	ValidArgsFunction:  hostsArgsFunc,
	Run: func(cmd *cobra.Command, args []string) {
		if err := sshFunc(args); err != nil {
			log.Fatal(err)
		}
	},
}

func sshFunc(args []string) (err error) {
	var user string

	var host datastructs.Host

	if strings.Contains(args[0], "@") {
		login := strings.Split(args[0], "@")
		user = login[0]

		host, err = viewHostByHostname(login[1])
		if err != nil {
			return err
		}
	} else {
		user = Conf.SSH.User
		host, err = viewHostByHostname(args[0])
		if err != nil {
			return err
		}
	}

	sshArgs := []string{}

	// set auth
	if Conf.SSH.KeyPath != "" {
		sshArgs = []string{"-i", Conf.SSH.KeyPath}
	} else if Conf.SSH.Password != "" {
		sshArgs = []string{"-P", Conf.SSH.Password}
	}

	// set proxy
	if Conf.SSH.Proxy && Conf.SSHProxy.Host != fmt.Sprintf("%v.%v", host.Hostname, host.Domain) {
		sshArgs = append(sshArgs, "-J", fmt.Sprintf("%v@%v:%v", Conf.SSHProxy.User, Conf.SSHProxy.Host, Conf.SSHProxy.Port))
	}

	// StrictHostKeyChecking, (won't apply on the proxy-jump host)
	if !Conf.SSH.StrictHostKeyChecking {
		sshArgs = append(sshArgs, "-o", "LogLevel=ERROR", "-o",
			"UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no")
	}

	// format host connection user@host.domain
	sshArgs = append(sshArgs, "-p", fmt.Sprint(Conf.SSH.Port),
		fmt.Sprintf("%v@%v.%v", user, host.Hostname, host.Domain))

	// pass additional arguments ass ssh command
	if len(args) > 1 {
		sshArgs = append(sshArgs, args[1:]...)
	}

	cmd := exec.Command("ssh", sshArgs...)
	// make the ssh session interactive
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok { // nolint: errorlint
			fmt.Sprintln(exitError.ExitCode())
		}
	}

	return nil
}

func viewHostByHostname(hostname string) (host datastructs.Host, err error) {
	host, err = DB.SelectHost(hostname)
	if err != nil {
		return host, err
	} else if host.Hostname == "" {
		return host, fmt.Errorf("requested host does not exists")
	}

	return host, nil
}
