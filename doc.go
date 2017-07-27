// usage: ssh-reg [<flags>] <command> [<args> ...]
//
// ssh-reg is a program to manage a user's ssh config file.
//
// Flags:
//   --help     Show context-sensitive help (also try --help-long and --help-man).
//   --version  Show application version.
//
// Commands:
//   help [<command>...]
//     Show help.
//
//
//   add [<flags>] <host> <hostname>
//     Add host entry.
//
//     -i, --identity=IDENTITY  The path to the identity key file.
//     -u, --user=USER          The SSH user.
//     -f, --force              Overwrite the specified host.
//     -e, --extra=EXTRA        Add Extra Keyword.
//
//   copy <source> <destination>
//     Copy host entry.
//
//
//   describe <host>
//     Describe host entry.
//
//
//   list
//     List all available host entries.
//
//
//   move <source> <destination>
//     Rename host entry.
//
//
//   remove <host>
//     Remove host entry.
//
//
//   update [<flags>] <host> [<hostname>]
//     Update host entry.
//
//     -i, --identity=IDENTITY  The path to the identity key file.
//     -u, --user=USER          The SSH User.
//     -e, --extra=EXTRA        Keyword=Value.
package main
