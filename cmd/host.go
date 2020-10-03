package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/tatsushid/go-prettytable"
	"github.com/via-justa/admiral/cli"
	"github.com/via-justa/admiral/datastructs"
)

var (
	ip        string
	domain    string
	mainGroup string
)

func init() {
	rootCmd.AddCommand(host)
	host.AddCommand(createHost)
	host.AddCommand(viewHost)
	host.AddCommand(editHost)
	host.AddCommand(deleteHost)
	host.AddCommand(listHost)

	createHost.Flags().StringVar(&ip, "ip", "", "host ip")
	createHost.Flags().StringVarP(&name, "hostname", "n", "", "base hostname")
	createHost.Flags().StringVarP(&domain, "domain", "d", "", "host domain name")
	createHost.Flags().StringVarP(&variables, "variables", "v", "", "json array of variables to set on the host level")
	createHost.Flags().StringVarP(&mainGroup, "group", "g", "", "main group the host will be in")
	createHost.Flags().BoolVarP(&enable, "enable", "e", false, "should the host be enabled")
	createHost.Flags().BoolVarP(&monitor, "monitor", "m", false, "should the host be monitored")

	viewHost.Flags().StringVar(&ip, "ip", "", "host ip")
	viewHost.Flags().IntVar(&id, "id", 0, "host id")
	viewHost.Flags().StringVarP(&name, "hostname", "n", "", "base hostname")
	viewHost.Flags().BoolVar(&toJSON, "json", false, "print as json")

	editHost.Flags().StringVar(&ip, "ip", "", "host ip")
	editHost.Flags().IntVar(&id, "id", 0, "host id")
	editHost.Flags().StringVarP(&name, "hostname", "n", "", "base hostname")

	deleteHost.Flags().StringVar(&ip, "ip", "", "host ip")
	deleteHost.Flags().IntVar(&id, "id", 0, "host id")
	deleteHost.Flags().StringVarP(&name, "hostname", "n", "", "base hostname")
}

var host = &cobra.Command{
	Use:   "host",
	Short: "Managing inventory hosts",
	Args:  cobra.MinimumNArgs(1),
}

var createHost = &cobra.Command{
	Use:   "create",
	Short: "Create new host or update existing one",
	Run:   createHostFunc,
}

func createHostFunc(cmd *cobra.Command, args []string) {
	var host datastructs.Host

	if jsonPath != "" {
		hostF, err := ioutil.ReadFile(jsonPath)
		if err != nil {
			log.Fatal(err)
		}

		if len(hostF) == 0 {
			log.Fatal("File is empty or could not be found")
		}

		err = json.Unmarshal(hostF, &host)
		if err != nil {
			log.Fatal(err)
		}

		err = host.MarshalVars()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		host = datastructs.Host{
			Host:      ip,
			Hostname:  name,
			Domain:    domain,
			Variables: variables,
			Enabled:   enable,
			Monitored: monitor,
		}
	}

	if err := cli.CreateHost(&host); err != nil {
		log.Fatal(err)
	}

	createdHost, err := cli.ViewHostByHostname(host.Hostname)
	if err != nil {
		log.Fatal(err)
	}

	printHosts([]datastructs.Host{createdHost})
}

var viewHost = &cobra.Command{
	Use:   "view",
	Short: "view host details",
	Long:  "View host information by hostname (-n,--hostname) ip (--ip) or host id (--id)",
	Run:   viewHostFunc,
}

func viewHostFunc(cmd *cobra.Command, args []string) {
	var host datastructs.Host

	var err error

	switch {
	case len(name) > 0:
		host, err = cli.ViewHostByHostname(name)
		if err != nil {
			log.Fatal(err)
		}
	case len(ip) > 0:
		host, err = cli.ViewHostByIP(ip)
		if err != nil {
			log.Fatal(err)
		}
	case id != 0:
		host, err = cli.ViewHostByID(id)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Missing selector flag use --help to get available options")
	}

	if toJSON {
		err = host.UnmarshalVars()
		if err != nil {
			log.Fatal(err)
		}

		hostB, err := json.MarshalIndent(host, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", hostB)
	} else {
		printHosts([]datastructs.Host{host})
	}
}

var editHost = &cobra.Command{
	Use:   "edit",
	Short: "interactively edit host",
	Long:  "interactively edit host by hostname (-n,--hostname) ip (--ip) or host id (--id)",
	Run:   editHostFunc,
}

func editHostFunc(cmd *cobra.Command, args []string) {
	var host datastructs.Host

	var err error

	switch {
	case len(name) > 0:
		host, err = cli.ViewHostByHostname(name)
		if err != nil {
			log.Fatal(err)
		}
	case len(ip) > 0:
		host, err = cli.ViewHostByIP(ip)
		if err != nil {
			log.Fatal(err)
		}
	case id != 0:
		host, err = cli.ViewHostByID(id)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Missing selector flag use --help to get available options")
	}

	err = cli.EditHost(&host)
	if err != nil {
		log.Fatal(err)
	}
}

var deleteHost = &cobra.Command{
	Use:   "delete",
	Short: "delete host from inventory (irevertable)",
	Long:  "delete host from inventory by hostname (-n,--hostname) ip (--ip) or host id (--id)",
	Run:   deleteHostFunc,
}

func deleteHostFunc(cmd *cobra.Command, args []string) {
	var host datastructs.Host

	var err error

	switch {
	case len(name) > 0:
		host, err = cli.ViewHostByHostname(name)
		if err != nil {
			log.Fatal(err)
		}
	case len(ip) > 0:
		host, err = cli.ViewHostByIP(ip)
		if err != nil {
			log.Fatal(err)
		}
	case id != 0:
		host, err = cli.ViewHostByID(id)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Missing selector flag use --help to get available options")
	}

	affected, err := cli.DeleteHost(&host)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Lines deleted %v\n", affected)
	}
}

var listHost = &cobra.Command{
	Use:   "list",
	Short: "list all hosts in the inventory",
	Run:   listHostFunc,
}

func listHostFunc(cmd *cobra.Command, args []string) {
	hosts, err := cli.ListHosts()
	if err != nil {
		log.Fatal(err)
	}

	printHosts(hosts)
}

func printHosts(hosts []datastructs.Host) {
	tbl, err := prettytable.NewTable([]prettytable.Column{
		{Header: "ID"},
		{Header: "IP", MinWidth: 12},
		{Header: "Hostname", MinWidth: 12},
		{Header: "domain", MinWidth: 12},
		{Header: "Enabled", MinWidth: 12},
		{Header: "Monitored", MinWidth: 12},
		{Header: "Direct Groups", MinWidth: 12},
		{Header: "Inherited Groups", MinWidth: 12},
	}...)
	if err != nil {
		log.Fatal(err)
	}

	tbl.Separator = " | "

	for i := range hosts {
		err = tbl.AddRow(hosts[i].ID, hosts[i].Host, hosts[i].Hostname, hosts[i].Domain, hosts[i].Enabled,
			hosts[i].Monitored, hosts[i].DirectGroup, hosts[i].InheritedGroups)
		if err != nil {
			log.Fatal(err)
		}
	}

	// nolint: errcheck,gosec
	tbl.Print()
}
