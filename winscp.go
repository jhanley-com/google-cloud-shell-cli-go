package main

import (
	"fmt"
	"os/exec"
)

var path_winscp =  "C:\\Program Files (x86)\\WinSCP\\WinSCP.exe"

func exec_winscp(params CloudShellEnv) {
	key, err := env_get_ssh_ppk()

	if err != nil {
		fmt.Println("\nTip: Run the command: \"gcloud alpha cloud-shell ssh --dry-run\" to setup Cloud Shell SSH keys")
		return
	}

	sshUsername := params.SshUsername
	sshHost := params.SshHost
	sshPort := fmt.Sprint(params.SshPort)
	sshUrl := "sftp://" + sshUsername + "@" + sshHost + ":" + sshPort

	args := append([]string{"/ini=nul", "/privatekey=" + key, "/hostkey=*", sshUrl}, config.WinscpFlags...)

	if config.Debug == true {
		fmt.Println(key)
		fmt.Println(sshUsername)
		fmt.Println(sshHost)
		fmt.Println(sshPort)
		fmt.Println(sshUrl)
		fmt.Println(args)
	}

	path := path_winscp

	cmd := exec.Command(path, args...)

	err = cmd.Start()

	if err != nil {
		fmt.Println(err)
		return
	}
}
