package cli

import (
	"encoding/json"

	"github.com/via-justa/admiral/datastructs"
)

func GenInventory() ([]byte, error) {
	hosts, err := db.getHosts()
	if err != nil {
		return nil, err
	}

	groups, err := db.getGroups()
	if err != nil {
		return nil, err
	}

	childGroups, err := db.getChildGroups()
	if err != nil {
		return nil, err
	}

	hostGroups, err := db.getHostGroups()
	if err != nil {
		return nil, err
	}

	// generate inventory hosts
	inventoryHosts := datastructs.InventoryHosts{}

	for _, host := range hosts {
		if host.Enabled {
			var hostVars datastructs.InventoryVars

			err := json.Unmarshal([]byte(host.Variables), &hostVars)
			if err != nil {
				return nil, err
			}

			inventoryHosts[host.Hostname+"."+host.Domain] = hostVars
		}
	}

	// generate inventory groups
	inventoryGroups := datastructs.InventoryGroups{}

	// loop over parents
	for _, parent := range groups {
		if parent.Enabled {
			children := []string{}
			groupHosts := []string{}

			// Get group children
			for _, childGroup := range childGroups {
				if childGroup.Parent == parent.ID {
					// Get child group name
					for _, child := range groups {
						if childGroup.Child == child.ID {
							children = append(children, child.Name)
						}
					}
				}
			}

			// get group vars
			var GroupVars datastructs.InventoryVars

			err := json.Unmarshal([]byte(parent.Variables), &GroupVars)
			if err != nil {
				return nil, err
			}

			// get group hosts
			for _, hostGroup := range hostGroups {
				if hostGroup.Group == parent.ID {
					for _, host := range hosts {
						if hostGroup.Host == host.ID {
							groupHosts = append(groupHosts, host.Hostname+"."+host.Domain)
						}
					}
				}
			}

			inventoryGroups[parent.Name] = datastructs.InventoryGroupsData{
				Children: children,
				Vars:     GroupVars,
				Hosts:    groupHosts,
			}
		}
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
