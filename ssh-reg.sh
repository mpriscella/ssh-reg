#!/bin/bash

SSH_CONFIG="$HOME/.ssh/config"

function reg_help() {
  echo "usage: reg <ip_address> <machine_name> [-i, --ssh-key <path_to_ssh_key>] [-u, --user <username>] [--help]"
  echo ""
}

for opt in $@
do
  case $opt in
    -i)
            shift
            echo $opt
            exit 0
            ;;
    --help)
            reg_help
            exit 0
            ;;
    *) reg_help
  esac
done

function set_config() {
  echo ""
}

if [ -z "$1" ]; then
  reg_help
fi
