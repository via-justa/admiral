package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/via-justa/admiral/datastructs"
)

var viewAsJSON bool

func init() {
	rootCmd.AddCommand(view)

	view.AddCommand(viewHostVar)
	viewHostVar.Flags().BoolVarP(&viewAsJSON, "json", "j", false, "view in json format (present vars)")
	view.AddCommand(viewGroupVar)
	viewGroupVar.Flags().BoolVarP(&viewAsJSON, "json", "j", false, "view in json format (present vars)")
	view.AddCommand(viewChildVar)
	view.AddCommand(viewHostGroupVar)
}

var view = &cobra.Command{
	Use:        "view",
	Aliases:    []string{"list", "ls", "get"},
	ValidArgs:  []string{"host", "group", "child"},
	ArgAliases: []string{"hosts", "groups"},
	Short:      "view existing record",
}

var viewHostVar = &cobra.Command{
	Use:   "host [hostname | 'host fqdn']",
	Short: "view existing host",
	Long: "view existing host by substring of hostname or IP or view all records when no argument passed." +
		"pass the flag `-j,--json` to view the host in json structure with host variables",
	Example:           "admiral view host\nadmiral view host host1\nadmiral view host host1 -j",
	ValidArgsFunction: hostsArgsFunc,
	Run: func(cmd *cobra.Command, args []string) {
		var hosts []datastructs.Host

		var err error

		switch len(args) {
		case 0:
			hosts, err = listHosts()
			if err != nil {
				log.Fatal(err)
			}
		case 1:
			hosts, err = scanHosts(args[0])
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("received too many arguments")
		}

		if viewAsJSON {
			if len(hosts) > 0 {
				for i := range hosts {
					_ = hosts[i].UnmarshalVars()
				}

				b, _ := json.MarshalIndent(hosts, "", "    ")

				fmt.Printf("%s\n", b)
			} else {
				log.Println("No host matched")
			}
		} else {
			printHosts(hosts)
		}
	},
}

func listHosts() (hosts []datastructs.Host, err error) {
	hosts, err = DB.GetHosts()
	if err != nil {
		return hosts, err
	}

	return hosts, nil
}

func scanHosts(val string) (hosts []datastructs.Host, err error) {
	hosts, err = DB.ScanHosts(val)
	if err != nil {
		return hosts, err
	}

	return hosts, nil
}

var viewHostGroupVar = &cobra.Command{
	Use:               "host-group ['group name']",
	Short:             "view direct hosts for groups",
	Long:              "view direct hosts for group or view all records when no argument passed",
	Example:           "admiral view host-group\nadmiral view host-group group1",
	ValidArgsFunction: groupsArgsFunc,
	Run: func(cmd *cobra.Command, args []string) {
		var hgs []datastructs.HostGroup

		var err error

		switch len(args) {
		case 0:
			hgs, err = listHostGroups()
			if err != nil {
				log.Fatal(err)
			}
		case 1:
			hgs, err = scanHostGroups(args[0])
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("received too many arguments")
		}

		printHostGroups(hgs)
	},
}

func listHostGroups() (hg []datastructs.HostGroup, err error) {
	hg, err = DB.GetHostGroups()
	if err != nil {
		return hg, err
	}

	return hg, nil
}

func scanHostGroups(val string) (hostGroups []datastructs.HostGroup, err error) {
	hostGroups, err = DB.ScanHostGroups(val)
	if err != nil {
		return hostGroups, err
	}

	return hostGroups, nil
}

var viewGroupVar = &cobra.Command{
	Use:   "group ['group name']",
	Short: "view existing group",
	Long: "view existing group by substring of group name or view all records when no argument passed" +
		"pass the flag `-j,--json` to view the group in json structure with group variables",
	Example:           "admiral view group\nadmiral view group group1\nadmiral view group group1 -j",
	ValidArgsFunction: groupsArgsFunc,
	Run: func(cmd *cobra.Command, args []string) {
		var groups []datastructs.Group

		var err error

		switch len(args) {
		case 0:
			groups, err = listGroups()
			if err != nil {
				log.Fatal(err)
			}
		case 1:
			groups, err = scanGroups(args[0])
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("received too many arguments")
		}

		if viewAsJSON {
			if len(groups) > 0 {
				for i := range groups {
					_ = groups[i].UnmarshalVars()
				}
				b, _ := json.MarshalIndent(groups, "", "    ")
				fmt.Printf("%s\n", b)
			} else {
				log.Println("No groups matched")
			}
		} else {
			printGroups(groups)
		}
	},
}

func viewGroupByName(name string) (group datastructs.Group, err error) {
	group, err = DB.SelectGroup(name)
	if err != nil {
		return group, err
	} else if group.ID == 0 {
		return group, fmt.Errorf("requested group does not exists")
	}

	return group, nil
}

func listGroups() (groups []datastructs.Group, err error) {
	groups, err = DB.GetGroups()
	if err != nil {
		return groups, err
	}

	return groups, nil
}

func scanGroups(val string) (groups []datastructs.Group, err error) {
	groups, err = DB.ScanGroups(val)
	if err != nil {
		return groups, err
	}

	return groups, nil
}

var viewChildVar = &cobra.Command{
	Use:   "child ['child group' | 'parent group']",
	Short: "view existing child-group relationship",
	Long: "view existing child-group relationship by parent" +
		" or child or view all records when no argument passed",
	Example:           "admiral view child\nadmiral view child parent-group\nadmiral view child child-group",
	ValidArgsFunction: groupsArgsFunc,
	Run: func(cmd *cobra.Command, args []string) {
		var childGroups []datastructs.ChildGroup

		var err error

		switch len(args) {
		case 0:
			childGroups, err = listChildGroups()
			if err != nil {
				log.Fatal(err)
			}
		case 1:
			childGroups, err = scanChildGroups(args[0])
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("received too many arguments")
		}

		printChildGroups(childGroups)
	},
}

func listChildGroups() (childGroups []datastructs.ChildGroup, err error) {
	childGroups, err = DB.GetChildGroups()
	if err != nil {
		return childGroups, err
	}

	return childGroups, nil
}

func viewChildGroup(child, parent string) (childGroups []datastructs.ChildGroup, err error) {
	childGroups, err = DB.SelectChildGroup(child, parent)
	if err != nil {
		return childGroups, err
	} else if childGroups == nil {
		return childGroups, fmt.Errorf("no record matched request")
	}

	return childGroups, nil
}

func scanChildGroups(val string) (childGroups []datastructs.ChildGroup, err error) {
	childGroups, err = DB.ScanChildGroups(val)
	if err != nil {
		return childGroups, err
	}

	return childGroups, nil
}
