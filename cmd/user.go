package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/tatsushid/go-prettytable"
	"github.com/via-justa/admiral/datastructs"
)

const separator = " |"

func printHosts(hosts []datastructs.Host) {
	tbl, err := prettytable.NewTable([]prettytable.Column{
		{Header: "IP", MinWidth: 12},
		{Header: "Hostname", MinWidth: 12},
		{Header: "domain", MinWidth: 12},
		{Header: "Enabled", MinWidth: 12},
		{Header: "Monitored", MinWidth: 12},
		{Header: "Direct Groups", MinWidth: 12},
		{Header: "Inherited Groups", MinWidth: 12},
	}...)
	if err != nil {
		log.Fatal(err)
	}

	tbl.Separator = separator

	for i := range hosts {
		err = tbl.AddRow(hosts[i].Host, hosts[i].Hostname, hosts[i].Domain, hosts[i].Enabled,
			hosts[i].Monitored, hosts[i].DirectGroup, hosts[i].InheritedGroups)
		if err != nil {
			log.Fatal(err)
		}
	}

	// nolint: errcheck,gosec
	tbl.Print()
}

func printGroups(groups []datastructs.Group) {
	tbl, err := prettytable.NewTable([]prettytable.Column{
		{Header: "name", MinWidth: 12},
		{Header: "Enabled", MinWidth: 12},
		{Header: "Monitored", MinWidth: 12},
		{Header: "Children count", MinWidth: 12},
		{Header: "Hosts count", MinWidth: 12},
	}...)
	if err != nil {
		log.Fatal(err)
	}

	tbl.Separator = separator

	for _, group := range groups {
		err = tbl.AddRow(group.Name, group.Enabled, group.Monitored, group.NumChildren, group.NumHosts)
		if err != nil {
			log.Fatal(err)
		}
	}

	// nolint: errcheck,gosec
	tbl.Print()
}

func printChildGroups(childGroups []datastructs.ChildGroup) {
	tbl, err := prettytable.NewTable([]prettytable.Column{
		{Header: "Parent", MinWidth: 12},
		{Header: "Parent ID", MinWidth: 12},
		{Header: "Child", MinWidth: 12},
		{Header: "Child ID", MinWidth: 12},
	}...)
	if err != nil {
		log.Fatal(err)
	}

	// nolint: goconst
	tbl.Separator = separator

	for _, childGroup := range childGroups {
		err = tbl.AddRow(childGroup.Parent, childGroup.ParentID, childGroup.Child, childGroup.ChildID)
		if err != nil {
			log.Fatal(err)
		}
	}

	// nolint: errcheck,gosec
	tbl.Print()
}

func printHostGroups(hostGroups []datastructs.HostGroup) {
	tbl, err := prettytable.NewTable([]prettytable.Column{
		{Header: "Group", MinWidth: 12},
		{Header: "Group ID", MinWidth: 12},
		{Header: "Hostname", MinWidth: 12},
		{Header: "Host ID", MinWidth: 12},
	}...)
	if err != nil {
		log.Fatal(err)
	}

	tbl.Separator = " | "

	for _, hostGroup := range hostGroups {
		err = tbl.AddRow(hostGroup.Group, hostGroup.GroupID, hostGroup.Host, hostGroup.HostID)
		if err != nil {
			log.Fatal(err)
		}
	}

	// nolint: errcheck,gosec
	tbl.Print()
}

const defaultEditor = "vim"

func getPreferredEditorFromEnvironment() string {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		return defaultEditor
	}

	return editor
}

// resolveEditorArguments for now only with VS-Code, others can be added if needed
func resolveEditorArguments(executable string, filename string) []string {
	args := []string{filename}

	if strings.Contains(executable, "code") {
		args = append([]string{"--wait"}, args...)
	}

	return args
}

func openFileInEditor(filename string, editor string) error {
	// Get the full executable path for the editor.
	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, resolveEditorArguments(executable, filename)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Edit opens a temporary file in a text editor, write data into it for editing and
// returns the written bytes on success or an error on failure. It handles deletion
// of the temporary file behind the scenes.
// nolint: gosec
func (u user) Edit(data []byte) ([]byte, error) {
	resolveEditor := getPreferredEditorFromEnvironment()

	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		return []byte{}, err
	}

	filename := file.Name()

	_, err = file.Write(data)
	if err != nil {
		return []byte{}, err
	}

	// Defer removal of the temporary file in case any of the next steps fail.
	// nolint: errcheck
	defer os.Remove(filename)

	if err = file.Close(); err != nil {
		return []byte{}, err
	}

	if err = openFileInEditor(filename, resolveEditor); err != nil {
		return []byte{}, err
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

type userInt interface {
	confirm() bool
	Edit(data []byte) ([]byte, error)
}

type user struct{}

func (u user) confirm() bool {
	r := bufio.NewReader(os.Stdin)

	for tries := 2; tries > 0; tries-- {
		fmt.Print("Please confirm [y/n]: ")

		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Empty input (i.e. "\n")
		if len(res) < 2 {
			continue
		}

		return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
	}

	return false
}

func newUser() userInt {
	u := user{}
	return u
}
