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
	host         = kingpin.Arg("host", "The name of the host").String()
	hostName     = kingpin.Arg("hostname", "The HostName of the specified host").String()
	remove       = kingpin.Flag("remove", "Sest to remove the specified host").Bool()
	identityFile = kingpin.Flag("identity", "The location of the hosts private key").Default("").Short('i').String()
	user         = kingpin.Flag("user", "The SSH User").Default("").Short('u').String()
	force        = kingpin.Flag("force", "Overwrite the specified host").Bool()
	update       = kingpin.Flag("update", "Update the specified host").Bool()
	list         = kingpin.Flag("list", "List the available hosts").Short('l').Bool()
)

var ssh_config string

func main() {
	kingpin.Version("0.0.2")
	kingpin.Parse()

	usr, _ := User.Current()
	dir := usr.HomeDir
	ssh_config = dir + "/.ssh/config"
	fh, _ := os.OpenFile(ssh_config, os.O_RDWR|os.O_APPEND, 0777)

	hostExists := searchHost(*host)
	if *list {
		listHosts()
	} else if hostExists {
		if *force {
			removeHost(*host)
			addHost(*host, *hostName, *identityFile, *user)
		} else if *update {
			updateHost(*host, *hostName, *identityFile, *user)
		} else if *remove {
			removeHost(*host)
		} else {
			fmt.Println("Host exists, use --force to overwrite.")
			kingpin.Usage()
		}
	} else {
		addHost(*host, *hostName, *identityFile, *user)
	}
	defer fh.Close()
}

func searchHost(host string) bool {
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

func listHosts() {
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

func addHost(host string, hostName string, identityFile string, user string) {
	if hostName == "" {
		kingpin.Usage()
		return
	}
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

func updateHost(host string, hostName string, identityFile string, user string) {

}

func removeHost(host string) {
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
