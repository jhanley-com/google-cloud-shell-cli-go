package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
	"github.com/kirinlabs/HttpRequest"
)

// Web browser to launch to authenticate
// This path is valid for Windows x64 only
// FIX - Test for Windows x86
var CHROME = "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe"

var ENDPOINT = "https://accounts.google.com/o/oauth2/v2/auth"

type ClientSecrets struct {
	Installed struct {
		ClientID                string   `json:"client_id"`
		ProjectID               string   `json:"project_id"`
		AuthURI                 string   `json:"auth_uri"`
		TokenURI                string   `json:"token_uri"`
		AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
		ClientSecret            string   `json:"client_secret"`
		RedirectUris            []string `json:"redirect_uris"`
	} `json:"installed"`
}

type UserCredentials struct {
	ClientID                string   `json:"client_id"`
	ClientSecret            string   `json:"client_secret"`
	RefreshToken            string   `json:"refresh_token"`
	Scope            	string   `json:"scope"`
	Type             	string   `json:"type"`
	// The following two fields are option and exist after authentication
	AccessToken		string   `json:"access_token"`
	IDToken			string   `json:"id_token"`
	Email			string	 `json:"email"`
	ExpiresAt		int64	 `json:"expires_at"`
}

type OAuthTokens struct {
	AccessToken		string   `json:"access_token"`
	ExpiresIn               int      `json:"expires_in"`
	RefreshToken            string   `json:"refresh_token"`
	Scope            	string   `json:"scope"`
	TokenType            	string   `json:"token_type"`
	IDToken			string   `json:"id_token"`
	Error			string   `json:"error"`
	ErrorDescription	string   `json:"error_description"`
}

func readCredentials(filename string) ([]byte, error) {
	in, err := os.Open(filename)

	if err != nil {
		return []byte(""), err
	}

	defer in.Close()

	b, err := ioutil.ReadAll(in)

	return b, err
}

func loadClientSecrets(filename string) (ClientSecrets, error) {
	var secrets ClientSecrets

	data, err := readCredentials(filename)

	if err != nil {
		return secrets, err
	}

	// fmt.Println(string(data))

	err = json.Unmarshal(data, &secrets)

	if err != nil {
		fmt.Println("Error: Cannot unmarshal JSON: ", err)
		fmt.Println("File:", filename)
		return secrets, err
	}

	config.ProjectId = secrets.Installed.ProjectID

	// fmt.Println("ClientID:", secrets.Installed.ClientID)

	return secrets, err
}

func loadUserCredentials(filename string) (UserCredentials, error) {
	var secrets UserCredentials

	data, err := readCredentials(filename)

	if err != nil {
		return secrets, err
	}

	// fmt.Println(string(data))

	err = json.Unmarshal(data, &secrets)

	if err != nil {
		fmt.Println("Error: Cannot unmarshal JSON: ", err)
		fmt.Println("File:", filename)
		return secrets, err
	}

	// fmt.Println("ClientID:", secrets.ClientID)

	return secrets, err
}

