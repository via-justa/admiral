package cli

import (
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
	hostGroups, err := db.selectHostGroup(host.ID, 0)
	if err != nil {
		return res, err
	}

	groups := []int{}
	groupsName := []string{}

	// loop over direct host groups to get parents
	for _, hostGroup := range hostGroups {
		// add group to list
		groups = append(groups, hostGroup.Group)

		parents, err := getParents(hostGroup.Group, []int{})
		if err != nil {
			return res, err
		}

		groups = append(groups, parents...)
	}

	for _, g := range groups {
		group, err := db.selectGroup("", g)
		if err != nil {
			return res, err
		}

		groupsName = append(groupsName, group.Name)
	}

	host.Groups = groupsName

	return *host, nil
}
