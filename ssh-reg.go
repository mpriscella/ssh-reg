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

	add             = app.Command("add", "Add host")
	addHost         = add.Arg("host", "The name of the host").Required().String()
	addHostName     = add.Arg("hostname", "The HostName of the specified host").Required().String()
	addIdentityFile = add.Flag("identity", "The location of the hosts private key").Default("").Short('i').String()
	addUser         = add.Flag("user", "The SSH User").Default("").Short('u').String()
	addForce        = add.Flag("force", "Overwrite the specified host").Bool()

	remove     = app.Command("remove", "Remove host")
	removeHost = remove.Arg("host", "The name of the host").String()

	update             = app.Command("update", "Update host")
	updateHost         = update.Arg("host", "The name of the host").String()
	updateHostName     = update.Arg("hostname", "The HostName of the specified host").String()
	updateIdentityFile = update.Flag("identity", "The location of the hosts private key").Default("").Short('i').String()
	updateUser         = update.Flag("user", "The SSH User").Default("").Short('u').String()
)

var ssh_config string

func main() {
	kingpin.Version("0.0.3")
	usr, _ := User.Current()
	dir := usr.HomeDir
	ssh_config = dir + "/.ssh/config"
	fh, _ := os.OpenFile(ssh_config, os.O_RDWR|os.O_APPEND, 0777)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case list.FullCommand():
		_listHosts()
		break
	case add.FullCommand():
		hostExists := _searchHost(*addHost)
		if hostExists {
			if *addForce {
				_removeHost(*addHost)
				_addHost(*addHost, *addHostName, *addIdentityFile, *addUser)
			} else {
				fmt.Println(fmt.Sprintf("Host '%v' already exists. Use --force to overwrite.", *addHost))
				app.Usage(os.Args[1:])
			}
		} else {
			_addHost(*addHost, *addHostName, *addIdentityFile, *addUser)
		}
		break
	case remove.FullCommand():
		hostExists := _searchHost(*removeHost)
		if hostExists {
			_removeHost(*removeHost)
		} else {
			fmt.Println(fmt.Printf("Host '%v' doesn't exist.", *removeHost))
			app.Usage(os.Args[1:])
		}
		break
	case update.FullCommand():
		hostExists := _searchHost(*updateHost)
		if hostExists {
			_updateHost(*updateHost, *updateHostName, *updateIdentityFile, *updateUser)
		} else {
			fmt.Println(fmt.Printf("Host '%v' doesn't exist.", *updateHost))
			app.Usage(os.Args[1:])
		}
		break
	}

	defer fh.Close()
}

func _searchHost(host string) bool {
	regex, _ := regexp.Compile(fmt.Sprintf("^Host %v$", host))

	input, _ := ioutil.ReadFile(ssh_config)
	lines := strings.Split(string(input), "\n")

	for _, line := range lines {
		if regex.MatchString(line) {
			return true
		}
	}
	return false
}

func _listHosts() {
	regex, _ := regexp.Compile(`^Host (.+)$`)

	input, _ := ioutil.ReadFile(ssh_config)
	lines := strings.Split(string(input), "\n")

	for _, line := range lines {
		if regex.MatchString(line) {
			match := regex.FindStringSubmatch(line)
			fmt.Println(fmt.Sprintf("%v", match[1]))
		}
	}
}

func _addHost(host string, hostName string, identityFile string, user string) {
	fh, _ := os.OpenFile(ssh_config, os.O_RDWR|os.O_APPEND, 0777)
	fh.WriteString(fmt.Sprintf("Host %v\n", host))
	fh.WriteString(fmt.Sprintf("  HostName %v\n", hostName))
	if identityFile != "" {
		fh.WriteString(fmt.Sprintf("  IdentityFile %v\n", identityFile))
	}
	if user != "" {
		fh.WriteString(fmt.Sprintf("  User %v\n", user))
	}
	fh.WriteString("\n")
	defer fh.Close()
}

func _updateHost(host string, hostName string, identityFile string, user string) {
	fmt.Printf("Not implemented yet.")
}

func _removeHost(host string) {
	hostRegex, _ := regexp.Compile(fmt.Sprintf("^Host %v$", host))
	regex, _ := regexp.Compile("^Host .+$")
	input, _ := ioutil.ReadFile(ssh_config)
	fh, _ := os.OpenFile(ssh_config, os.O_RDWR|os.O_TRUNC, 0777)
	lines := strings.Split(string(input), "\n")

	for i := 0; i < len(lines); i++ {
		if hostRegex.MatchString(lines[i]) {
			for k := i + 1; k < len(lines); k++ {
				i++
				if regex.MatchString(lines[k]) {
					fh.WriteString(fmt.Sprintf("%v\n", lines[k]))
					break
				}
			}
		} else {
			fh.WriteString(fmt.Sprintf("%v\n", lines[i]))
		}
	}

	defer fh.Close()
}
