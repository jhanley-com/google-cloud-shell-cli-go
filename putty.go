package main

import (
	"fmt"
	"os/exec"
)

func exec_putty(params CloudShellEnv) {
	key, err := env_get_ssh_ppk()

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

	cmd := exec.Command("putty.exe", "-t", "-P", sshPort, "-i", key, sshUrl)

	err = cmd.Start()

	if err != nil {
		fmt.Println(err)
		return
	}
}
