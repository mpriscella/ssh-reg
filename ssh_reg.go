package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("ssh-reg", "ssh-reg is a program to manage a user's ssh config file.")

	add             = app.Command("add", "Add host entry.")
	addHost         = add.Arg("host", "The name of the host to add.").Required().String()
	addHostName     = add.Arg("hostname", "The host's URL.").Required().String()
	addIdentityFile = add.Flag("identity", "The path to the identity key file.").Default("").Short('i').String()
	addUser         = add.Flag("user", "The SSH user.").Default("").Short('u').String()
	addForce        = add.Flag("force", "Overwrite the specified host.").Short('f').Bool()
	addExtra        = add.Flag("extra", "Add Extra Keyword.").Short('e').String()

	copy        = app.Command("copy", "Copy host entry.")
	copyHost    = copy.Arg("source", "Source host.").Required().String()
	copyNewHost = copy.Arg("destination", "Destination host").Required().String()

	describe     = app.Command("describe", "Describe host entry.")
	describeHost = describe.Arg("host", "The host name.").Required().String()

	list = app.Command("list", "List all available host entries.")

	move        = app.Command("move", "Rename host entry.")
	moveHost    = move.Arg("source", "Source host.").Required().String()
	moveNewHost = move.Arg("destination", "Destination host.").Required().String()

	remove     = app.Command("remove", "Remove host entry.")
	removeHost = remove.Arg("host", "The name of the host to remove.").Required().String()

	update             = app.Command("update", "Update host entry.")
	updateHost         = update.Arg("host", "The name of the host to update.").Required().String()
	updateHostName     = update.Arg("hostname", "The host's URL.").Default("").String()
	updateIdentityFile = update.Flag("identity", "The path to the identity key file.").Short('i').String()
	updateUser         = update.Flag("user", "The SSH User.").Default("").Short('u').String()
	updateExtra        = update.Flag("extra", "Keyword=Value.").Short('e').String()
)

var sshConfig string
var entries map[string]host
var context *kingpin.ParseContext

type host struct {
	Host         string
	HostName     string
	IdentityFile string
	User         string
	Extras       map[string]string
}

func main() {
	app.Version("1.1.1")
	app.Author("Mike Priscella & Dario Castañé")

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	dir := usr.HomeDir
	sshConfig = dir + filepath.FromSlash("/.ssh/config")
	entries = make(map[string]host)
	if _, err := os.Stat(sshConfig); err == nil {
		input, err := ioutil.ReadFile(sshConfig)
		if err != nil {
			panic(err)
		}
		err = parseConfig(string(input))
		if err != nil {
			panic(err)
		}
	} else {
		if !os.IsNotExist(err) {
			panic(err)
		}
	}

	context, err := app.ParseContext(os.Args[1:])
	if err != nil {
		panic(err)
	}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case add.FullCommand():
		if _, exists := entries[*addHost]; exists {
			if *addForce {
				delete(entries, *addHost)
				err = addEntry(*addHost, *addHostName, *addIdentityFile, *addUser, *addExtra)
			} else {
				app.FatalUsageContext(context, fmt.Sprintf("Host '%v' already exists. Use --force to overwrite.", *addHost))
			}
		} else {
			err = addEntry(*addHost, *addHostName, *addIdentityFile, *addUser, *addExtra)
		}
	case copy.FullCommand():
		if _, exists := entries[*copyHost]; exists {
			entry := entries[*copyHost]
			entries[*copyNewHost] = host{Host: *copyNewHost, HostName: entry.HostName, IdentityFile: entry.IdentityFile, User: entry.User}
			saveEntries()
		} else {
			app.FatalUsageContext(context, fmt.Sprintf("Host '%v' doesn't exist.", *copyHost))
		}
	case describe.FullCommand():
		if _, exists := entries[*describeHost]; exists {
			fmt.Printf(printEntry(entries[*describeHost]))
		} else {
			app.FatalUsageContext(context, fmt.Sprintf("Host '%v' doesn't exist.", *describeHost))
		}
	case list.FullCommand():
		listEntries()
	case move.FullCommand():
		if _, exists := entries[*moveHost]; exists {
			entry := entries[*moveHost]
			entries[*moveNewHost] = host{Host: *moveNewHost, HostName: entry.HostName, IdentityFile: entry.IdentityFile, User: entry.User}
			delete(entries, *moveHost)
			saveEntries()
		} else {
			app.FatalUsageContext(context, fmt.Sprintf("Host '%v' doesn't exist.", *moveHost))
		}
	case remove.FullCommand():
		if _, exists := entries[*removeHost]; exists {
			delete(entries, *removeHost)
			saveEntries()
		} else {
			app.FatalUsageContext(context, fmt.Sprintf("Host '%v' doesn't exist.", *removeHost))
		}
	case update.FullCommand():
		if _, exists := entries[*updateHost]; exists {
			updateEntry(*updateHost, *updateHostName, *updateIdentityFile, *updateUser, *updateExtra)
		} else {
			app.FatalUsageContext(context, fmt.Sprintf("Host '%v' doesn't exist.", *updateHost))
		}
	}
	if err != nil {
		panic(err)
	}
}

