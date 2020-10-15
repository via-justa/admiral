package datastructs

import (
	"encoding/json"
)

// Host represents inventory host
type Host struct {
	ID              int           `json:"-" db:"id"`
	HostID          int           `json:"id" db:"host_id"`
	Host            string        `json:"ip" db:"host"`
	Hostname        string        `json:"hostname" db:"hostname"`
	Domain          string        `json:"domain" db:"domain"`
	Variables       string        `json:"-" db:"variables"`
	PrettyVariables InventoryVars `json:"variables"`
	Enabled         bool          `json:"enable" db:"enabled"`
	Monitored       bool          `json:"monitor" db:"monitored"`
	DirectGroup     string        `json:"direct_group" db:"direct_group"`
	InheritedGroups string        `json:"inherited_groups" db:"inherited_groups"`
}

// UnmarshalVars convert string json `Host.Variables` to json value of
// `Host.PrettyVariables`
func (h *Host) UnmarshalVars() error {
	var vars InventoryVars

	err := json.Unmarshal([]byte(h.Variables), &vars)
	if err != nil {
		return err
	}

	h.PrettyVariables = vars

	return nil
}

// MarshalVars convert json value of `Host.PrettyVariables` to string in
// `Host.Variables`
func (h *Host) MarshalVars() error {
	varsB, err := json.Marshal(h.PrettyVariables)
	if err != nil {
		return err
	}

	h.Variables = string(varsB)

	return nil
}

// Group represent inventory group
type Group struct {
	ID              int           `json:"-" db:"id"`
	GroupID         int           `json:"id" db:"group_id"`
	Name            string        `json:"name" db:"name"`
	Variables       string        `json:"-" db:"variables"`
	PrettyVariables InventoryVars `json:"variables"`
	Enabled         bool          `json:"enable" db:"enabled"`
	Monitored       bool          `json:"monitor" db:"monitored"`
	NumChildren     int           `json:"-" db:"num_children"`
	NumHosts        int           `json:"-" db:"num_hosts"`
}

// UnmarshalVars convert string json `Group.Variables` to json value of
// `Group.PrettyVariables`
func (g *Group) UnmarshalVars() error {
	var vars InventoryVars

	err := json.Unmarshal([]byte(g.Variables), &vars)
	if err != nil {
		return err
	}

	g.PrettyVariables = vars

	return nil
}

// MarshalVars convert json value of `Group.PrettyVariables` to string in
// `Group.Variables`
func (g *Group) MarshalVars() error {
	varsB, err := json.Marshal(g.PrettyVariables)
	if err != nil {
		return err
	}

	g.Variables = string(varsB)

	return nil
}

// ChildGroup represent child-group relationship
type ChildGroup struct {
	ID     int `json:"id" db:"id"`
	Child  int `json:"child_id" db:"child_id"`
	Parent int `json:"parent_id" db:"parent_id"`
}

// ChildGroupView represent child-group view data
type ChildGroupView struct {
	ID       int    `json:"id" db:"relationship_id"`
	Child    string `json:"child" db:"child"`
	ChildID  int    `json:"child_id" db:"child_id"`
	Parent   string `json:"parent" db:"parent"`
	ParentID int    `json:"parent_id" db:"parent_id"`
}

// HostGroup represents host-group relationship
type HostGroup struct {
	ID    int `json:"id" db:"id"`
	Host  int `json:"host_id" db:"host_id"`
	Group int `json:"group_id" db:"group_id"`
}

// HostGroupView represents host-group view data
type HostGroupView struct {
	ID      int    `json:"id" db:"relationship_id"`
	Host    string `json:"host" db:"host"`
	HostID  int    `json:"host_id" db:"host_id"`
	Group   string `json:"group" db:"group"`
	GroupID int    `json:"group_id" db:"group_id"`
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
		Group           string `json:"group"`
		InheritedGroups string `json:"inherited_groups"`
	} `json:"labels"`
}
