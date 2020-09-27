package cli

import (
	"encoding/json"

	"github.com/via-justa/admiral/datastructs"
)

// GenPrometheusSDFile return the entire inventory in Prometheus static file acceptable json structure
func GenPrometheusSDFile() (promSDFile []byte, err error) {
	hostsWithGroups := []datastructs.Host{}

	hosts, err := db.getHosts()
	if err != nil {
		return nil, err
	}

	for _, host := range hosts {
		if host.Monitored {
			var updated datastructs.Host

			this := host

			updated, err = getHostGroups(&this)
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