func saveUserCredentials(filename string, creds UserCredentials) error {
	fmt.Println("Save Credentials to:", filename)

	j, err := json.MarshalIndent(creds, "", " ")

	if err != nil {
		fmt.Println("Error: Cannot marshall JSON:", err)
		return err
	}

	// err = ioutil.WriteFile(filename + ".test", j, 0644)
	err = ioutil.WriteFile(filename, j, 0644)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func debug_PrintUserCredentials(creds UserCredentials) {
	fmt.Println("************************************************************")
	fmt.Println("ClientID:", creds.ClientID)
	fmt.Println("ClientSecret:", creds.ClientSecret)
	fmt.Println("RefreshToken:", creds.RefreshToken)
	fmt.Println("Scope:", creds.Scope)
	fmt.Println("Type:", creds.Type)
	fmt.Println("AccessToken:", creds.AccessToken)
	fmt.Println("IDToken:", creds.IDToken)
	fmt.Println("ExpiresAt:", creds.ExpiresAt)

	fmt.Println("Expires At:", time.Unix(creds.ExpiresAt, 0))

	var t time.Time = time.Unix(creds.ExpiresAt, 0)
	var expires_in int64 = 0

	if time.Now().Before(t) {
		expires_in = int64(creds.ExpiresAt) - int64(time.Now().UTC().Unix())
		fmt.Println("Expires In:", expires_in)
	} else {
		fmt.Println("Expires In: Expired")
	}
	fmt.Println("************************************************************")
}

func doRefresh(filename string) (string, string, bool) {
	endpoint := "https://www.googleapis.com/oauth2/v4/token"

	creds, err := loadUserCredentials(filename)

	if err != nil {
		fmt.Println(err)
		return "", "", false
	}

	// debug_PrintUserCredentials(creds)

	// We want an access token that is good for a while.
	// Brand new tokens are valid for 3600 seconds
	// For testing require 15 minutes or 900 seconds

	var t time.Time = time.Unix(creds.ExpiresAt - (15 * 60), 0)
	// fmt.Println(t)

	// fmt.Println(time.Now())

	if time.Now().Before(t) {
		if config.Debug == true {
			fmt.Println("Saved credentials (Access Token) have not expired")
		}

		return creds.AccessToken, creds.IDToken, true
	}

	if config.Debug == true {
		fmt.Println("Must Refresh Token")
	}

	content := "client_id=" + creds.ClientID + "&"
	content += "client_secret=" + creds.ClientSecret + "&"
	content += "grant_type=refresh_token&"
	content += "refresh_token=" + creds.RefreshToken

	req := HttpRequest.NewRequest()

	req.SetHeaders(map[string]string{"Content-Type": "application/x-www-form-urlencoded"})

	res, err := req.Post(endpoint, content)

	if err != nil {
		fmt.Println("Error: ", err)
		return "", "", false
	}

	body, err := res.Body()

	if err != nil {
		fmt.Println("Error: ", err)
		return "", "", false
	}

	var tokens OAuthTokens

	err = json.Unmarshal(body, &tokens)

	if err != nil {
		fmt.Println("Error: Cannot unmarshal JSON: ", err)
		return "", "", false
	}

	var expires_at int64 = int64(time.Now().UTC().Unix()) + int64(tokens.ExpiresIn)

/*
	fmt.Println("AccessToken:", tokens.AccessToken)
	fmt.Println("ExpiresIn:", tokens.ExpiresIn)
	fmt.Println("ExpiresAt:", expires_at)
	fmt.Println("Scope:", tokens.Scope)
	fmt.Println("TokenType:", tokens.TokenType)
	fmt.Println("IDToken:", tokens.IDToken)
*/

	creds.AccessToken = tokens.AccessToken
	creds.IDToken = tokens.IDToken
	creds.ExpiresAt = expires_at

	email, err := get_email_address(tokens.AccessToken)

	if err == nil {
		if config.Debug == true {
			fmt.Println("Email:", email)
		}

		creds.Email = email
	}

	err = saveUserCredentials(filename, creds)

	if err != nil {
		fmt.Println("Error: Cannot save user credentials: ", err)
		return "", "", false
	}

	return creds.AccessToken, creds.IDToken, true
}

func debug_displayAccessToken(accessToken string) {
	endpoint := "https://www.googleapis.com/oauth2/v3/tokeninfo"

	req := HttpRequest.NewRequest()

	req.SetHeaders(map[string]string{"Authorization": "Bearer " + accessToken})

	res, err := req.Get(endpoint)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	body, err := res.Body()

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println(string(body))
}

func debug_displayUserInfo(accessToken string) {
	endpoint := "https://www.googleapis.com/oauth2/v3/userinfo"

	req := HttpRequest.NewRequest()

	req.SetHeaders(map[string]string{"Authorization": "Bearer " + accessToken})

	res, err := req.Get(endpoint)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	body, err := res.Body()

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println(string(body))
}

func debug_displayIDToken(accessToken, idToken string) {
	endpoint := "https://www.googleapis.com/oauth2/v3/tokeninfo"

	endpoint += "?id_token=" + idToken

	req := HttpRequest.NewRequest()

	req.SetHeaders(map[string]string{"Authorization": "Bearer " + accessToken})

	res, err := req.Get(endpoint)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	body, err := res.Body()

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println(string(body))
}

func get_email_address(accessToken string) (string, error) {
	type Access_Token struct {
		Azp		string   `json:"azp"`
		Aud		string   `json:"aud"`
		Sub		string   `json:"sub"`
		Scope		string   `json:"scope"`
		Exp		string   `json:"exp"`
		Expires_in	string   `json:"expires_in"`
		Email		string   `json:"email"`
		Email_verified	string   `json:"email_verified"`
		Access_type	string   `json:"access_type"`
	}

	//************************************************************
	//
	//************************************************************

	endpoint := "https://www.googleapis.com/oauth2/v3/tokeninfo"

	req := HttpRequest.NewRequest()

	req.SetHeaders(map[string]string{"Authorization": "Bearer " + accessToken})

	//************************************************************
	//
	//************************************************************

	res, err := req.Get(endpoint)

	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}

	body, err := res.Body()

	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}

	//************************************************************
	//
	//************************************************************

	var tokens Access_Token

	err = json.Unmarshal(body, &tokens)

	if err != nil {
		fmt.Println("Error: Cannot unmarshal JSON: ", err)
		return "", err
	}

	return tokens.Email, nil
}

