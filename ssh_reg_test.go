package main

import (
	"fmt"
	"testing"
)

var test_config_file = `Host mike
  HostName mpriscella.com
  IdentityFile /path/to/ssh/file
  User mpriscella
	
Host test
  HostName
  IdentityFile

Host something
  HostName
  User
	
Host test
  HostName`

func TestDescribe(t *testing.T) {

	t.Log(fmt.Sprintf("\n%v", test_config_file))
}
