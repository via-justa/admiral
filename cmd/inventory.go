package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/via-justa/admiral/datastructs"
)

func init() {
	rootCmd.AddCommand(genInventory)
}

var genInventory = &cobra.Command{
	Use:     "inventory",
	Aliases: []string{"inv"},
	Short:   "Output Ansible compatible inventory structure",
	Example: "admiral inventory\nadmiral inventory > inventory.json",
	Run: func(cmd *cobra.Command, args []string) {
		inv, err := inventory()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", inv)
	},
}

type inventoryData struct {
	hosts       []datastructs.Host
	groups      []datastructs.Group
	childGroups []datastructs.ChildGroup
}

func getInventoryData() (inv inventoryData, err error) {
	inv.hosts, err = db.getHosts()
	if err != nil {
		return inv, err
	}

	inv.groups, err = db.getGroups()
	if err != nil {
		return inv, err
	}

	inv.childGroups, err = db.getChildGroups()
	if err != nil {
		return inv, err
	}

	return inv, nil
}

func (inv *inventoryData) getChildren(parent *datastructs.Group) (children []string) {
	// Get group children
	for _, childGroup := range inv.childGroups {
		if childGroup.ParentID == parent.ID {
			children = append(children, childGroup.Child)
		}
	}

	return children
}

func (inv *inventoryData) getGroupHosts(group *datastructs.Group) (groupHosts []string) {
	for i := range inv.hosts {
		if inv.hosts[i].DirectGroup == group.Name {
			groupHosts = append(groupHosts, inv.hosts[i].Hostname+"."+inv.hosts[i].Domain)
		}
	}

	return groupHosts
}

func (inv *inventoryData) buildInventoryGroups() (datastructs.InventoryGroups, error) {
	inventoryGroups := datastructs.InventoryGroups{}

	for _, parent := range inv.groups {
		if parent.Enabled {
			p := parent
			children := inv.getChildren(&p)

			// get group vars
			var GroupVars datastructs.InventoryVars

			err := json.Unmarshal([]byte(parent.Variables), &GroupVars)
			if err != nil {
				return nil, err
			}

			groupHosts := inv.getGroupHosts(&p)

			inventoryGroups[parent.Name] = datastructs.InventoryGroupsData{
				Children: children,
				Vars:     GroupVars,
				Hosts:    groupHosts,
			}
		}
	}

	return inventoryGroups, nil
}

// inventory return the entire inventory in Ansible acceptable json structure
func inventory() ([]byte, error) {
	invData, err := getInventoryData()
	if err != nil {
		return nil, err
	}

	// generate inventory hosts
	inventoryHosts := datastructs.InventoryHosts{}

	for i := range invData.hosts {
		if invData.hosts[i].Enabled {
			var hostVars datastructs.InventoryVars

			err = json.Unmarshal([]byte(invData.hosts[i].Variables), &hostVars)
			if err != nil {
				return nil, err
			}

			hostVars["ansible_ssh_host"] = invData.hosts[i].Host

			inventoryHosts[invData.hosts[i].Hostname+"."+invData.hosts[i].Domain] = hostVars
		}
	}

	inventoryGroups, err := invData.buildInventoryGroups()
	if err != nil {
		return nil, err
	}

	inv := datastructs.Inventory{}
	inv.Meta.HostVars = inventoryHosts

	// As we cannot skip the top level key, we do some black marshel
	// magic to combine the groups and hosts into the inventory
	var m map[string]interface{}

	hostsBytes, _ := json.Marshal(inv)
	if err := json.Unmarshal(hostsBytes, &m); err != nil {
		return nil, err
	}

	groupsBytes, _ := json.Marshal(inventoryGroups)
	if err := json.Unmarshal(groupsBytes, &m); err != nil {
		return nil, err
	}

	invBytes, _ := json.MarshalIndent(m, "", "    ")

	return invBytes, nil
}
