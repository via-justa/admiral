package datastructs

// Host represents inventory host
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

// Group represent inventory group
type Group struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Variables string `json:"variables" db:"variables"`
	Enabled   bool   `json:"enable" db:"enabled"`
	Monitored bool   `json:"monitor" db:"monitored"`
}

// ChildGroup represent child-group relationship
type ChildGroup struct {
	ID     int `json:"id" db:"id"`
	Child  int `json:"child_id" db:"child_id"`
	Parent int `json:"parent_id" db:"parent_id"`
}

// HostGroup represents host-group relationship
type HostGroup struct {
	ID    int `json:"id" db:"id"`
	Host  int `json:"host_id" db:"host_id"`
	Group int `json:"group_id" db:"group_id"`
}

// Inventory struct

// InventoryVars is map used to cast inventory json vars to Ansible inventory host / group vars
type InventoryVars map[string]interface{}

// InventoryHosts implement a map to attach vars to host for Ansible inventory export
type InventoryHosts map[string]InventoryVars

// InventoryGroupsData is used in Ansible inventory export to generate JSON structure for
// group-children and group-hosts relationship
type InventoryGroupsData struct {
	Children []string      `json:"children,omitempty"`
	Vars     InventoryVars `json:"vars,omitempty"`
	Hosts    []string      `json:"hosts,omitempty"`
}

// InventoryGroups represent a map of group name and it's InventoryGroupsData
type InventoryGroups map[string]InventoryGroupsData

// Inventory is the root structure for the hosts part of Ansible inventory export
type Inventory struct {
	Meta struct {
		HostVars map[string]InventoryVars `json:"hostvars"`
	} `json:"_meta"`
}

// Prometheus struct

// Prometheus is used to construct the inventory data in Prometheus static file structure
type Prometheus struct {
	Targets []string `json:"targets"`
	Lables  struct {
		Groups []string `json:"groups"`
	} `json:"labels"`
}
