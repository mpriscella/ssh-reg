package main

import (
	"bufio"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	User "os/user"
	"regexp"
)

var (
	host         = kingpin.Arg("host", "The name of the host").Required().String()
	hostName     = kingpin.Arg("hostname", "The HostName of the specified host").String()
	remove       = kingpin.Flag("remove", "Sest to remove the specified host").Bool()
	identityFile = kingpin.Flag("identity", "The location of the hosts private key").Default("").Short('i').String()
	user         = kingpin.Flag("user", "The SSH User").Default("").Short('u').String()
	force        = kingpin.Flag("force", "Overwrite the specified host").Bool()
	update       = kingpin.Flag("update", "Update the specified host").Bool()
)

var ssh_config string
var fh *os.File

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	usr, _ := User.Current()
	dir := usr.HomeDir
	ssh_config = dir + "/configtest"
	fh, _ = os.OpenFile(ssh_config, os.O_RDWR, 0777)

	hostExists := searchHost(*host)
	// logic here should be cleaned up. Probably check 'remove' first, then update, then force
	if hostExists {
		if *force {
			// Remove host and re-add it.
			removeHost(*host)
			addHost(*host, *hostName, *identityFile, *user)
		} else if *update {
			// Update host with whatever new values, keeping untouched values.
			updateHost(*host, *hostName, *identityFile, *user)
		} else if *remove {
			// Remove host entirely
			removeHost(*host)
		} else {
			fmt.Println("Host exists, use --force to overwrite.")
		}
	} else {
		// Add host
		addHost(*host, *hostName, *identityFile, *user)
	}
	defer fh.Close()
}

func searchHost(host string) bool {
	regex, err := regexp.Compile(fmt.Sprintf("^Host %v$", host))
	if err != nil {
		return false
	}

	f := bufio.NewReader(fh)

	buf := make([]byte, 1024)
	for {
		buf, _, err = f.ReadLine()
		if err != nil {
			return false
		}

		s := string(buf)
		if regex.MatchString(s) {
			return true
		}
	}
}

func addHost(host string, hostName string, identityFile string, user string) {
	fh.WriteString(fmt.Sprintf("Host %v\n", host))
	fh.WriteString(fmt.Sprintf("  HostName %v\n", hostName))
	if identityFile != "" {
		fh.WriteString(fmt.Sprintf("  IdentityFile %v\n", identityFile))
	}
	if user != "" {
		fh.WriteString(fmt.Sprintf("  User %v\n", user))
	}
	fh.WriteString("\n")
}

func updateHost(host string, hostName string, identityFile string, user string) {

}

func removeHost(host string) {

}
