package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/via-justa/admiral/datastructs"
)

var (
	monitor bool
	enable  bool
	accept  bool
	ip      string
	group   string
)

func init() {
	rootCmd.AddCommand(create)
	create.PersistentFlags().BoolVarP(&accept, "accept", "y", false, "auto accept changes")
	create.PersistentFlags().BoolVarP(&monitor, "monitor", "m", true, "set monitor value [true|false] (default: true)")
	create.PersistentFlags().BoolVarP(&enable, "enable", "e", true, "set enable value [true|false] (default: true)")
	createHostVar.Flags().StringVar(&ip, "ip", "", "set host ip")
	createHostVar.Flags().StringVarP(&group, "group", "g", "", "default host group")

	create.AddCommand(createHostVar)
	create.AddCommand(createGroupVar)
	create.AddCommand(createChildVar)
}

var create = &cobra.Command{
	Use:        "create [host | group]",
	Aliases:    []string{"add", "edit"},
	ValidArgs:  []string{"host", "group"},
	ArgAliases: []string{"hosts", "groups"},
	Short:      "create or modify existing record",
}

var createHostVar = &cobra.Command{
	Use:   "host {hostname | 'host fqdn'}",
	Short: "create or modify host",
	Long: "create new host or modify existing one, expecting argument host hostname/fqdn as the host to create or edit" +
		"the new or edited host would open in your favorite editor as editable json",
	Example: "admiral create host new-host\nadmiral create" +
		" host new-host.domain.com\nadmiral edit host existing-host",
	ValidArgsFunction: hostsArgsFunc,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := createHostCase(args); err != nil {
			log.Fatal(err)
		}
	},
}

func createHostCase(args []string) error {
	var hosts datastructs.Hosts

	var err error

	hosts, err = returnHosts(args[0])
	if err != nil {
		return err
	}

	enableF := create.Flag("enable")
	if enableF.Changed {
		for i := range hosts {
			hosts[i].Enabled = enable
		}
	}

	monitorF := create.Flag("monitor")
	if monitorF.Changed {
		for i := range hosts {
			hosts[i].Monitored = monitor
		}
	}

	if ip != "" {
		if len(hosts) == 1 {
			hosts[0].Host = ip
		} else {
			return fmt.Errorf("cannot set host ip, too many host matches")
		}
	}

	if group != "" {
		for i := range hosts {
			hosts[i].DirectGroup = group
			hosts[i].InheritedGroups = ""
		}
	}

	if !enableF.Changed && !monitorF.Changed && ip == "" && group == "" {
		hosts, err = editHosts(&hosts)
		if err != nil {
			return err
		}
	}

	printHosts(hosts)

	if accept || User.confirm() {
		err = confirmedHosts(&hosts)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("aborted")
	}

	return nil
}

// returnHosts return existing records or list of hosts with one new record
func returnHosts(val string) (hosts []datastructs.Host, err error) {
	fqdn := strings.SplitN(val, ".", 2)

	hosts, err = scanHosts(fqdn[0])
	if err != nil {
		return hosts, err
	}

	switch len(hosts) {
	case 0:
		host := Conf.NewDefaultHost()
		host.Hostname = fqdn[0]

		if len(fqdn) > 1 {
			host.Domain = fqdn[1]
		}

		return []datastructs.Host{host}, nil
	case 1:
		if len(fqdn) > 1 {
			if fqdn[1] == hosts[0].Domain {
				return hosts, err
			}

			host := Conf.NewDefaultHost()
			host.Hostname = fqdn[0]
			host.Domain = fqdn[1]

			return []datastructs.Host{host}, nil
		}

		return hosts, err
	default:
		return hosts, err
	}
}

func marshalHosts(hosts *datastructs.Hosts) (b []byte, err error) {
	marshaledHosts := *hosts
	for i := range marshaledHosts {
		err = marshaledHosts[i].UnmarshalVars()
		if err != nil {
			return b, err
		}
	}

	b, err = json.MarshalIndent(marshaledHosts, "", "  ")
	if err != nil {
		return b, err
	}

	return b, err
}

