package cmd

import (
	"strings"
	"testing"

	// sqlite driver

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func Test_hostsArgsFunc(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	rootCmd := &cobra.Command{
		Use:               "root",
		ValidArgsFunction: hostsArgsFunc,
		Run:               emptyRun,
	}

	// Test that both sub-commands and validArgsFunction are completed
	output, err := executeCommand(rootCmd, cobra.ShellCompNoDescRequestCmd, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Join([]string{"host1", "host2", "host3", ":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}
}

func Test_groupsArgsFunc(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	rootCmd := &cobra.Command{
		Use:               "root",
		ValidArgsFunction: groupsArgsFunc,
		Run:               emptyRun,
	}

	// Test that both sub-commands and validArgsFunction are completed
	output, err := executeCommand(rootCmd, cobra.ShellCompNoDescRequestCmd, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Join([]string{"group1", "group2", "group3", "group4", "group5", ":0",
		"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}
}
