package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	Name                 string   `json:"name"`
	Id                   string   `json:"id"`
	DockerImage          string   `json:"dockerImage"`
	State                string   `json:"state"`
	SshUsername          string   `json:"sshUsername"`
	SshHost              string   `json:"sshHost"`
	SshPort              int32   `json:"sshPort"`
	Error struct {
		Code              int32   `json:"code"`
		Message           string  `json:"message"`
		Status            string  `json:"status"`
	} `json:"error"`
}

//******************************************************************************************
// Method: users.environments.get
// https://cloud.google.com/shell/docs/reference/rest/v1alpha1/users.environments/get
//******************************************************************************************

func cloud_shell_get_environment(accessToken string, flag_info bool) (CloudShellEnv, error) {
	//************************************************************
	//
	//************************************************************

	var params CloudShellEnv

	endpoint := "https://cloudshell.googleapis.com/v1alpha1/users/me/environments/default"
	endpoint += "?alt=json"

	req := HttpRequest.NewRequest()

	req.SetHeaders(map[string]string{
			"Authorization": "Bearer " + accessToken,
			"X-Goog-User-Project": config.ProjectId})

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
		fmt.Println("")
		fmt.Println("************************************************************")
		fmt.Println("Cloud Shell Info:")
		fmt.Println(string(body))
		fmt.Println("************************************************************")
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

func cloudshell_start(accessToken string) {
	//************************************************************
	//
	//************************************************************

	fmt.Println("Starting Cloud Shell")

	endpoint := "https://cloudshell.googleapis.com/v1alpha1/users/me/environments/default"
	endpoint += ":start"
	endpoint += "?alt=json"

	req := HttpRequest.NewRequest()

	// req.Header.Set("Authorization", "Bearer " + accessToken)
	// req.Header.Set("X-Goog-User-Project", config.ProjectId)

	req.SetHeaders(map[string]string{
			"Authorization": "Bearer " + accessToken,
			"X-Goog-User-Project": config.ProjectId})

	//************************************************************
	//
	//************************************************************

	res, err := req.Post(endpoint)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	body, err := res.Body()

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("")
	fmt.Println("************************************************************")
	fmt.Println("Cloud Shell Info:")
	fmt.Println(string(body))
	fmt.Println("************************************************************")

	return
}

func env_get_ssh_pkey() string {
	//*************************************************************
	// Return the Google Cloud SSH Key for the current Windows User
	//*************************************************************

	home_drive := os.Getenv("HOMEDRIVE")
	home_path := os.Getenv("HOMEPATH")

	if home_drive == "" {
		fmt.Println("Error: Missing Environment Variable: HOMEDRIVE")
		return ""
	}

	if home_path == "" {
		fmt.Println("Error: Missing Environment Variable: HOMEPATH")
		return ""
	}

	path := home_drive + home_path

	path += "\\.ssh\\google_compute_engine"

	if config.Debug == true {
		fmt.Println("Path:", path)
	}

	if fileExists(path) == false {
		fmt.Println("Error: Google SSH Key does not exist")
		fmt.Println("File:", path)
		return ""
	}

	return path
}

func env_get_ssh_ppk() string {
	//*************************************************************
	// Return the Google Cloud SSH Key for the current Windows User
	//*************************************************************

	home_drive := os.Getenv("HOMEDRIVE")
	home_path := os.Getenv("HOMEPATH")

	if home_drive == "" {
		fmt.Println("Error: Missing Environment Variable: HOMEDRIVE")
		return ""
	}

	if home_path == "" {
		fmt.Println("Error: Missing Environment Variable: HOMEPATH")
		return ""
	}

	path := home_drive + home_path

	path += "\\.ssh\\google_compute_engine.ppk"

	if config.Debug == true {
		fmt.Println("Path:", path)
	}

	if fileExists(path) == false {
		fmt.Println("Error: Google SSH Key does not exist")
		fmt.Println("File:", path)
		return ""
	}

	return path
}

func call_cloud_shell(accessToken string) {
	//************************************************************
	//
	//************************************************************

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
		fmt.Println("CloudShell State:", params.State)

		cloudshell_start(accessToken)

		for x := 0; x < 20; x++ {
			time.Sleep(500 * time.Millisecond)

			params, err = cloud_shell_get_environment(accessToken, flag_info)

			if err != nil {
				return
			}

			if params.State == "RUNNING" {
				break;
			}
		}
	}

	if params.State != "RUNNING" {
		fmt.Println("CloudShell State:", params.State)
		return
	}

	if config.Command == CMD_PUTTY {
		exec_putty(params)
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
}