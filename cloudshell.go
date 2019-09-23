package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/kirinlabs/HttpRequest"
)

//******************************************************************************************
// Cloud Shell State
//
// https://cloud.google.com/shell/docs/reference/rest/Shared.Types/State
//
// STATE_UNSPECIFIED	The environment's states is unknown.
// DISABLED		The environment is not running and can't be connected to.
//			Starting the environment will transition it to the STARTING state.
// STARTING		The environment is being started but is not yet ready to accept
//			connections.
// RUNNING		The environment is running and ready to accept connections. It
//			will automatically transition back to DISABLED after a period of
//			inactivity or if another environment is started.
//******************************************************************************************

//******************************************************************************************
// https://cloud.google.com/shell/docs/reference/rest/Shared.Types/Environment
//******************************************************************************************

type CloudShellEnv struct {
	Name        string `json:"name"`
	Id          string `json:"id"`
	DockerImage string `json:"dockerImage"`
	State       string `json:"state"`
	SshUsername string `json:"sshUsername"`
	SshHost     string `json:"sshHost"`
	SshPort     int32  `json:"sshPort"`
	Error       struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

//******************************************************************************************
// Method: users.environments.get
// https://cloud.google.com/shell/docs/reference/rest/v1alpha1/users.environments/get
//******************************************************************************************

func cloud_shell_get_environment(accessToken string, flag_info bool) (CloudShellEnv, error) {

	var params CloudShellEnv

	endpoint := "https://cloudshell.googleapis.com/v1alpha1/users/me/environments/default"
	endpoint += "?alt=json"

	if config.UrlFetch != "" {
		endpoint = config.UrlFetch + endpoint
	}

	req := HttpRequest.NewRequest()

	req.SetHeaders(map[string]string{
		"Authorization":       "Bearer " + accessToken,
		"X-Goog-User-Project": config.ProjectId})

	if config.Debug == true {
		fmt.Println("Access Token:", accessToken)
		fmt.Println("ProjectId:", config.ProjectId)
	}

	//************************************************************
	//
	//************************************************************

	res, err := req.Get(endpoint)

	if err != nil {
		fmt.Println("Error: ", err)
		return params, err
	}

	body, err := res.Body()

	if err != nil {
		fmt.Println("Error: ", err)
		return params, err
	}

	if flag_info == true {
		fmt.Println("Cloud Shell Info:")
		fmt.Println(string(body))
	}

	err = json.Unmarshal(body, &params)

	if err != nil {
		fmt.Println("Error: Cannot unmarshal JSON: ", err)
		return params, err
	}

	if params.Error.Code != 0 {
		fmt.Println("")
		fmt.Println(params.Error.Message)
	}

	return params, nil
}

//******************************************************************************************
// Method: users.environment.start
// https://cloud.google.com/shell/docs/reference/rest/v1alpha1/users.environments/start
//******************************************************************************************

func cloudshell_start(accessToken string) error {
	//************************************************************
	//
	//************************************************************

	if config.Debug == true {
		fmt.Println("Request users.environment.start")
	}

	endpoint := "https://cloudshell.googleapis.com/v1alpha1/users/me/environments/default"
	endpoint += ":start"
	endpoint += "?alt=json"

	if config.UrlFetch != "" {
		endpoint = config.UrlFetch + endpoint
	}

	req := HttpRequest.NewRequest()

	req.SetHeaders(map[string]string{
		"Authorization":       "Bearer " + accessToken,
		"X-Goog-User-Project": config.ProjectId})

	res, err := req.JSON().Post(endpoint, "{\"accessToken\": \""+accessToken+"\"}")

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	body, err := res.Body()

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	if config.Debug == true {
		fmt.Println("")
		fmt.Println("************************************************************")
		fmt.Println("Cloud Shell Info:")
		fmt.Println(string(body))
		fmt.Println("************************************************************")
	}

	var params CloudShellEnv

	err = json.Unmarshal(body, &params)

	if err != nil {
		fmt.Println("Error: Cannot unmarshal JSON: ", err)
		return err
	}

	if params.Error.Code != 0 {
		err = errors.New(params.Error.Message)
		fmt.Println("Error Code:", params.Error.Code)
		fmt.Println(err)
		return err
	}

	return nil
}

func cloud_shell_create_publickeys(accessToken string) error {

	fmt.Println("Pushing your public key to Cloud Shell...")

	endpoint := "https://cloudshell.googleapis.com/v1alpha1/users/me/environments/default/publicKeys"
	endpoint += "?alt=json"

	if config.UrlFetch != "" {
		endpoint = config.UrlFetch + endpoint
	}

	req := HttpRequest.NewRequest()

	req.SetHeaders(map[string]string{
		"Authorization":       "Bearer " + accessToken,
		"X-Goog-User-Project": config.ProjectId})

	if config.Debug == true {
		fmt.Println("Access Token:", accessToken)
		fmt.Println("ProjectId:", config.ProjectId)
	}

	//************************************************************
	//
	//************************************************************

	path, err := env_get_ssh_pub_key()
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	keyFields := strings.Fields(string(bytes))

	pubkey := keyFields[1]

	res, err := req.JSON().Post(endpoint, "{\"key\":{\"format\": \"SSH_RSA\",\"key\": \""+pubkey+"\"}}")
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	body, err := res.Body()
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	type Params struct {
		Name   string `json:"name"`
		Format string `json:"format"`
		Key    string `json:"key"`
		Error  struct {
			Code    int32  `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	var params Params

	err = json.Unmarshal(body, &params)

	if err != nil {
		fmt.Println("Error: Cannot unmarshal JSON: ", err)
		return err
	}

	if config.Debug == true {
		fmt.Println("Body:", params)
	}

	if params.Error.Code != 0 {
		fmt.Println("Error:", params.Error.Message)
	}

	return nil
}

func env_get_ssh_key(filename string, ext string) (string, error) {
	//*************************************************************
	// Return the Google Cloud SSH Key for the current User
	//*************************************************************

	// var path string
	homepath, err := get_home_directory()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// if isWindows() == true {
	// 	path = homepath + "\\.ssh\\" + filename + ext
	// } else {
	// 	path = homepath + "/.ssh/" + filename + ext
	// }

	path := homepath + "/.ssh/" + filename + ext
	if fileExists(path) == false {
		path = homepath + "/.ssh/id_rsa_cloudshell" + ext
		if fileExists(path) == false {
			path = homepath + "/.ssh/id_rsa" + ext
			if fileExists(path) == false {
				err = errors.New("SSH Key does not exist")
				fmt.Println("Error:", err)
				fmt.Println("File:", homepath+"/.ssh/id_rsa"+ext)
				fmt.Println("Please create and upload SSH key first")
				fmt.Println("\n\n$ ssh-keygen -t rsa -C \"Cloud Shell\"")
				fmt.Println("$ cloudshell push_pubkey\n\n ")
				return "", err
			}
		}
	}

	if config.Debug == true {
		fmt.Println("Path:", path)
	}

	return path, nil
}

func env_get_ssh_pkey() (string, error) {
	return env_get_ssh_key("google_compute_engine", "")
}

func env_get_ssh_ppk() (string, error) {
	return env_get_ssh_key("google_compute_engine", ".ppk")
}

func env_get_ssh_pub_key() (string, error) {
	return env_get_ssh_key("google_compute_engine", ".pub")
}

func call_cloud_shell(accessToken string) {
	//************************************************************
	//
	//************************************************************

	if config.Command == CREATE_PUBKEY {
		err := cloud_shell_create_publickeys(accessToken)
		if err != nil {
			fmt.Println("Error:", err)
		}
		os.Exit(0)
	}

	flag_info := false

	if config.Command == CMD_EXEC {
		flag_info = false
	}

	if config.Debug == true {
		flag_info = true
	}

	if config.Command == CMD_INFO {
		flag_info = true
	}

	var params CloudShellEnv

	params, err := cloud_shell_get_environment(accessToken, flag_info)

	if err != nil {
		return
	}

	if config.Command == CMD_INFO {
		return
	}
	if params.State == "DISABLED" {
		if config.Debug == true {
			fmt.Println("CloudShell State:", params.State)
		}

		fmt.Println("Starting your Cloud Shell machine...")

		err = cloudshell_start(accessToken)

		if err != nil {
			return
		}

		for x := 0; x < 60; x++ {
			time.Sleep(500 * time.Millisecond)

			params, err = cloud_shell_get_environment(accessToken, flag_info)

			if err != nil {
				return
			}

			if params.Error.Code != 0 {
				return
			}

			if params.State == "RUNNING" {
				fmt.Println("Waiting for your Cloud Shell machine to start...")
				// Increase waiting time
				// time.Sleep(5000 * time.Millisecond)

				break
			}
		}

		// waiting
		host := params.SshHost
		port := fmt.Sprint(params.SshPort)

		for x := 0; x < 120; x++ {
			time.Sleep(1000 * time.Millisecond)

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

	if params.State != "RUNNING" {
		fmt.Println("CloudShell State:", params.State)
		return
	}

	if config.Command == CMD_WINSSH {
		fmt.Println("Your Cloud Shell machine is RUNNING, connecting...")
		exec_winssh(params)
	}

	if config.Command == CMD_SSH {
		exec_ssh(params)
	}

	if config.Command == CMD_EXEC {
		exec_command(params)
	}

	if config.Command == CMD_DOWNLOAD {
		sftp_download(params)
	}

	if config.Command == CMD_UPLOAD {
		sftp_upload(params)
	}

	if config.Command == CMD_PUTTY {
		exec_putty(params)
	}

	if config.Command == CMD_WINSCP {
		exec_winscp(params)
	}
}
