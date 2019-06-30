package main

import (
	"fmt"
	"os/exec"
)

var path_winssh =  "C:/Windows/System32/OpenSSH/ssh.exe"

func exec_winssh(params CloudShellEnv) {
	key := env_get_ssh_pkey()

	if key == "" {
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

	cmd := exec.Command("cmd.exe", "/C", "start", path_winssh, sshUrl, "-p", sshPort, "-i", key)

	err := cmd.Start()

	if err != nil {
		fmt.Println(err)
		return
	}
}
