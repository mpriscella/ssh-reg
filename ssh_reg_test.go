package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	config := "Host dev\n  HostName dev.mpriscella.com\n  IdentityFile ~/.ssh/mpriscella\n  User mpriscella\n\nHost prod\n  HostName mpriscella.com\n  IdentityFile ~/.ssh/mpriscella\n  User mpriscella\n\nHost staging\n  HostName staging.mpriscella.com\n  IdentityFile ~/.ssh/mpriscella\n  User mpriscella\n\n"
	sshConfig = "test-config"
	entries = make(map[string]host)

	testConfig, _ := os.Create(sshConfig)
	testConfig.WriteString(config)

	input, _ := ioutil.ReadFile(sshConfig)
	parseConfig(string(input))
	exitStatus := m.Run()
	os.Remove(sshConfig)
	os.Exit(exitStatus)
}

func TestParseConfig(t *testing.T) {
	// TODO this should check content as well.
	if len(entries) != 3 {
		t.Fatalf("Expected entries: %d. Actual: %d.", 3, len(entries))
	}
}

func TestAddEntry(t *testing.T) {
	addEntry("db", "db.mpriscella.com", "~/.ssh/mpriscella", "mpriscella", "Port=3336")
	entry := entries["db"]
	fail := false
	if entry.Host != "db" {
		fail = true
	} else if entry.HostName != "db.mpriscella.com" {
		fail = true
	} else if entry.IdentityFile != "~/.ssh/mpriscella" {
		fail = true
	} else if entry.User != "mpriscella" {
		fail = true
	} else if entry.Extras["Port"] != "3336" {
		fail = true
	}
	if fail {
		t.Fatalf("Entry was not added correctly.")
	}
}

func TestUpdateEntry(t *testing.T) {
	updateEntry("db", "database.mpriscella.com", "", "mike", "Port=")
	entry := entries["db"]
	fail := false
	if entry.Host != "db" {
		fail = true
	} else if entry.HostName != "database.mpriscella.com" {
		fail = true
	} else if entry.IdentityFile != "~/.ssh/mpriscella" {
		fail = true
	} else if entry.User != "mike" {
		fail = true
	} else if entry.Extras["Port"] == "3336" {
		fail = true
	}
	if fail {
		t.Fatalf("Entry was not added correctly.")
	}
}

func TestPrintEntry(t *testing.T) {
	expected := "Host dev\n  HostName dev.mpriscella.com\n  IdentityFile ~/.ssh/mpriscella\n  User mpriscella\n"
	output := printEntry(entries["dev"])
	if expected != output {
		t.Fatalf("printEntry() output does not mach expected output.")
	}
}

func TestValidateExtras(t *testing.T) {
	keyword := "Port"
	if !validateExtras([]string{keyword}) {
		t.Fatalf("Validation for keyword '%s' failed", keyword)
	}
}
