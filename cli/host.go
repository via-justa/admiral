package cli

import (
	"fmt"

	"github.com/via-justa/admiral/database"
	"github.com/via-justa/admiral/datastructs"
)

func (conf *Config) CreateHost(host datastructs.Host) error {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return err
	}

	i, err := db.InsertHost(host)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("No lines affected")
	}

	return nil
}

func (conf *Config) ViewHostByHostname(hostname string) (host datastructs.Host, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return host, err
	}

	selected, err := db.SelectHost(hostname, "", 0)
	if err != nil {
		return host, err
	}

	host, err = conf.getHostGroups(selected)
	if err != nil {
		return host, err
	}

	return host, nil
}

func (conf *Config) ViewHostByIP(ip string) (host datastructs.Host, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return host, err
	}

	selected, err := db.SelectHost("", ip, 0)
	if err != nil {
		return host, err
	}

	host, err = conf.getHostGroups(selected)
	if err != nil {
		return host, err
	}

	return host, nil
}

func (conf *Config) ViewHostByID(id int) (host datastructs.Host, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return host, err
	}

	selected, err := db.SelectHost("", "", id)
	if err != nil {
		return host, err
	}

	host, err = conf.getHostGroups(selected)
	if err != nil {
		return host, err
	}

	return host, nil
}

func (conf *Config) ListHosts() (hosts []datastructs.Host, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return hosts, err
	}

	selected, err := db.GetHosts()
	if err != nil {
		return hosts, err
	}

	for _, host := range selected {
		host, err = conf.getHostGroups(host)
		if err != nil {
			return hosts, err
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func (conf *Config) DeleteHost(host datastructs.Host) (affected int64, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return affected, err
	}

	affected, err = db.DeleteHost(host)
	if err != nil {
		return affected, err
	}

	return affected, nil
}

func (conf *Config) getHostGroups(host datastructs.Host) (res datastructs.Host, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return res, err
	}

	hostGroups, err := db.SelectHostGroup(host.ID, 0)
	if err != nil {
		return res, err
	}

	groups := []int{}
	groupsName := []string{}

	// loop over direct host groups to get parents
	for _, hostGroup := range hostGroups {
		// add group to list
		groups = append(groups, hostGroup.Group)

		parents, err := conf.getParents(hostGroup.Group, []int{})
		if err != nil {
			return res, err
		}
		groups = append(groups, parents...)

	}

	for _, g := range groups {
		group, err := db.SelectGroup("", g)
		if err != nil {
			return res, err
		}
		groupsName = append(groupsName, group.Name)
	}

	host.Groups = groupsName

	return host, nil
}

func (conf *Config) getParents(child int, parents []int) ([]int, error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return nil, err
	}

	childGroups, err := db.SelectChildGroup(child, 0)
	if err != nil {
		return nil, err
	}
	for _, group := range childGroups {
		parents = append(parents, group.Parent)
		child = group.Parent
		parents, err = conf.getParents(child, parents)
		if err != nil {
			return nil, err
		}
	}
	return parents, nil
}
