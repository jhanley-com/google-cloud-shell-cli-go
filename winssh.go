package main

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"time"

	"strconv"

	"io"
	"os"
	"runtime"
)

// var path_winssh =  "C:/Windows/System32/OpenSSH/ssh.exe"
var proxy *exec.Cmd

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

	// go ShadowSocks(sshHost, sshPort)

	if config.Proxy == "v2ray" {
		if config.Debug == true {
			fmt.Println("Proxy: V2ray")
		}

		CheckUrlConfig(sshHost, sshPort)
		V2ray(sshHost, sshPort)
		CheckPort("127.0.0.1", "8022")
	} else if config.Proxy == "shadowsocks" {
		if config.Debug == true {
			fmt.Println("Proxy: shadowsocks")
		}

		ShadowSocks(sshHost, sshPort)
		CheckPort("127.0.0.1", "8022")
	} else if config.Proxy != "" {
		config.sshFlags = append(config.sshFlags, "-o", "ProxyCommand=connect.exe -S "+config.Proxy+" %h %p")
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

	client, err := NewExternalClient(sshBinaryPath, sshUsername, sshHost, sshPortInt, &auth, config.sshFlags)
	if err != nil {
		fmt.Println("Failed to create new client - ", err)
		return
	}

	// Set heartbeats
	go client.HeartBeats()

	// "curl -sSL https://raw.githubusercontent.com/ixiumu/google-cloud-shell-cli-go/patch-1/scripts-remote/heartbeats | sh & bash"
	err = client.Shell()
	if err != nil && err.Error() != "exit status 255" {
		fmt.Println("Failed to request shell - ", err)
		// return
	}

	if config.Proxy != "" {
		proxy.Process.Kill()
	}

}

func V2ray(sshHost string, sshPort string) {
	ssBinaryPath, err := exec.LookPath(config.AbsPath + "/v2ray/v2ray.exe")
	if err != nil {
		fmt.Println("Failed to create proxy client - ", err)
		return
	}

	url := "http://" + sshHost + ":" + sshPort + "/api/v2ray.json"

	proxy = exec.Command(ssBinaryPath, "-config", url)

	err = proxy.Start()
	if err != nil {
		fmt.Println("Failed to create proxy client - ", err)
		return
	}

	if config.Debug == true {
		fmt.Println("Proxy Config:", url)
	}

	// config.sshFlags = append(config.sshFlags, "-o", "ProxyCommand=\"C:\\Program Files\\Git\\mingw64\\bin\\connect.exe -S 127.0.0.1:8022 %h %p\"")

	config.sshFlags = append(config.sshFlags, "-o", "ProxyCommand=connect.exe -S 127.0.0.1:8022 %h %p")

	// https://stackoverflow.com/questions/34730941/ensure-executables-called-in-go-process-get-killed-when-process-is-killed

}

func ShadowSocks(sshHost string, sshPort string) {
	// ss-local -s 1.1.1.1 -p 6000 -k 7758521 -m rc4-md5 -l 8022 --plugin v2ray-plugin --plugin-opts "tls;path=/api/shadowsocks/;host=localhost"
	ssArgs := []string{"-m", "rc4-md5", "-l", "8022", "--plugin-opts", "tls;path=/api/shadowsocks/;host=localhost"}

	ssArgs = append(ssArgs, "-k", "7758521")
	ssArgs = append(ssArgs, "-s", sshHost)
	ssArgs = append(ssArgs, "-p", sshPort)

	ssBinaryPath, err := exec.LookPath(config.AbsPath + "/shadowsocks-libev/ss-local.exe")

	// fmt.Println(ssBinaryPath)

	if err != nil {
		fmt.Println("Failed to create proxy client - ", err)
		return
	}

	ssPluginBinaryPath, err := exec.LookPath(config.AbsPath + "/shadowsocks-libev/v2ray-plugin.exe")
	if err != nil {
		fmt.Println("Failed to find proxy plugin - ", err)
		return
	}

	ssArgs = append(ssArgs, "--plugin", ssPluginBinaryPath)

	proxy = exec.Command(ssBinaryPath, ssArgs...)
	// output, err := proxy.CombinedOutput()
	// fmt.Println(string(output))
	// return
	err = proxy.Start()
	if err != nil {
		fmt.Println("Failed to create proxy client - ", err)
		return
	}

	if config.Debug == true {
		fmt.Println("Proxy:", ssArgs)
	}

	config.sshFlags = append(config.sshFlags, "-o", "ProxyCommand=connect.exe -S 127.0.0.1:8022 %h %p")

}

func CheckUrlConfig(sshHost string, sshPort string) {
	for x := 0; x < 60; x++ {
		time.Sleep(500 * time.Millisecond)

		resp, err := http.Get("http://" + sshHost + ":" + sshPort + "/api/v2ray.json")

		if err != nil {
			return
		}

		defer resp.Body.Close()

		// if resp.StatusCode != http.StatusOK {
		// 	if config.Debug == true {
		// 		fmt.Println("Error: status code", resp.StatusCode)
		// 	}
		// 	continue
		// }

		if resp.StatusCode == 200 {
			if config.Debug == true {
				fmt.Println("CheckPort: status code", resp.StatusCode)
			}
			break
		}
	}
}

func CheckPort(host string, port string) {
	for x := 0; x < 60; x++ {
		time.Sleep(500 * time.Millisecond)

		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err != nil {
			if config.Debug == true {
				fmt.Println("Connecting error:", err)
			}
			continue
		}
		if conn != nil {
			defer conn.Close()
			if config.Debug == true {
				fmt.Println("Opened", net.JoinHostPort(host, port))
			}
			break
		}
	}
}

func (client *ExternalClient) HeartBeats() {

	time.Sleep(10 * time.Second)

	if config.Debug == true {
		fmt.Print("Start Send HeartBeats")
	}

	for true {
		args := append(client.BaseArgs, "curl -I -H \"Devshell-Vm-Ip-Address:${DEVSHELL_IP_ADDRESS}\" -X POST -s -w %{http_code} -o /dev/null ${DEVSHELL_SERVER_URL}/devshell/vmheartbeat")

		cmd := getSSHCmd(client.BinaryPath, args...)
		output, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Print("HeartBeats error: ", err)
		}

		if config.Debug == true {
			fmt.Print("HeartBeats: ", string(output))
		}

		time.Sleep(360 * time.Second)
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
		// "-F", "/dev/null", // disable load ~/.ssh/config
		"-o", "ConnectionAttempts=5", // retry 5 times if SSH connection fails
		// "-o", "ConnectTimeout=5", // timeout after 5 seconds
		"-o", "ControlMaster=no", // disable ssh multiplexing
		"-o", "ControlPath=none",
		"-o", "LogLevel=quiet", // suppress "Warning: Permanently added '[localhost]:2022' (ECDSA) to the list of known hosts."
		"-o", "PasswordAuthentication=no",
		"-o", "ServerAliveInterval=30", // prevents connection to be dropped if command takes too long
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
	}
	defaultClientType = External
)

func NewExternalClient(sshBinaryPath, user, host string, port int, auth *Auth, flags []string) (*ExternalClient, error) {
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

	if len(flags) > 0 {
		args = append(args, flags...)
		// fmt.Println(args)
	}

	client.BaseArgs = args

	if config.Debug == true {
		fmt.Println("sshArgs:", client.BaseArgs)
	}

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
		fmt.Println("Error closing SSH Client: ", err)
	}
}
