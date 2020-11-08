// nolint: rowserrcheck,lll,golint
package db

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/via-justa/admiral/datastructs"
)

//DatabaseConfig specific configurations from file
type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	DB       string
}

// Database exposes a database connection
type Database struct {
	Conn *sqlx.DB
}

// Connect returns a Database connection
func Connect(conf DatabaseConfig) (Database, error) {
	var db Database

	dbConfig := mysql.Config{
		User:                 conf.User,
		Passwd:               conf.Password,
		Net:                  "tcp",
		Addr:                 conf.Host,
		DBName:               conf.DB,
		AllowNativePasswords: true,
	}

	var err error

	db.Conn, err = sqlx.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		return db, err
	}

	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}

	return db, err
}

// Hosts

// SelectHost return host information. The function will search for the host in the following order:
// By hostname, if hostname is empty by host and if both hostname and host are empty by id
func (db *Database) SelectHost(hostname string) (returnedHost datastructs.Host, err error) {
	if len(hostname) != 0 {
		err = db.Conn.Get(&returnedHost, "SELECT host_id, host, hostname, domain, variables, enabled, monitored, direct_group, inherited_groups FROM host_view WHERE hostname=?", hostname)
		if err == sql.ErrNoRows {
			return returnedHost, nil
		} else if err != nil {
			return returnedHost, err
		}

		return returnedHost, nil
	}

	return returnedHost, fmt.Errorf("please provide either hostname, host or id")
}

// GetHosts return all hosts in the inventory
func (db *Database) GetHosts() (hosts []datastructs.Host, err error) {
	rows, err := db.Conn.Query("SELECT host_id, host, hostname, domain, variables, enabled, monitored, direct_group, inherited_groups FROM host_view")
	if err == sql.ErrNoRows {
		return hosts, nil
	} else if err != nil {
		return hosts, err
	}

	for rows.Next() {
		host := new(datastructs.Host)
		if err = rows.Scan(&host.ID, &host.Host, &host.Hostname, &host.Domain, &host.Variables, &host.Enabled, &host.Monitored, &host.DirectGroup, &host.InheritedGroups); err != nil {
			return hosts, err
		}

		hosts = append(hosts, *host)
	}

	return hosts, nil
}

// InsertHost accept Host to insert or update and return the number of affected rows and error if exists
func (db *Database) InsertHost(host *datastructs.Host) (affected int64, err error) {
	sql := `INSERT INTO host (host, hostname, domain, variables, enabled, monitored) VALUES (?,?,?,?,?,?) 
	ON DUPLICATE KEY UPDATE host=?, hostname=?, domain=?, variables=?, enabled=?, monitored=?`

	res, err := db.Conn.Exec(sql, host.Host, host.Hostname, host.Domain, host.Variables, host.Enabled, host.Monitored,
		host.Host, host.Hostname, host.Domain, host.Variables, host.Enabled, host.Monitored)
	if err != nil {
		return 0, err
	}

	affected, err = res.RowsAffected()

	return affected, err
}

// DeleteHost accept Host to delete and return the number of affected rows and error if exists
func (db *Database) DeleteHost(host *datastructs.Host) (affected int64, err error) {
	res, err := db.Conn.Exec("DELETE FROM host WHERE id=?", host.ID)
	if err != nil {
		return 0, err
	}

	affected, err = res.RowsAffected()

	return affected, err
}

// ScanHosts get hosts where hostname or IP is like requested string
func (db *Database) ScanHosts(val string) (hosts []datastructs.Host, err error) {
	rows, err := db.Conn.Query("Select host_id, host, hostname, domain, variables, enabled, monitored, direct_group, inherited_groups FROM host_view WHERE hostname LIKE ? OR host LIKE ?;", "%"+val+"%", "%"+val+"%")
	if err == sql.ErrNoRows {
		return hosts, nil
	} else if err != nil {
		return hosts, err
	}

	for rows.Next() {
		host := new(datastructs.Host)
		if err = rows.Scan(&host.ID, &host.Host, &host.Hostname, &host.Domain, &host.Variables, &host.Enabled, &host.Monitored, &host.DirectGroup, &host.InheritedGroups); err != nil {
			return hosts, err
		}

		hosts = append(hosts, *host)
	}

	return hosts, nil
}

// Groups

// SelectGroup return group information. The function will search for the group in the following order:
// By name, if name is empty by id
func (db *Database) SelectGroup(name string) (returnedGroup datastructs.Group, err error) {
	if len(name) != 0 {
		err = db.Conn.Get(&returnedGroup, "SELECT group_id, name, variables, enabled, monitored, num_children, num_hosts, child_groups FROM `groups_view` WHERE name=?", name)
		if err == sql.ErrNoRows {
			return returnedGroup, nil
		} else if err != nil {
			return returnedGroup, err
		}

		return returnedGroup, nil
	}

	return returnedGroup, fmt.Errorf("please provide either name or id")
}

