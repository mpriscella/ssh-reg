package main

import (
	"fmt"
)

type Config struct {
	Path    string
	Entries []ConfigEntry
}

type ConfigEntry struct {
	Host         string
	HostName     string
	IdentityFile string
	User         string
}

// This function searches the initialized config file
// for the specified option with the provided value.
func (*Config) Search(option string, value string) {

	fmt.Println("hello")
}
