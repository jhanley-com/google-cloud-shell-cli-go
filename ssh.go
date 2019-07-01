package main

// https://github.com/kirinlabs/HttpRequest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"golang.org/x/crypto/ssh"
)

func PublicKeyFile(file string) ssh.AuthMethod {

	buffer, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return ssh.PublicKeys(key)
}

func exec_command(params CloudShellEnv) {
	file, err := env_get_ssh_pkey()

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

	// fmt.Println("Dial")
	// fmt.Println(host)

	connection, err := ssh.Dial("tcp", host, sshConfig)

	if err != nil {
		// return nil, fmt.Errorf("Failed to dial: %s", err)
		fmt.Println(err)
		return
	}

	defer connection.Close()

	if config.Debug == true {
		fmt.Println("Connect")
	}

	session, err := connection.NewSession()

	if err != nil {
		// return nil, fmt.Errorf("Failed to create session: %s", err)
		fmt.Println(err)
		return
	}

	defer session.Close()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	if config.Debug == true {
		fmt.Println("Run Command:", config.RemoteCommand)
	}

	session.Run(config.RemoteCommand)

	fmt.Printf("%s\n", stdoutBuf.String())
	fmt.Printf("%s\n", stderrBuf.String())
}

func exec_ssh(params CloudShellEnv) {
	key, err := env_get_ssh_pkey()

	if err != nil {
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

	cmd := exec.Command("ssh", "-p", sshPort, "-i", key, sshUrl)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		fmt.Println(err)
		return
	}
}
