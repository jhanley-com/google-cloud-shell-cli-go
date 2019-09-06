package main

import (
	"fmt"
	"os/exec"

	"strconv"

	"os"
	"io"
	"runtime"
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

	sshBinaryPath, err := exec.LookPath(config.AbsPath + "/ssh")
	if err != nil {
		sshBinaryPath, err = exec.LookPath("ssh")
		if err != nil {
			if runtime.GOOS != "windows" {
				sshBinaryPath = "ssh"
			} else {
				sshBinaryPath = "ssh.exe"
			}
		}
	}

	if config.Debug == true {
		fmt.Println("Use:", sshBinaryPath)
	}

	auth := Auth{Keys: []string{key}}
	sshPortInt, err := strconv.Atoi(sshPort)

	client, err := NewExternalClient(sshBinaryPath, sshUsername, sshHost, sshPortInt, config.Flags.BindAddress, &auth)
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


// "github.com/docker/machine/libmachine/ssh"

type ExternalClient struct {
	BaseArgs   []string
	BinaryPath string
	cmd        *exec.Cmd
}

type Auth struct {
	Passwords []string
	Keys      []string
}

const (
	External ClientType = "external"
	Native   ClientType = "native"
)

type ClientType string

var (
	baseSSHArgs = []string{
		"-F", "/dev/null",
		"-o", "ConnectionAttempts=3", // retry 3 times if SSH connection fails
		"-o", "ConnectTimeout=10", // timeout after 10 seconds
		"-o", "ControlMaster=no", // disable ssh multiplexing
		"-o", "ControlPath=none",
		"-o", "LogLevel=quiet", // suppress "Warning: Permanently added '[localhost]:2022' (ECDSA) to the list of known hosts."
		"-o", "PasswordAuthentication=no",
		"-o", "ServerAliveInterval=60", // prevents connection to be dropped if command takes too long
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
	}
	defaultClientType = External
)

func NewExternalClient(sshBinaryPath, user, host string, port int, bind string, auth *Auth) (*ExternalClient, error) {
	client := &ExternalClient{
		BinaryPath: sshBinaryPath,
	}

	args := append(baseSSHArgs, fmt.Sprintf("%s@%s", user, host))

	// If no identities are explicitly provided, also look at the identities
	// offered by ssh-agent
	if len(auth.Keys) > 0 {
		args = append(args, "-o", "IdentitiesOnly=yes")
	}

	// Specify which private keys to use to authorize the SSH request.
	for _, privateKeyPath := range auth.Keys {
		if privateKeyPath != "" {
			// Check each private key before use it
			fi, err := os.Stat(privateKeyPath)
			if err != nil {
				// Abort if key not accessible
				return nil, err
			}
			if runtime.GOOS != "windows" {
				mode := fi.Mode()
				// log.Debugf("Using SSH private key: %s (%s)", privateKeyPath, mode)
				// Private key file should have strict permissions
				perm := mode.Perm()
				if perm&0400 == 0 {
					return nil, fmt.Errorf("'%s' is not readable", privateKeyPath)
				}
				if perm&0077 != 0 {
					return nil, fmt.Errorf("permissions %#o for '%s' are too open", perm, privateKeyPath)
				}
			}
			args = append(args, "-i", privateKeyPath)
		}
	}

	// Set which port to use for SSH.
	args = append(args, "-p", fmt.Sprintf("%d", port))

	if len(bind) > 0 {
		args = append(args, "-D", bind)
		// fmt.Println(args)
	}

	client.BaseArgs = args

	return client, nil
}

func getSSHCmd(binaryPath string, args ...string) *exec.Cmd {
	return exec.Command(binaryPath, args...)
}

func (client *ExternalClient) Output(command string) (string, error) {
	args := append(client.BaseArgs, command)
	cmd := getSSHCmd(client.BinaryPath, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (client *ExternalClient) Shell(args ...string) error {
	args = append(client.BaseArgs, args...)
	cmd := getSSHCmd(client.BinaryPath, args...)

	// log.Debug(cmd)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (client *ExternalClient) Start(command string) (io.ReadCloser, io.ReadCloser, error) {
	args := append(client.BaseArgs, command)
	cmd := getSSHCmd(client.BinaryPath, args...)

	// log.Debug(cmd)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		if closeErr := stdout.Close(); closeErr != nil {
			return nil, nil, fmt.Errorf("%s, %s", err, closeErr)
		}
		return nil, nil, err
	}
	if err := cmd.Start(); err != nil {
		stdOutCloseErr := stdout.Close()
		stdErrCloseErr := stderr.Close()
		if stdOutCloseErr != nil || stdErrCloseErr != nil {
			return nil, nil, fmt.Errorf("%s, %s, %s",
				err, stdOutCloseErr, stdErrCloseErr)
		}
		return nil, nil, err
	}

	client.cmd = cmd
	return stdout, stderr, nil
}

func (client *ExternalClient) Wait() error {
	err := client.cmd.Wait()
	client.cmd = nil
	return err
}

func closeConn(c io.Closer) {
	err := c.Close()
	if err != nil {
		// log.Debugf("Error closing SSH Client: %s", err)
		// fmt.Errorf("Error closing SSH Client: %s", err)
	}
}
