package cli

import (
	"encoding/json"

	"github.com/via-justa/admiral/datastructs"
)

type inventoryData struct {
	hosts       []datastructs.Host
	groups      []datastructs.Group
	childGroups []datastructs.ChildGroupView
	hostGroups  []datastructs.HostGroup
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

	inv.hostGroups, err = db.getHostGroups()
	if err != nil {
		return inv, err
	}

	return inv, nil
}

func (inv *inventoryData) getChildren(parent datastructs.Group) (children []string) {
	// Get group children
	for _, childGroup := range inv.childGroups {
		if childGroup.ParentID == parent.ID {
			children = append(children, childGroup.Child)
		}
	}

	return children
}

func (inv *inventoryData) getGroupHosts(parent datastructs.Group) (groupHosts []string) {
	for _, hostGroup := range inv.hostGroups {
		if hostGroup.Group == parent.ID {
			for _, host := range inv.hosts {
				if hostGroup.Host == host.ID {
					groupHosts = append(groupHosts, host.Hostname+"."+host.Domain)
				}
			}
		}
	}

	return groupHosts
}

func (inv *inventoryData) buildInventoryGroups() (datastructs.InventoryGroups, error) {
	inventoryGroups := datastructs.InventoryGroups{}

	for _, parent := range inv.groups {
		if parent.Enabled {
			children := inv.getChildren(parent)

			// get group vars
			var GroupVars datastructs.InventoryVars

			err := json.Unmarshal([]byte(parent.Variables), &GroupVars)
			if err != nil {
				return nil, err
			}

			groupHosts := inv.getGroupHosts(parent)

			inventoryGroups[parent.Name] = datastructs.InventoryGroupsData{
				Children: children,
				Vars:     GroupVars,
				Hosts:    groupHosts,
			}
		}
	}

	return inventoryGroups, nil
}

// GenInventory return the entire inventory in Ansible acceptable json structure
func GenInventory() ([]byte, error) {
	invData, err := getInventoryData()
	if err != nil {
		return nil, err
	}

	// generate inventory hosts
	inventoryHosts := datastructs.InventoryHosts{}

	for _, host := range invData.hosts {
		if host.Enabled {
			var hostVars datastructs.InventoryVars

			err = json.Unmarshal([]byte(host.Variables), &hostVars)
			if err != nil {
				return nil, err
			}

			inventoryHosts[host.Hostname+"."+host.Domain] = hostVars
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

	invBytes, _ := json.Marshal(m)

	return invBytes, nil
}