func parseConfig(input string) error {
	entryRegex, err := regexp.Compile("((.+) (.+)\\s?)+")
	if err != nil {
		return err
	}
	hostOption, err := regexp.Compile("(?:\\s+)?(.+) (.+)\\s?")
	if err != nil {
		return err
	}
	matches := entryRegex.FindAllStringSubmatch(string(input), -1)

	var hosts []string
	for i := 0; i < len(matches); i++ {
		hosts = append(hosts, matches[i][0])
	}

	for _, entry := range hosts {
		options := strings.Split(entry, "\n")
		output := host{}
		output.Extras = make(map[string]string)

		for _, option := range options {
			if len(option) > 1 {
				optionMatches := hostOption.FindAllStringSubmatch(option, -1)
				key := optionMatches[0][1]
				value := optionMatches[0][2]
				switch key {
				case "Host":
					output.Host = value
					break
				case "HostName":
					output.HostName = value
					break
				case "IdentityFile":
					output.IdentityFile = value
					break
				case "User":
					output.User = value
					break
				default:
					output.Extras[key] = value
					break
				}
			}
		}
		entries[output.Host] = output
	}
	return nil
}

func validateExtras(input []string) bool {
	validKeywords := extraKeywords()
	for _, extra := range input {
		keyword := strings.Split(extra, "=")
		if !stringInSlice(keyword[0], validKeywords) {
			app.FatalUsageContext(context, fmt.Sprintf("Invalid Keyword: %s", keyword[0]))
			return false
		}
	}
	return true
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func listEntries() {
	var keys []string
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(fmt.Sprintf("%v", k))
	}
}

func printEntry(h host) string {
	hostTemplate := []string{fmt.Sprintf("Host %v\n", h.Host), fmt.Sprintf("  HostName %v\n", h.HostName)}
	if h.IdentityFile != "" {
		hostTemplate = append(hostTemplate, fmt.Sprintf("  IdentityFile %v\n", h.IdentityFile))
	}
	if h.User != "" {
		hostTemplate = append(hostTemplate, fmt.Sprintf("  User %v\n", h.User))
	}
	if len(h.Extras) > 0 {
		for k, v := range h.Extras {
			hostTemplate = append(hostTemplate, fmt.Sprintf("  %v %v\n", k, v))
		}
	}
	return strings.Join(hostTemplate, "")
}

func addEntry(hostID string, hostName string, identityFile string, user string, extra string) error {
	entry := host{Host: hostID, HostName: hostName, IdentityFile: identityFile, User: user}
	entry.Extras = make(map[string]string)

	if extra != "" {
		extraArray := strings.Split(extra, ",")
		if validateExtras(extraArray) {
			for _, v := range extraArray {
				option := strings.Split(v, "=")
				if option[1] != "" {
					entry.Extras[option[0]] = option[1]
				}
			}
		}
	}
	entries[hostID] = entry
	return saveEntries()
}

func updateEntry(host string, hostName string, identityFile string, user string, extra string) {
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
	if extra != "" {
		extraArray := strings.Split(extra, ",")
		if validateExtras(extraArray) {
			for _, v := range extraArray {
				option := strings.Split(v, "=")
				if option[1] != "" {
					entry.Extras[option[0]] = option[1]
				} else {
					delete(entry.Extras, option[0])
				}
			}
		}
	}
	entries[host] = entry
	saveEntries()
}

func saveEntries() error {
	fh, err := os.OpenFile(sshConfig, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	var keys []string
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fh.WriteString(fmt.Sprintf("%v\n", printEntry(entries[k])))
	}
	return nil
}
