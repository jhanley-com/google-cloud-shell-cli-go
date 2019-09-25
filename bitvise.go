package main

import (
	"fmt"
	"os/exec"
)

// Bitvise parameter list
// https://www.bitvise.com/files/guides/BvSshClient-Params.txt
// Also supports loading a profile
// BvSsh.exe -profile=<full-path-to-profile.tlp>

var path_bitvise =  "C:\\Program Files (x86)\\Bitvise SSH Client\\BvSsh.exe"

func exec_bitvise(params CloudShellEnv) {
	key, err := env_get_ssh_ppk()

	if err != nil {
		fmt.Println("\nTip: Run the command: \"gcloud alpha cloud-shell ssh --dry-run\" to setup Cloud Shell SSH keys")
		return
	}

	sshUsername := params.SshUsername
	sshHost := params.SshHost
	sshPort := fmt.Sprint(params.SshPort)
	sshUrl := "sftp://" + sshUsername + "@" + sshHost + ":" + sshPort

	args := []string{}

	args = append(args, "-host=" + sshHost)
	args = append(args, "-port=" + sshPort)
	args = append(args, "-user=" + sshUsername)
	args = append(args, "-keypairFile=" + key)
	args = append(args, "-loginOnStartup")

	if config.Debug == true {
		fmt.Println(key)
		fmt.Println(sshUsername)
		fmt.Println(sshHost)
		fmt.Println(sshPort)
		fmt.Println(sshUrl)
		fmt.Println(args)
	}

	path := path_bitvise

	cmd := exec.Command(path, args...)

	err = cmd.Start()

	if err != nil {
		fmt.Println(err)
		return
	}
}