func unmarshalHosts(b []byte) (datastructs.Hosts, error) {
	var hosts datastructs.Hosts

	err := json.Unmarshal(b, &hosts)
	if err != nil {
		return hosts, err
	}

	for i := range hosts {
		err := hosts[i].MarshalVars()
		if err != nil {
			return hosts, err
		}
	}

	return hosts, nil
}

func editHosts(hosts *datastructs.Hosts) (returnHosts datastructs.Hosts, err error) {
	var hostsB []byte

	hostsB, err = marshalHosts(hosts)
	if err != nil {
		return returnHosts, err
	}

	modifiedHostB, err := User.Edit(hostsB)
	if err != nil {
		return returnHosts, err
	}

	return unmarshalHosts(modifiedHostB)
}

// nolint: gocognit
func confirmedHosts(hostsToCreate *datastructs.Hosts) (err error) {
	hosts := *hostsToCreate
	for i := range hosts {
		var group datastructs.Group

		err = createHost(&hosts[i])
		if err != nil && err.Error() != "no lines affected" {
			return err
		}

		if hosts[i].DirectGroup == "" {
			log.Println("created host without group. please make sure to add the host to default group")
		} else {
			group, err = viewGroupByName(hosts[i].DirectGroup)
			if err != nil {
				return err
			}

			var existingHostGroup []datastructs.HostGroup

			// if host already got host-group relationship first delete it
			existingHostGroup, err = viewHostGroupByHost(hosts[i].Hostname)
			if err != nil && err.Error() != "no record matched request" {
				return err
			} else if existingHostGroup != nil {
				_, err = deleteHostGroup(&existingHostGroup[0])
				if err != nil {
					return err
				}
			}

			var created datastructs.Hosts

			// retrieving the created host to get its ID
			created, err = scanHosts(hosts[i].Hostname)
			if err != nil {
				return err
			}

			err = createHostGroup(&created[0], &group)
			if err != nil && !strings.Contains(err.Error(), "Duplicate entry") {
				return err
			}
		}
	}

	return err
}

func createHost(host *datastructs.Host) error {
	if host.Hostname == "" || host.Host == "" {
		return fmt.Errorf("missing mandatory field ip or hostname")
	}

	i, err := DB.InsertHost(host)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("no lines affected")
	}

	return nil
}

func viewHostGroupByHost(host string) (hostGroup []datastructs.HostGroup, err error) {
	hostGroup, err = DB.SelectHostGroup(host)
	if err != nil {
		return hostGroup, err
	} else if hostGroup == nil {
		return hostGroup, fmt.Errorf("no record matched request")
	}

	return hostGroup, nil
}

func deleteHostGroup(hostGroup *datastructs.HostGroup) (affected int64, err error) {
	affected, err = DB.DeleteHostGroup(hostGroup)
	if err != nil {
		return affected, err
	} else if affected == 0 {
		return affected, fmt.Errorf("no record matched request")
	}

	return affected, nil
}

func createHostGroup(host *datastructs.Host, group *datastructs.Group) error {
	hostGroup := &datastructs.HostGroup{
		HostID:  host.ID,
		GroupID: group.ID,
	}

	i, err := DB.InsertHostGroup(hostGroup)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("no lines affected")
	}

	return nil
}

var createGroupVar = &cobra.Command{
	Use:   "group 'group name'",
	Short: "create or modify group",
	Long: "create new group or modify existing one by passing argument group name" +
		"the new or edited group would open in your favorite editor as editable json",
	Example:           "admiral create group new-group\nadmiral edit group existing-group",
	ValidArgsFunction: groupsArgsFunc,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := createGroupCase(args); err != nil {
			log.Fatal(err)
		}
	},
}

