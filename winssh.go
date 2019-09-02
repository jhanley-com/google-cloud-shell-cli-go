package main

import (
	"fmt"
	
	"github.com/docker/machine/libmachine/ssh"
	"strconv"
)

var path_winssh =  "C:/Windows/System32/OpenSSH/ssh.exe"

func exec_winssh(params CloudShellEnv) {
	key, err := env_get_ssh_pkey()

	if err != nil {
		fmt.Println("\nTip: Run the command: \"gcloud alpha cloud-shell ssh --dry-run\" to setup Cloud Shell SSH keys")
		return
	}

	sshUsername := params.SshUsername
	sshHost := params.SshHost
	sshPort := fmt.Sprint(params.SshPort)
	sshUrl := sshUsername + "@" + sshHost

	if config.Debug == true {
		fmt.Println(key)
		fmt.Println(sshUsername)
		fmt.Println(sshHost)
		fmt.Println(sshPort)
		fmt.Println(sshUrl)
	}

	auth := ssh.Auth{Keys: []string{key}}
	sshPortInt, err := strconv.Atoi(sshPort)
	client, err := ssh.NewClient(sshUsername, sshHost, sshPortInt, &auth)
	if err != nil {
		fmt.Errorf("Failed to create new client - %s", err)
		return
	}

	err = client.Shell()
	if err != nil && err.Error() != "exit status 255" {
		fmt.Errorf("Failed to request shell - %s", err)
		return
	}

}
