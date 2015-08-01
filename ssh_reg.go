package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
	User "os/user"
	"regexp"
	"strings"
)

var (
	app = kingpin.New("ssh-reg", "A ssh config management tool.")

	list = app.Command("list", "List all available hosts")

	describe     = app.Command("describe", "Describe host")
	describeHost = describe.Arg("host", "The name of the host").Required().String()

	add             = app.Command("add", "Add host")
	addHost         = add.Arg("host", "The name of the host").Required().String()
	addHostName     = add.Arg("hostname", "The HostName of the specified host").Required().String()
	addIdentityFile = add.Flag("identity", "The location of the hosts private key").Default("").Short('i').String()
	addUser         = add.Flag("user", "The SSH User").Default("").Short('u').String()
	addForce        = add.Flag("force", "Overwrite the specified host").Short('f').Bool()

	remove     = app.Command("remove", "Remove host")
	removeHost = remove.Arg("host", "The name of the host").Required().String()

	update             = app.Command("update", "Update host")
	updateHost         = update.Arg("host", "The name of the host").Required().String()
	updateHostName     = update.Arg("hostname", "The HostName of the specified host").Default("").String()
	updateIdentityFile = update.Flag("identity", "The location of the hosts private key").Short('i').String()
	updateUser         = update.Flag("user", "The SSH User").Default("").Short('u').String()
)

var ssh_config string
var entries map[string]Host

type Host struct {
	Host         string
	HostName     string
	IdentityFile string
	User         string
}

func main() {
	kingpin.Version("0.5.0")
	usr, _ := User.Current()
	dir := usr.HomeDir
	ssh_config = dir + "/configtest"
	input, _ := ioutil.ReadFile(ssh_config)
	entries = make(map[string]Host)
	_parseConfig(string(input))

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case list.FullCommand():
		_listHosts()
		break
	case describe.FullCommand():
		_, exists := entries[*describeHost]
		if exists {
			fmt.Printf(_printHost(entries[*describeHost]))
		} else {
			fmt.Println(fmt.Sprintf("ssh-reg: Host '%v' doesn't exist.", *describeHost))
			app.Usage(os.Args[1:])
		}
		break
	case add.FullCommand():
		_, exists := entries[*addHost]
		if exists {
			if *addForce {
				delete(entries, *addHost)
				_addHost(*addHost, *addHostName, *addIdentityFile, *addUser)
			} else {
				fmt.Println(fmt.Sprintf("ssh-reg: Host '%v' already exists. Use --force to overwrite.", *addHost))
				app.Usage(os.Args[1:])
			}
		} else {
			_addHost(*addHost, *addHostName, *addIdentityFile, *addUser)
		}
		break
	case remove.FullCommand():
		_, exists := entries[*removeHost]
		if exists {
			delete(entries, *removeHost)
		} else {
			fmt.Println(fmt.Printf("ssh-reg: Host '%v' doesn't exist.", *removeHost))
			app.Usage(os.Args[1:])
		}
		break
	case update.FullCommand():
		_, exists := entries[*updateHost]
		if exists {
			_updateHost(*updateHost, *updateHostName, *updateIdentityFile, *updateUser)
		} else {
			fmt.Println(fmt.Printf("ssh-reg: Host '%v' doesn't exist.", *updateHost))
			app.Usage(os.Args[1:])
		}
		break
	}
}

func _parseConfig(input string) {
	regex, _ := regexp.Compile("Host (.+)\\s+HostName (.+)\\s+((IdentityFile|User) (.+)\\s+)?((IdentityFile|User) (.+)\\s+)?")
	matches := regex.FindAllStringSubmatch(string(input), -1)

	for _, match := range matches {
		output := Host{Host: match[1], HostName: match[2]}
		switch match[4] {
		case "IdentityFile":
			output.IdentityFile = match[5]
			break
		case "User":
			output.User = match[5]
			break
		}
		switch match[7] {
		case "IdentityFile":
			output.IdentityFile = match[8]
			break
		case "User":
			output.User = match[8]
			break
		}
		entries[match[1]] = output
	}
}

func _listHosts() {
	regex, _ := regexp.Compile(`Host (.+)`)
	input, _ := ioutil.ReadFile(ssh_config)
	match := regex.FindAllStringSubmatch(string(input), -1)

	for _, host := range match {
		fmt.Println(fmt.Sprintf("%v", host[1]))
	}
}

func _printHost(host Host) string {
	hostTemplate := []string{fmt.Sprintf("Host %v\n", host.Host), fmt.Sprintf("  HostName %v\n", host.HostName)}
	if host.IdentityFile != "" {
		hostTemplate = append(hostTemplate, fmt.Sprintf("  IdentityFile %v\n", host.IdentityFile))
	}
	if host.User != "" {
		hostTemplate = append(hostTemplate, fmt.Sprintf("  User %v\n", host.User))
	}
	return strings.Join(hostTemplate, "")
}

func _addHost(host string, hostName string, identityFile string, user string) {
	entries[host] = Host{Host: host, HostName: hostName, IdentityFile: identityFile, User: user}
	_saveEntries()
}

func _updateHost(host string, hostName string, identityFile string, user string) {
	entry := entries[host]
	if hostName != "" {
		entry.HostName = hostName
	}
	if identityFile != "" {
		entry.IdentityFile = identityFile
	}
	if user != "" {
		entry.User = user
	}
	entries[host] = entry
	_saveEntries()
}

func _saveEntries() {
	fh, _ := os.OpenFile(ssh_config, os.O_RDWR|os.O_TRUNC, 0777)

	for _, entry := range entries {
		fh.WriteString(fmt.Sprintf("%v\n", _printHost(entry)))
	}
}
