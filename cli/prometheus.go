package cli

import (
	"encoding/json"
	"strings"

	"github.com/via-justa/admiral/datastructs"
)

// GenPrometheusSDFile return the entire inventory in Prometheus static file acceptable json structure
func GenPrometheusSDFile() (promSDFile []byte, err error) {
	hosts, err := db.getHosts()
	if err != nil {
		return nil, err
	}

	prom := []datastructs.Prometheus{}

	for i := range hosts {
		if hosts[i].Enabled && hosts[i].Monitored {
			pHost := datastructs.Prometheus{}
			pHost.Targets = []string{hosts[i].Hostname + "." + hosts[i].Domain}
			pHost.Lables.Groups = mergeGroups(&hosts[i])
			prom = append(prom, pHost)
		}
	}

	promSDFile, err = json.Marshal(prom)
	if err != nil {
		return nil, err
	}

	return promSDFile, err
}

func mergeGroups(host *datastructs.Host) []string {
	var groups = new([]string)
	if len(host.DirectGroup) != 0 {
		*groups = append(*groups, strings.Split(host.DirectGroup, ",")...)

		if len(host.InheritedGroups) != 0 {
			for _, group := range strings.Split(host.InheritedGroups, ",") {
				for _, g := range *groups {
					if g == group {
						break
					}
				}

				*groups = append(*groups, group)
			}
		}
	}

	return *groups
}
