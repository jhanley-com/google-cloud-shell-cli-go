package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func sftp_download(params CloudShellEnv) {
	file := env_get_ssh_pkey()

	if file == "" {
		fmt.Println("Error: Cannot get SSH private key file")
		return
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
		fmt.Println("Dial")
		fmt.Println(host)
	}

	connection, err := ssh.Dial("tcp", host, sshConfig)

	if err != nil {
		// return nil, fmt.Errorf("Failed to dial: %s", err)
		fmt.Println(err)
		return
	}

	defer connection.Close()

	client, err := sftp.NewClient(connection)

	if err != nil {
		// return nil, fmt.Errorf("Failed to dial: %s", err)
		fmt.Println(err)
		return
	}

	defer client.Close()

	// open source file
	fmt.Println("open source file")
	srcFile, err := client.Open(config.SrcFile)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()
 
	// create destination file
	fmt.Println("create destination file")
	dstFile, err := os.Create(config.DstFile)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	// copy source file to destination file
	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d bytes copied\n", bytes)
}

func sftp_upload(params CloudShellEnv) {
	file := env_get_ssh_pkey()

	if file == "" {
		fmt.Println("Error: Cannot get SSH private key file")
		return
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
		fmt.Println("Dial")
		fmt.Println(host)
	}

	connection, err := ssh.Dial("tcp", host, sshConfig)

	if err != nil {
		// return nil, fmt.Errorf("Failed to dial: %s", err)
		fmt.Println(err)
		return
	}

	defer connection.Close()

	client, err := sftp.NewClient(connection)

	if err != nil {
		// return nil, fmt.Errorf("Failed to dial: %s", err)
		fmt.Println(err)
		return
	}

	defer client.Close()

	// open source file
	srcFile, err := os.Open(config.SrcFile)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	// create destination file
	dstFile, err := client.Create(config.DstFile)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()
 
	// copy source file to destination file
	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d bytes copied\n", bytes)
}
