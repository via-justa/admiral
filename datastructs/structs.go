package datastructs

type Host struct {
	ID        int      `json:"id" db:"id"`
	Host      string   `json:"ip" db:"host"`
	Hostname  string   `json:"hostname" db:"hostname"`
	Domain    string   `json:"domain" db:"domain"`
	Variables string   `json:"variables" db:"variables"`
	Enabled   bool     `json:"enable" db:"enabled"`
	Monitored bool     `json:"monitor" db:"monitored"`
	Groups    []string `json:"groups" db:"groups"`
}

type Group struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Variables string `json:"variables" db:"variables"`
	Enabled   bool   `json:"enable" db:"enabled"`
	Monitored bool   `json:"monitor" db:"monitored"`
}

type ChildGroup struct {
	ID     int `json:"id" db:"id"`
	Child  int `json:"child_id" db:"child_id"`
	Parent int `json:"parent_id" db:"parent_id"`
}

type HostGroup struct {
	ID    int `json:"id" db:"id"`
	Host  int `json:"host_id" db:"host_id"`
	Group int `json:"group_id" db:"group_id"`
}

// Inventory struct

type InventoryVars map[string]interface{}

type InventoryHosts map[string]InventoryVars

type InventoryGroupsData struct {
	Children []string      `json:"children,omitempty"`
	Vars     InventoryVars `json:"vars,omitempty"`
	Hosts    []string      `json:"hosts,omitempty"`
}

type InventoryGroups map[string]InventoryGroupsData

type Inventory struct {
	Meta struct {
		HostVars map[string]InventoryVars `json:"hostvars"`
	} `json:"_meta"`
}

// Prometheus struct

type Prometheus struct {
	Targets []string `json:"targets"`
	Lables  struct {
		Groups []string `json:"groups"`
	} `json:"labels"`
}
