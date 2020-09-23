package cli

import (
	"fmt"

	"github.com/via-justa/admiral/database"
	"github.com/via-justa/admiral/datastructs"
)

func (conf *Config) CreateHostGroup(host datastructs.Host, group datastructs.Group) error {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return err
	}

	hostGroup := datastructs.HostGroup{
		Host:  host.ID,
		Group: group.ID,
	}

	i, err := db.InsertHostGroup(hostGroup)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("No lines affected")
	}

	return nil
}

func (conf *Config) ViewHostGroupByHost(hostID int) (hostGroup []datastructs.HostGroup, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return hostGroup, err
	}

	hostGroup, err = db.SelectHostGroup(hostID, 0)
	if err != nil {
		return hostGroup, err
	}

	return hostGroup, nil
}

func (conf *Config) ViewHostGroupByGroup(groupID int) (hostGroup []datastructs.HostGroup, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return hostGroup, err
	}

	hostGroup, err = db.SelectHostGroup(0, groupID)
	if err != nil {
		return hostGroup, err
	}

	return hostGroup, nil
}

func (conf *Config) ListHostGroup() (hostGroups []datastructs.HostGroup, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return hostGroups, err
	}

	hostGroups, err = db.GetHostGroups()
	if err != nil {
		return hostGroups, err
	}

	return hostGroups, nil
}

func (conf *Config) DeleteHostGroup(hostGroup datastructs.HostGroup) (affected int64, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return affected, err
	}

	affected, err = db.DeleteHostGroup(hostGroup)
	if err != nil {
		return affected, err
	}

	return affected, nil
}
