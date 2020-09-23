package cli

import (
	"encoding/json"

	"github.com/via-justa/admiral/database"
	"github.com/via-justa/admiral/datastructs"
)

func (conf *Config) GenPrometheusSDFile() (promSDFile []byte, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return nil, err
	}

	hostsWithGroups := []datastructs.Host{}

	hosts, err := db.GetHosts()
	if err != nil {
		return nil, err
	}

	for _, host := range hosts {
		if host.Monitored {
			updated, err := conf.getHostGroups(host)
			if err != nil {
				return nil, err
			}
			hostsWithGroups = append(hostsWithGroups, updated)
		}
	}

	prom := []datastructs.Prometheus{}

	for _, host := range hostsWithGroups {
		pHost := datastructs.Prometheus{}
		pHost.Targets = []string{host.Hostname + "." + host.Domain}
		pHost.Lables.Groups = host.Groups
		prom = append(prom, pHost)
	}

	promSDFile, err = json.Marshal(prom)
	if err != nil {
		return nil, err
	}

	return promSDFile, err

}