// GetGroups return all groups in the inventory
func (db *Database) GetGroups() (groups []datastructs.Group, err error) {
	rows, err := db.Conn.Query("SELECT group_id, name, variables, enabled, monitored, num_children, num_hosts, child_groups FROM `groups_view`")
	if err == sql.ErrNoRows {
		return groups, nil
	} else if err != nil {
		return groups, err
	}

	for rows.Next() {
		group := new(datastructs.Group)
		if err = rows.Scan(&group.ID, &group.Name, &group.Variables, &group.Enabled, &group.Monitored, &group.NumChildren, &group.NumHosts, &group.ChildGroups); err != nil {
			return groups, err
		}

		groups = append(groups, *group)
	}

	return groups, nil
}

// InsertGroup accept Group to insert or update and return the number of affected rows and error if exists
func (db *Database) InsertGroup(group *datastructs.Group) (affected int64, err error) {
	sql := "INSERT INTO `group` (name, variables, enabled, monitored) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE variables=?, enabled=?, monitored=?"

	res, err := db.Conn.Exec(sql, group.Name, group.Variables, group.Enabled,
		group.Monitored, group.Variables, group.Enabled, group.Monitored)
	if err != nil {
		return 0, err
	}

	affected, err = res.RowsAffected()

	return affected, err
}

// DeleteGroup accept Group to delete and return the number of affected rows and error if exists
func (db *Database) DeleteGroup(group *datastructs.Group) (affected int64, err error) {
	res, err := db.Conn.Exec("DELETE FROM `group` WHERE id=?", group.ID)
	if err != nil {
		return 0, err
	}

	affected, err = res.RowsAffected()

	return affected, err
}

// ScanGroups get group where group name in like requested string
func (db *Database) ScanGroups(val string) (groups []datastructs.Group, err error) {
	rows, err := db.Conn.Query("SELECT group_id, name, variables, enabled, monitored, num_children, num_hosts FROM `groups_view` WHERE name LIKE ?;", "%"+val+"%")
	if err == sql.ErrNoRows {
		return groups, nil
	} else if err != nil {
		return groups, err
	}

	for rows.Next() {
		group := new(datastructs.Group)
		if err = rows.Scan(&group.ID, &group.Name, &group.Variables, &group.Enabled, &group.Monitored, &group.NumChildren, &group.NumHosts); err != nil {
			return groups, err
		}

		groups = append(groups, *group)
	}

	return groups, nil
}

// ChildGroups

// SelectChildGroup accept either child or parent id and return slice of ids for parent or child groups respectively.
// If child is provided will return slice of parent ids
// If parent is provided will return slice of child ids
// will error if none is provided
func (db *Database) SelectChildGroup(child, parent string) (childGroups []datastructs.ChildGroup, err error) {
	if child != "" && parent != "" {
		var rows *sql.Rows

		rows, err = db.Conn.Query("SELECT relationship_id,parent, parent_id, child, child_id FROM childgroups_view WHERE parent=? AND child=?", parent, child)
		if err == sql.ErrNoRows {
			return childGroups, nil
		} else if err != nil {
			return childGroups, err
		}

		for rows.Next() {
			childGroup := new(datastructs.ChildGroup)
			if err = rows.Scan(&childGroup.ID, &childGroup.Parent, &childGroup.ParentID, &childGroup.Child, &childGroup.ChildID); err != nil {
				return childGroups, err
			}

			childGroups = append(childGroups, *childGroup)
		}

		return childGroups, err
	}

	return childGroups, fmt.Errorf("please provide child and parent group names")
}

// GetChildGroups return all child groups relationships in the inventory
func (db *Database) GetChildGroups() (childGroups []datastructs.ChildGroup, err error) {
	rows, err := db.Conn.Query("SELECT relationship_id,parent, parent_id, child, child_id FROM childgroups_view")
	if err == sql.ErrNoRows {
		return childGroups, nil
	} else if err != nil {
		return childGroups, err
	}

	for rows.Next() {
		childGroup := new(datastructs.ChildGroup)
		if err = rows.Scan(&childGroup.ID, &childGroup.Parent, &childGroup.ParentID, &childGroup.Child, &childGroup.ChildID); err != nil {
			return childGroups, err
		}

		childGroups = append(childGroups, *childGroup)
	}

	return childGroups, nil
}

// InsertChildGroup accept ChildGroup to insert and return the number of affected rows and error if exists
func (db *Database) InsertChildGroup(childGroup *datastructs.ChildGroup) (affected int64, err error) {
	sql := `INSERT INTO childgroups (child_id, parent_id) VALUES (?,?)`

	res, err := db.Conn.Exec(sql, childGroup.ChildID, childGroup.ParentID)
	if err != nil {
		return 0, err
	}

	affected, err = res.RowsAffected()

	return affected, err
}

