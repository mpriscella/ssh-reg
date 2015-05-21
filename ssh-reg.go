package main

// add error handling

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	// Config path should be set in hidden directory (maybe symlink?)
	// .ssh-reg/config
	config_path := "/Users/mikepriscella/.ssh/config"

	parseConfig(config_path)
	// take []byte and pass it to a function that will parse the config
	// write function to parse config

	// search

	// config.Search("Host", "prod")
}

func parseConfig(path string) Config {
	var conf Config

	file, _ := os.Open(path)
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		regex := regexp.MustCompile("(Host|HostName|IdentityFile|User) (.+)")
		fmt.Println(scanner.Text())
	}

	return conf
}

func parseConfigEntry([]byte) ConfigEntry {
	var ce ConfigEntry

	return ce
}
