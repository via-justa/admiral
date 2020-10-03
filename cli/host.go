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
	host, err = db.selectHost(hostname, "", 0)
	if err != nil {
		return host, err
	} else if host.Hostname == "" {
		return host, fmt.Errorf("requested host does not exists")
	}

	return host, nil
}

// ViewHostByIP accept IP of host and return the host struct
func ViewHostByIP(ip string) (host datastructs.Host, err error) {
	host, err = db.selectHost("", ip, 0)
	if err != nil {
		return host, err
	} else if host.Hostname == "" {
		return host, fmt.Errorf("requested host does not exists")
	}

	return host, nil
}

// ViewHostByID accept ID of host and return the host struct
func ViewHostByID(id int) (host datastructs.Host, err error) {
	host, err = db.selectHost("", "", id)
	if err != nil {
		return host, err
	} else if host.Hostname == "" {
		return host, fmt.Errorf("requested host does not exists")
	}

	return host, nil
}

// ListHosts return all existing hosts
func ListHosts() (hosts []datastructs.Host, err error) {
	hosts, err = db.getHosts()
	if err != nil {
		return hosts, err
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
