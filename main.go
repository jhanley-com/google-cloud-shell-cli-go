package main

// go get github.com/kirinlabs/HttpRequest
// go get github.com/pkg/sftp
// go get golang.org/x/crypto/ssh

import (
	"fmt"
)

// This is the file where User Credentails are saved after authorization
// This credentials are loaded on program start and refreshed if previously saved
var SavedUserCredentials = "user_credentials.json"
var SavedAdcCredentials = "adc_credentials.json"

// If you change the scopes, delete the saved user_credentials.json
var SCOPE = "https://www.googleapis.com/auth/cloud-platform openid https://www.googleapis.com/auth/userinfo.email"

func main() {
	//************************************************************
	//
	//************************************************************

	err := init_config()

	if err != nil {
		return
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
	accessToken, _, err = get_tokens()

	if err != nil {
		return
	}

	if accessToken == "" {
		fmt.Println("Error: Empty Access Token")
		return
	}

	call_cloud_shell(accessToken)
}
