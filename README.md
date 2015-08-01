# ssh-reg
An SSH config management tool written in go.

# Installation
Download the compiled [binary](https://github.com/mpriscella/ssh-reg/releases/download/v0.5.1/ssh-reg) to a directory in your $PATH (like `/usr/local/bin`).

# Examples
`$ ssh-reg add development dev.mpriscella.com -i ~/.ssh/id_rsa -u mpriscella`

```
$ ssh-reg describe development
> Host development
>   HostName dev.mpriscella.com
>   IdentityFile /home/mpriscella/.ssh/id_rsa
>   User mpriscella
```

```
$ ssh-reg update development staging.mpriscella.com
> Host development
>   HostName staging.mpriscella.com
>   IdentityFile /home/mpriscella/.ssh/id_rsa
>   User mpriscella
```

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


  describe <host>
    Describe host


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
