package cli

import (
	"fmt"

	"github.com/via-justa/admiral/datastructs"
)

func CreateHost(host datastructs.Host) error {
	i, err := db.insertHost(host)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("No lines affected")
	}

	return nil
}

func ViewHostByHostname(hostname string) (host datastructs.Host, err error) {
	selected, err := db.selectHost(hostname, "", 0)
	if err != nil {
		return host, err
	}

	host, err = getHostGroups(selected)
	if err != nil {
		return host, err
	}

	return host, nil
}

func ViewHostByIP(ip string) (host datastructs.Host, err error) {
	selected, err := db.selectHost("", ip, 0)
	if err != nil {
		return host, err
	}

	host, err = getHostGroups(selected)
	if err != nil {
		return host, err
	}

	return host, nil
}

func ViewHostByID(id int) (host datastructs.Host, err error) {
	selected, err := db.selectHost("", "", id)
	if err != nil {
		return host, err
	}

	host, err = getHostGroups(selected)
	if err != nil {
		return host, err
	}

	return host, nil
}

func ListHosts() (hosts []datastructs.Host, err error) {
	selected, err := db.getHosts()
	if err != nil {
		return hosts, err
	}

	for _, host := range selected {
		host, err = getHostGroups(host)
		if err != nil {
			return hosts, err
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func DeleteHost(host datastructs.Host) (affected int64, err error) {
	affected, err = db.deleteHost(host)
	if err != nil {
		return affected, err
	}

	return affected, nil
}

func getHostGroups(host datastructs.Host) (res datastructs.Host, err error) {
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

	return host, nil
}

func getParents(child int, parents []int) ([]int, error) {
	childGroups, err := db.selectChildGroup(child, 0)
	if err != nil {
		return nil, err
	}
	for _, group := range childGroups {
		parents = append(parents, group.Parent)
		child = group.Parent
		parents, err = getParents(child, parents)
		if err != nil {
			return nil, err
		}
	}
	return parents, nil
}