func get_tokens() (string, string, error) {
	//************************************************************
	//
	//************************************************************

	// fmt.Println("Auth:", config.Flags.Auth)
	// fmt.Println("Login:", config.Flags.Login)

	if config.Flags.Auth == false {
		if fileExists(SavedUserCredentials) {
			accessToken, idToken, valid := doRefresh(SavedUserCredentials)

			if valid == true {
				// fmt.Println("Access Token: ", accessToken)
				// fmt.Println("ID Token:     ", idToken)

				// debug_displayAccessToken(accessToken)
				// debug_displayUserInfo(accessToken)
				// debug_displayIDToken(accessToken, idToken)


				return accessToken, idToken, nil
			}
		}
	}

	//************************************************************
	// FIX - For Linux, this code does not yet support launching
	// a web browser to authenticate with Google
	//************************************************************

	if isWindows() == false {
		err := errors.New("Cannot launch Google Chrome on Linux")
		fmt.Println("Error:", err)
		return "", "", err
	}

	secrets, err := loadClientSecrets(config.ClientSecretsFile)

	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	//************************************************************
	url := ENDPOINT
	url += "?client_id=" + secrets.Installed.ClientID
	url += "&response_type=code"
	url += "&scope=" + SCOPE
	url += "&access_type=offline"
	if len(config.Flags.Login) != 0 {
		url += "&login_hint=" + config.Flags.Login
	}
	url += "&redirect_uri=http://localhost:9000"
	//************************************************************

	cmd := exec.Command(CHROME, url)

	err = cmd.Start()

	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	fmt.Println("Chrome running")

	//************************************************************
	// Start the web server
	//
	// FIX: This is coded in Python.
	//************************************************************

	fmt.Println("Web server starting")

	out, err := exec.Command("python", "webserver.py").Output()

	if err != nil {
		fmt.Println("Error: Web server failed to start")
		fmt.Println(err)
		return "", "", err
	}

	if len(out) == 0 {
		fmt.Println("************************************************************")
		fmt.Println(out)
		log.Fatal("Error: Missing OAuth2 Code")
	}

	if config.Debug == true {
		fmt.Println("OAuth2 Code:", string(out))
	}

	AUTH_CODE := string(out)

	//************************************************************
	content := "client_id=" + secrets.Installed.ClientID
	content += "&client_secret=" + secrets.Installed.ClientSecret
	content += "&code=" + AUTH_CODE
	content += "&redirect_uri=http://localhost:9000"
	content += "&grant_type=authorization_code"
	//************************************************************

	endpoint := "https://www.googleapis.com/oauth2/v4/token"

	req := HttpRequest.NewRequest()

	req.SetHeaders(map[string]string{"Content-Type": "application/x-www-form-urlencoded"})

	res, err := req.Post(endpoint, content)

	if err != nil {
		fmt.Println("Error: ", err)
		return "", "", err
	}

	body, err := res.Body()

	if err != nil {
		fmt.Println("Error: ", err)
		return "", "", err
	}

	if config.Debug == true {
		fmt.Println("BODY:", string(body))
	}

	//************************************************************
	//
	//************************************************************

	var tokens OAuthTokens

	err = json.Unmarshal(body, &tokens)

	if err != nil {
		fmt.Println("Error: Cannot unmarshal JSON: ", err)
		return "", "", err
	}

	if config.Debug == true {
		fmt.Println("JSON:", tokens)
	}

	if tokens.Error != "" {
		fmt.Println("Error: Cannot authenticate")
		fmt.Println(tokens.Error)
		fmt.Println(tokens.ErrorDescription)
		return "", "", errors.New(tokens.ErrorDescription)
	}

	//************************************************************
	//
	//************************************************************

	var expires_at int64 = int64(time.Now().UTC().Unix()) + int64(tokens.ExpiresIn)

	var creds UserCredentials

	creds.ClientID = secrets.Installed.ClientID
	creds.ClientSecret = secrets.Installed.ClientSecret

	creds.RefreshToken = tokens.RefreshToken
	creds.Scope = tokens.Scope
	creds.Type = tokens.TokenType

	creds.AccessToken = tokens.AccessToken
	creds.IDToken = tokens.IDToken
	creds.ExpiresAt = expires_at

	//************************************************************
	//
	//************************************************************

	email, err := get_email_address(creds.AccessToken)

	if err == nil {
		fmt.Println("Email:", email)

		creds.Email = email
	}

	//************************************************************
	//
	//************************************************************

	err = saveUserCredentials(SavedUserCredentials, creds)

	if err != nil {
		fmt.Println("Error: Cannot save user credentials: ", err)
		return "", "", err
	}

	//************************************************************
	//
	//************************************************************

	// debug_displayAccessToken(creds.AccessToken)
	// debug_displayUserInfo(creds.AccessToken)
	// debug_displayIDToken(creds.AccessToken, creds.IDToken)

	return creds.AccessToken, creds.IDToken, nil
}