func createGroupCase(args []string) error {
	var group datastructs.Group

	var err error

	group, err = viewGroupByName(args[0])
	if err != nil {
		if err.Error() == "requested group does not exists" {
			group = Conf.NewDefaultGroup()
			group.Name = args[0]
		} else {
			return err
		}
	}

	enableF := create.Flag("enable")
	if enableF.Changed {
		group.Enabled = enable
	}

	monitorF := create.Flag("monitor")
	if monitorF.Changed {
		group.Monitored = monitor
	}

	if !enableF.Changed && !monitorF.Changed {
		group, err = editGroup(&group)
		if err != nil {
			return err
		}
	}

	printGroups([]datastructs.Group{group})

	if accept || User.confirm() {
		err := createGroup(&group)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("aborted")
	}

	return nil
}

func unmarshalGroups(group *datastructs.Group) (b []byte, err error) {
	err = group.UnmarshalVars()
	if err != nil {
		return b, err
	}

	b, err = json.MarshalIndent(group, "", "  ")
	if err != nil {
		return b, err
	}

	return b, err
}

func editGroup(group *datastructs.Group) (returnGroup datastructs.Group, err error) {
	var groupB []byte

	groupB, err = unmarshalGroups(group)
	if err != nil {
		return returnGroup, err
	}

	modifiedgroupB, err := User.Edit(groupB)
	if err != nil {
		return returnGroup, err
	}

	err = json.Unmarshal(modifiedgroupB, &returnGroup)
	if err != nil {
		return returnGroup, err
	}

	err = returnGroup.MarshalVars()
	if err != nil {
		return returnGroup, err
	}

	return returnGroup, err
}

func createGroup(group *datastructs.Group) error {
	if group.Name == "" {
		return fmt.Errorf("missing mandatory field name")
	}

	i, err := DB.InsertGroup(group)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("no lines affected")
	}

	return nil
}

var createChildVar = &cobra.Command{
	Use:   "child 'child group' 'parent group'",
	Short: "create or modify existing child-group relationship",
	Long: "create or modify existing child-group relationship expecting ordered arguments child and parent group names." +
		" If the created relationship creates relationship loop an error will be returned",
	Example:           "admiral create child child-group parent-group",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: groupsArgsFunc,
	Run: func(cmd *cobra.Command, args []string) {
		if err := createChildCase(args); err != nil {
			log.Fatal(err)
		}
	},
}

func createChildCase(args []string) error {
	var childGroups []datastructs.ChildGroup

	var err error

	// check if relationship already exists
	childGroups, _ = viewChildGroup(args[0], args[1])
	if len(childGroups) != 0 {
		return fmt.Errorf("Group relationship already exists")
	}

	child, err := viewGroupByName(args[0])
	if err != nil {
		return err
	}

	parent, err := viewGroupByName(args[1])
	if err != nil {
		return err
	}

	childGroups = []datastructs.ChildGroup{
		{
			Parent:   parent.Name,
			ParentID: parent.ID,
			Child:    child.Name,
			ChildID:  child.ID,
		},
	}

	printChildGroups(childGroups)

	if User.confirm() {
		err = createChildGroup(&parent, &child)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("aborted")
	}

	return nil
}

func createChildGroup(parent *datastructs.Group, child *datastructs.Group) error {
	if child.ID == parent.ID {
		return fmt.Errorf("child and parent cannot be the same group")
	}

	isLoop := isRelationshipLoop(parent, child)
	if isLoop {
		return fmt.Errorf("relationship loop detected")
	}

	childGroup := &datastructs.ChildGroup{
		ParentID: parent.ID,
		ChildID:  child.ID,
	}

	i, err := DB.InsertChildGroup(childGroup)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("no lines affected")
	}

	return nil
}

func isRelationshipLoop(parent, child *datastructs.Group) bool {
	children := strings.Split(child.ChildGroups, ",")

	for _, c := range children {
		if parent.Name == c {
			return true
		}
	}

	return false
}
