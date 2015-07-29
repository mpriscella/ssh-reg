# ssh-reg

````
usage: ssh-reg [<flags>] <command> [<args> ...]

A ssh config management tool.

Flags:
  --help  Show help (also see --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.


  list
    List all available hosts


  add [<flags>] <host> <hostname>
    Add host

    -i, --identity=IDENTITY
                     The location of the hosts private key
    -u, --user=USER  The SSH User
    -f, --force      Overwrite the specified host

  remove <host>
    Remove host


  update [<flags>] <host> [<hostname>]
    Update host

    -i, --identity=IDENTITY
                     The location of the hosts private key
    -u, --user=USER  The SSH User
````
