package datastructs

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Host represents inventory host
type Host struct {
	ID              int           `json:"-" db:"host_id"`
	Host            string        `json:"ip" db:"host"`
	Hostname        string        `json:"hostname" db:"hostname"`
	Domain          string        `json:"domain" db:"domain"`
	Variables       string        `json:"-" db:"variables"`
	PrettyVariables InventoryVars `json:"variables"`
	Enabled         bool          `json:"enable" db:"enabled"`
	Monitored       bool          `json:"monitor" db:"monitored"`
	DirectGroup     string        `json:"direct_group" db:"direct_group"`
	InheritedGroups string        `json:"-" db:"inherited_groups"`
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

// ByHostname implements sort.Interface for []Host based on
// the hostname field.
type ByHostname []Host

func (h ByHostname) Len() int {
	return len(h)
}

func (h ByHostname) Less(i, j int) bool {
	return h[i].Hostname < h[j].Hostname
}

func (h ByHostname) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// ByIP implements sort.Interface for []Host based on
// the host (IP) field.
type ByIP []Host

func (h ByIP) Len() int {
	return len(h)
}

func (h ByIP) Less(i, j int) bool {
	return h[i].Host < h[j].Host
}

func (h ByIP) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// ByDomain implements sort.Interface for []Host based on
// the Domain field.
type ByDomain []Host

func (h ByDomain) Len() int {
	return len(h)
}

func (h ByDomain) Less(i, j int) bool {
	return h[i].Domain < h[j].Domain
}

func (h ByDomain) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Hosts slice of Host
type Hosts []Host

// Sort implements sort.Interface for []Host based on
// the field field.
func (h Hosts) Sort(field string) error {
	switch strings.ToLower(field) {
	case "hostname":
		sort.Sort(ByHostname(h))
	case "ip":
		sort.Sort(ByIP(h))
	case "domain":
		sort.Sort(ByIP(h))
	default:
		return fmt.Errorf("%v is not a valid sort field", field)
	}

	return nil
}

// Group represent inventory group
type Group struct {
	ID              int           `json:"-" db:"group_id"`
	Name            string        `json:"name" db:"name"`
	Variables       string        `json:"-" db:"variables"`
	PrettyVariables InventoryVars `json:"variables"`
	Enabled         bool          `json:"enable" db:"enabled"`
	Monitored       bool          `json:"monitor" db:"monitored"`
	NumChildren     int           `json:"-" db:"num_children"`
	NumHosts        int           `json:"-" db:"num_hosts"`
	ChildGroups     string        `json:"-" db:"child_groups"`
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

// ByName implements sort.Interface for []Host based on
// the Name field.
type ByName []Group

func (g ByName) Len() int {
	return len(g)
}

func (g ByName) Less(i, j int) bool {
	return g[i].Name < g[j].Name
}

func (g ByName) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

// ByNumChildren implements sort.Interface for []Host based on
// the NumChildren (IP) field.
type ByNumChildren []Group

func (g ByNumChildren) Len() int {
	return len(g)
}

func (g ByNumChildren) Less(i, j int) bool {
	return g[i].NumChildren < g[j].NumChildren
}

func (g ByNumChildren) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

// ByNumHosts implements sort.Interface for []Host based on
// the NumHosts field.
type ByNumHosts []Group

func (g ByNumHosts) Len() int {
	return len(g)
}

func (g ByNumHosts) Less(i, j int) bool {
	return g[i].NumHosts < g[j].NumHosts
}

func (g ByNumHosts) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

// Groups slice of Group
type Groups []Group

// Sort implements sort.Interface for []Host based on
// the field field.
func (g Groups) Sort(field string) error {
	switch strings.ToLower(field) {
	case "name":
		sort.Sort(ByName(g))
	case "children-count":
		sort.Sort(ByNumChildren(g))
	case "hosts-count":
		sort.Sort(ByNumHosts(g))
	default:
		return fmt.Errorf("%v is not a valid sort field", field)
	}

	return nil
}

// ChildGroup represent child-group relationship
type ChildGroup struct {
	ID       int    `json:"-" db:"relationship_id"`
	Child    string `json:"child" db:"child"`
	ChildID  int    `json:"child_id" db:"child_id"`
	Parent   string `json:"parent" db:"parent"`
	ParentID int    `json:"parent_id" db:"parent_id"`
}

// ByParent implements sort.Interface for []Host based on
// the Parent (IP) field.
type ByParent []ChildGroup

func (g ByParent) Len() int {
	return len(g)
}

func (g ByParent) Less(i, j int) bool {
	return g[i].Parent < g[j].Parent
}

func (g ByParent) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

// ByChild implements sort.Interface for []Host based on
// the NumHosts field.
type ByChild []ChildGroup

func (g ByChild) Len() int {
	return len(g)
}

func (g ByChild) Less(i, j int) bool {
	return g[i].Child < g[j].Child
}

func (g ByChild) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

// ChildGroups slice of ChildGroup
type ChildGroups []ChildGroup

// Sort implements sort.Interface for []Host based on
// the field field.
func (g ChildGroups) Sort(field string) error {
	switch strings.ToLower(field) {
	case "parent":
		sort.Sort(ByParent(g))
	case "child":
		sort.Sort(ByChild(g))
	default:
		return fmt.Errorf("%v is not a valid sort field", field)
	}

	return nil
}

// HostGroup represents host-group
type HostGroup struct {
	ID      int    `json:"-" db:"relationship_id"`
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
	Labels  struct {
		Group           string `json:"group"`
		InheritedGroups string `json:"inherited_groups"`
	} `json:"labels"`
}
