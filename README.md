# reg
reg {ip} {name} {-i,--ssh-key [ssh-key]} {-u, --user [user]}

reg 192.168.0.1 A
reg 192.168.0.1 A -i ~/.ssh/thrillist -u mpriscella
confirm overwrite (y/n)
reg A