// DeleteChildGroup accept ChildGroup to delete and return the number of affected rows and error if exists
func (db *Database) DeleteChildGroup(childGroup *datastructs.ChildGroup) (affected int64, err error) {
	res, err := db.Conn.Exec("DELETE FROM childgroups WHERE child_id=? and parent_id=?", childGroup.ChildID, childGroup.ParentID)
	if err != nil {
		return 0, err
	}

	affected, err = res.RowsAffected()

	return affected, err
}

// ScanChildGroups get child-group relationships where parent or child is like requested string
func (db *Database) ScanChildGroups(val string) (childGroups []datastructs.ChildGroup, err error) {
	rows, err := db.Conn.Query("SELECT relationship_id,parent, parent_id, child, child_id FROM childgroups_view WHERE parent LIKE ? OR child LIKE ?;", "%"+val+"%", "%"+val+"%")
	if err == sql.ErrNoRows {
		return childGroups, nil
	} else if err != nil {
		return childGroups, err
	}

	for rows.Next() {
		childGroup := new(datastructs.ChildGroup)
		if err = rows.Scan(&childGroup.ID, &childGroup.Parent, &childGroup.ParentID, &childGroup.Child, &childGroup.ChildID); err != nil {
			return childGroups, err
		}

		childGroups = append(childGroups, *childGroup)
	}

	return childGroups, nil
}

// HostGroups

// SelectHostGroup accept either host or group id and return slice of ids for groups or hosts respectively.
// If host is provided will return slice of groups ids
// If group is provided will return slice of hosts ids
// will error if none is provided
func (db *Database) SelectHostGroup(host string) (hostGroups []datastructs.HostGroup, err error) {
	if host != "" {
		rows, err := db.Conn.Query("SELECT relationship_id, `group`, group_id, host, host_id FROM hostgroup_view WHERE host=?", host)
		if err == sql.ErrNoRows {
			return hostGroups, nil
		} else if err != nil {
			return hostGroups, err
		}

		for rows.Next() {
			hostGroup := new(datastructs.HostGroup)
			if err = rows.Scan(&hostGroup.ID, &hostGroup.Group, &hostGroup.GroupID, &hostGroup.Host, &hostGroup.HostID); err != nil {
				return hostGroups, err
			}

			hostGroups = append(hostGroups, *hostGroup)
		}

		return hostGroups, nil
	}

	return hostGroups, fmt.Errorf("please provide either host or group id")
}

// GetHostGroups return all host groups relationships in the inventory
func (db *Database) GetHostGroups() (hostGroups []datastructs.HostGroup, err error) {
	rows, err := db.Conn.Query("SELECT relationship_id, `group`, group_id, host, host_id FROM hostgroup_view")
	if err == sql.ErrNoRows {
		return hostGroups, nil
	} else if err != nil {
		return hostGroups, err
	}

	for rows.Next() {
		hostGroup := new(datastructs.HostGroup)
		if err = rows.Scan(&hostGroup.ID, &hostGroup.Group, &hostGroup.GroupID, &hostGroup.Host, &hostGroup.HostID); err != nil {
			return hostGroups, err
		}

		hostGroups = append(hostGroups, *hostGroup)
	}

	return hostGroups, nil
}

// InsertHostGroup accept HostGroup to insert and return the number of affected rows and error if exists
func (db *Database) InsertHostGroup(hostGroup *datastructs.HostGroup) (affected int64, err error) {
	sql := `INSERT INTO hostgroups (host_id, group_id) VALUES (?,?)`

	res, err := db.Conn.Exec(sql, hostGroup.HostID, hostGroup.GroupID)
	if err != nil {
		return 0, err
	}

	affected, err = res.RowsAffected()

	return affected, err
}

// DeleteHostGroup accept HostGroup to delete and return the number of affected rows and error if exists
func (db *Database) DeleteHostGroup(hostGroup *datastructs.HostGroup) (affected int64, err error) {
	res, err := db.Conn.Exec("DELETE FROM hostgroups WHERE host_id=? and group_id=?", hostGroup.HostID, hostGroup.GroupID)
	if err != nil {
		return 0, err
	}

	affected, err = res.RowsAffected()

	return affected, err
}

// ScanHostGroups get host-groups where group is like requested string
func (db *Database) ScanHostGroups(val string) (hostGroups []datastructs.HostGroup, err error) {
	rows, err := db.Conn.Query("Select relationship_id, host_id, host, group_id, `group` FROM hostgroup_view WHERE `group` LIKE ?", "%"+val+"%")
	if err == sql.ErrNoRows {
		return hostGroups, nil
	} else if err != nil {
		return hostGroups, err
	}

	for rows.Next() {
		hg := new(datastructs.HostGroup)
		if err = rows.Scan(&hg.ID, &hg.HostID, &hg.Host, &hg.GroupID, &hg.Group); err != nil {
			return hostGroups, err
		}

		hostGroups = append(hostGroups, *hg)
	}

	return hostGroups, nil
}
