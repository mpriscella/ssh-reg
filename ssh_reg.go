package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
	User "os/user"
	"regexp"
	"sort"
	"strings"
)

var (
	app = kingpin.New("ssh-reg", "\033[1mssh-reg\033[0m is a program to manage a user's ssh config file.")

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

var ssh_config string
var entries map[string]Host
var context kingpin.ParseContext

type Host struct {
	Host         string
	HostName     string
	IdentityFile string
	User         string
	Extras       map[string]string
}

func main() {
	app.Version("1.0.0")
	app.Author("Mike Priscella")

	usr, _ := User.Current()
	dir := usr.HomeDir
	ssh_config = dir + "/.ssh/config"
	input, _ := ioutil.ReadFile(ssh_config)
	entries = make(map[string]Host)
	parseConfig(string(input))
	context, _ := app.ParseContext(os.Args[1:])

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case add.FullCommand():
		_, exists := entries[*addHost]
		if exists {
			if *addForce {
				delete(entries, *addHost)
				addEntry(*addHost, *addHostName, *addIdentityFile, *addUser, *addExtra)
			} else {
				app.FatalUsageContext(context, fmt.Sprintf("ssh-reg: Host '%v' already exists. Use --force to overwrite.", *addHost))
			}
		} else {
			addEntry(*addHost, *addHostName, *addIdentityFile, *addUser, *addExtra)
		}
		break
	case copy.FullCommand():
		_, exists := entries[*copyHost]
		if exists {
			entry := entries[*copyHost]
			entries[*copyNewHost] = Host{Host: *copyNewHost, HostName: entry.HostName, IdentityFile: entry.IdentityFile, User: entry.User}
			saveEntries()
		} else {
			app.FatalUsageContext(context, fmt.Sprintf("ssh-reg: Host '%v' doesn't exist.", *copyHost))
		}
		break
	case describe.FullCommand():
		_, exists := entries[*describeHost]
		if exists {
			fmt.Printf(printEntry(entries[*describeHost]))
		} else {
			app.FatalUsageContext(context, fmt.Sprintf("ssh-reg: Host '%v' doesn't exist.", *describeHost))
		}
		break
	case list.FullCommand():
		listHosts()
		break
	case move.FullCommand():
		_, exists := entries[*moveHost]
		if exists {
			entry := entries[*moveHost]
			entries[*moveNewHost] = Host{Host: *moveNewHost, HostName: entry.HostName, IdentityFile: entry.IdentityFile, User: entry.User}
			delete(entries, *moveHost)
			saveEntries()
		} else {
			app.FatalUsageContext(context, fmt.Sprintf("ssh-reg: Host '%v' doesn't exist.", *moveHost))
		}
		break
	case remove.FullCommand():
		_, exists := entries[*removeHost]
		if exists {
			delete(entries, *removeHost)
			saveEntries()
		} else {
			app.FatalUsageContext(context, fmt.Sprintf("ssh-reg: Host '%v' doesn't exist.", *removeHost))
		}
		break
	case update.FullCommand():
		_, exists := entries[*updateHost]
		if exists {
			updateEntry(*updateHost, *updateHostName, *updateIdentityFile, *updateUser, *updateExtra)
		} else {
			app.FatalUsageContext(context, fmt.Sprintf("ssh-reg: Host '%v' doesn't exist.", *updateHost))
		}
		break
	}
}

func parseConfig(input string) {
	entry_regex, _ := regexp.Compile("((.+) (.+)\\s?)+")
	host_option, _ := regexp.Compile("(?:\\s+)?(.+) (.+)\\s?")
	matches := entry_regex.FindAllStringSubmatch(string(input), -1)

	var hosts []string
	for i := 0; i < len(matches); i++ {
		hosts = append(hosts, matches[i][0])
	}

	for _, entry := range hosts {
		options := strings.Split(entry, "\n")
		output := Host{}
		output.Extras = make(map[string]string)

		for _, option := range options {
			if len(option) > 1 {
				option_matches := host_option.FindAllStringSubmatch(option, -1)
				key := option_matches[0][1]
				value := option_matches[0][2]
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
}

func validateExtras(input []string) bool {
	valid_keywords := extraKeywords()
	for _, extra := range input {
		keyword := strings.Split(extra, "=")
		if !stringInSlice(keyword[0], valid_keywords) {
			app.FatalUsageContext(*context, fmt.Sprintf("Invalid Keyword: %s", keyword[0]))
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

func listHosts() {
	var keys []string
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(fmt.Sprintf("%v", k))
	}
}

func printEntry(host Host) string {
	hostTemplate := []string{fmt.Sprintf("Host %v\n", host.Host), fmt.Sprintf("  HostName %v\n", host.HostName)}
	if host.IdentityFile != "" {
		hostTemplate = append(hostTemplate, fmt.Sprintf("  IdentityFile %v\n", host.IdentityFile))
	}
	if host.User != "" {
		hostTemplate = append(hostTemplate, fmt.Sprintf("  User %v\n", host.User))
	}
	if len(host.Extras) > 0 {
		for k, v := range host.Extras {
			hostTemplate = append(hostTemplate, fmt.Sprintf("  %v %v\n", k, v))
		}
	}
	return strings.Join(hostTemplate, "")
}

func addEntry(host string, hostName string, identityFile string, user string, extra string) {
	entry := Host{Host: host, HostName: hostName, IdentityFile: identityFile, User: user}
	entry.Extras = make(map[string]string)

	if extra != "" {
		extra_array := strings.Split(extra, ",")
		if validateExtras(extra_array) {
			for _, v := range extra_array {
				option := strings.Split(v, "=")
				if option[1] != "" {
					entry.Extras[option[0]] = option[1]
				}
			}
		}
	}
	entries[host] = entry
	saveEntries()
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
		extra_array := strings.Split(extra, ",")
		if validateExtras(extra_array) {
			for _, v := range extra_array {
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

func saveEntries() {
	fh, _ := os.OpenFile(ssh_config, os.O_RDWR|os.O_TRUNC, 0777)

	var keys []string
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fh.WriteString(fmt.Sprintf("%v\n", printEntry(entries[k])))
	}
}
