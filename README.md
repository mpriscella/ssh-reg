# ssh-reg
An SSH config management tool written in go.

# Installation
Download the compiled [binary](https://github.com/mpriscella/ssh-reg/releases/download/v1.1.1/ssh-reg) to a directory in your $PATH (like `/usr/local/bin`).

# Examples
`$ ssh-reg add development dev.mpriscella.com -i ~/.ssh/id_rsa -u mpriscella -e "Port=8989"`

```
$ ssh-reg describe development
> Host development
>   HostName dev.mpriscella.com
>   IdentityFile /home/mpriscella/.ssh/id_rsa
>   User mpriscella
>   Port 8989
```

```
$ ssh-reg update development staging.mpriscella.com -e "Port="
> Host development
>   HostName staging.mpriscella.com
>   IdentityFile /home/mpriscella/.ssh/id_rsa
>   User mpriscella
```

````
usage: ssh-reg [<flags>] <command> [<args> ...]

ssh-reg is a program to manage a user's ssh config file.

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --version  Show application version.

Commands:
  help [<command>...]
    Show help.


  add [<flags>] <host> <hostname>
    Add host entry.

    -i, --identity=IDENTITY  The path to the identity key file.
    -u, --user=USER          The SSH user.
    -f, --force              Overwrite the specified host.
    -e, --extra=EXTRA        Add Extra Keyword.

  copy <source> <destination>
    Copy host entry.


  describe <host>
    Describe host entry.


  list
    List all available host entries.


  move <source> <destination>
    Rename host entry.


  remove <host>
    Remove host entry.


  update [<flags>] <host> [<hostname>]
    Update host entry.

    -i, --identity=IDENTITY  The path to the identity key file.
    -u, --user=USER          The SSH User.
    -e, --extra=EXTRA        Keyword=Value.
````
