package main

// go get github.com/kirinlabs/HttpRequest
// go get github.com/pkg/sftp
// go get golang.org/x/crypto/ssh
// go get golang.org/x/oauth2/google

import (
	"fmt"
	"os"
)

// This is the file where User Credentails are saved after authorization
// This credentials are loaded on program start and refreshed if previously saved
// ~/.config/cloudshell/user_credentials.json
var SavedUserCredentials = user_config_path("user_credentials.json")

// If you change the scopes, delete the saved user_credentials.json
// var SCOPE = "https://www.googleapis.com/auth/cloud-platform openid https://www.googleapis.com/auth/userinfo.email"
var SCOPE = "https://www.googleapis.com/auth/cloud-platform"

func main() {
	//************************************************************
	//
	//************************************************************

	err := init_config()

	if err != nil {
		os.Exit(1)
	}

	// vscode
	if config.Command == CMD_SSH_VSCODE {
		exec_vscode_ssh()
		os.Exit(1)
	}

	//************************************************************
	//
	//************************************************************

	err = check_os()

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	//************************************************************
	// Using Cloud SDK User Credentials does not work with Cloud Shell
	//
	// Your application has authenticated using end user credentials from the
	// Google Cloud SDK or Google Cloud Shell which are not supported by the
	// cloudshell.googleapis.com. We recommend that most server applications
	// use service accounts instead. For more information about service accounts
	// and how to use them in your application, see
	// https://cloud.google.com/docs/authentication/.
	//
	// Service Account credentials do not work with Cloud Shell, so we
	// must use OAuth 2.0 User Credentials
	//************************************************************

	var accessToken = ""

	// accessToken, idToken, err := get_tokens()
	accessToken, err = get_tokens()

	if err != nil {
		os.Exit(1)
	}

	if accessToken == "" {
		fmt.Println("Error: Empty Access Token")
		os.Exit(1)
	}

	// FIX - return error numbers so that scripts/tools can detect errors
	call_cloud_shell(accessToken)

	os.Exit(0)
}
