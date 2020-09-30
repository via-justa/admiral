package cli

import (
	"encoding/json"
	"fmt"

	"github.com/via-justa/admiral/datastructs"
)

// CreateHost accept host to create
func CreateHost(host *datastructs.Host) error {
	i, err := db.insertHost(host)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("no lines affected")
	}

	return nil
}

// ViewHostByHostname accept hostname of host and return the host struct
func ViewHostByHostname(hostname string) (host datastructs.Host, err error) {
	selected, err := db.selectHost(hostname, "", 0)
	if err != nil {
		return host, err
	} else if selected.Hostname == "" {
		return host, fmt.Errorf("requested host does not exists")
	}

	host, err = getHostGroups(&selected)
	if err != nil {
		return host, err
	}

	return host, nil
}

// ViewHostByIP accept IP of host and return the host struct
func ViewHostByIP(ip string) (host datastructs.Host, err error) {
	selected, err := db.selectHost("", ip, 0)
	if err != nil {
		return host, err
	} else if selected.Hostname == "" {
		return host, fmt.Errorf("requested host does not exists")
	}

	host, err = getHostGroups(&selected)
	if err != nil {
		return host, err
	}

	return host, nil
}

// ViewHostByID accept ID of host and return the host struct
func ViewHostByID(id int) (host datastructs.Host, err error) {
	selected, err := db.selectHost("", "", id)
	if err != nil {
		return host, err
	} else if selected.Hostname == "" {
		return host, fmt.Errorf("requested host does not exists")
	}

	host, err = getHostGroups(&selected)
	if err != nil {
		return host, err
	}

	return host, nil
}

// ListHosts return all existing hosts
func ListHosts() (hosts []datastructs.Host, err error) {
	selected, err := db.getHosts()
	if err != nil {
		return hosts, err
	}

	for _, host := range selected {
		this := host

		hostWithGroups, err := getHostGroups(&this)
		if err != nil {
			return hosts, err
		}

		hosts = append(hosts, hostWithGroups)
	}

	return hosts, nil
}

// EditHost accept host to edit
func EditHost(host *datastructs.Host) error {
	err := host.UnmarshalVars()
	if err != nil {
		return err
	}

	hostB, err := json.MarshalIndent(host, "", "  ")
	if err != nil {
		return err
	}

	modifiedHostB, err := Edit(hostB)
	if err != nil {
		return err
	}

	var modifiedHost datastructs.Host

	err = json.Unmarshal(modifiedHostB, &modifiedHost)
	if err != nil {
		return err
	}

	err = modifiedHost.MarshalVars()
	if err != nil {
		return err
	}

	i, err := db.insertHost(&modifiedHost)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("no lines affected")
	}

	return nil
}

// DeleteHost accept host to remove
func DeleteHost(host *datastructs.Host) (affected int64, err error) {
	affected, err = db.deleteHost(host)
	if err != nil {
		return affected, err
	} else if affected == 0 {
		return affected, fmt.Errorf("no record matched")
	}

	return affected, nil
}

func getHostGroups(host *datastructs.Host) (res datastructs.Host, err error) {
	hostGroups, err := db.selectHostGroup(host.Hostname, "")
	if err != nil {
		return res, err
	}

	groups := []string{}

	// loop over direct host groups to get parents
	for _, hostGroup := range hostGroups {
		// add group to list
		group, err := db.selectGroup("", hostGroup.GroupID)
		if err != nil {
			return res, err
		}

		groups = append(groups, group.Name)

		parents, err := getParents(group.Name, []string{})
		if err != nil {
			return res, err
		}

		groups = append(groups, parents...)
	}

	host.Groups = groups

	return *host, nil
}
