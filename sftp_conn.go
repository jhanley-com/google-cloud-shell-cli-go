package main

import (
	"fmt"
	"net"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func sftp_open_connection(params CloudShellEnv) (*ssh.Client, *sftp.Client, error) {
	file, err := env_get_ssh_pkey()

	if err != nil {
		fmt.Println("\nTip: Run the command: \"gcloud alpha cloud-shell ssh --dry-run\" to setup Cloud Shell SSH keys")
		return nil, nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User: params.SshUsername,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(file),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	sshHost := params.SshHost
	sshPort := fmt.Sprint(params.SshPort)

	host := sshHost + ":" + sshPort

	if config.Debug == true {
		fmt.Println("Dial: " + host)
	}

	connection, err := ssh.Dial("tcp", host, sshConfig)

	if err != nil {
		// return nil, fmt.Errorf("Failed to dial: %s", err)
		fmt.Println(err)
		return nil, nil, err
	}

	client, err := sftp.NewClient(connection)

	if err != nil {
		// return nil, fmt.Errorf("Failed to dial: %s", err)
		fmt.Println(err)
		return nil, nil, err
	}

	return connection, client, nil
}
