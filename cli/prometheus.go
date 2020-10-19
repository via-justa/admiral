package cli

import (
	"encoding/json"

	"github.com/via-justa/admiral/datastructs"
)

// GenPrometheusSDFile return the entire inventory in Prometheus static file acceptable json structure
func GenPrometheusSDFile() (promSDFile []byte, err error) {
	hosts, err := db.getHosts()
	if err != nil {
		return nil, err
	}

	groups, err := db.getGroups()
	if err != nil {
		return nil, err
	}

	prom := []datastructs.Prometheus{}

	for i := range hosts {
		if hosts[i].Enabled && hosts[i].Monitored {
			for j := range groups {
				if groups[j].Name == hosts[i].DirectGroup {
					if groups[j].Enabled && groups[j].Monitored {
						pHost := datastructs.Prometheus{}
						pHost.Targets = []string{hosts[i].Hostname + "." + hosts[i].Domain}
						pHost.Labels.Group = hosts[i].DirectGroup
						pHost.Labels.InheritedGroups = hosts[i].InheritedGroups
						prom = append(prom, pHost)
					} else {
						break
					}
				}
			}
		}
	}

	promSDFile, err = json.Marshal(prom)
	if err != nil {
		return nil, err
	}

	return promSDFile, err
}
