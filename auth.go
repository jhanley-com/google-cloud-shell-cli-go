package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/kirinlabs/HttpRequest"
	"golang.org/x/oauth2/google"
)

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
	// ClientID     string `json:"client_id"`
	// ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	// Scope        string `json:"scope"`
	// Type string `json:"type"`
	// The following two fields are option and exist after authentication
	AccessToken string `json:"access_token"`
	// IDToken     string `json:"id_token"`
	// Email       string `json:"email"`
	ExpiresAt int64 `json:"expires"`
}

type OAuthTokens struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	Scope            string `json:"scope"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
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
	var data []byte
	var err error

	if filename != "" {
		data, err = readCredentials(filename)

		if err != nil {
			fmt.Println("Error: Cannot read credentials JSON file \""+filename+"\"", err)
			os.Exit(1)
		}
	} else {
		//
		data, err = base64.StdEncoding.DecodeString("eyJpbnN0YWxsZWQiOnsiY2xpZW50X2lkIjoiNjU4OTM5NTQ0ODM3LXNtaHR1MG42N3A3MGdqM2o0Y2JtZGw2NmNma3RhcWx2LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwicHJvamVjdF9pZCI6InhjbG91ZHNoZWxsIiwiYXV0aF91cmkiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20vby9vYXV0aDIvYXV0aCIsInRva2VuX3VyaSI6Imh0dHBzOi8vb2F1dGgyLmdvb2dsZWFwaXMuY29tL3Rva2VuIiwiYXV0aF9wcm92aWRlcl94NTA5X2NlcnRfdXJsIjoiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vb2F1dGgyL3YxL2NlcnRzIiwiY2xpZW50X3NlY3JldCI6ImJEUEN4eTl3LTlac0ZoQ2hpX243Yk5ERyIsInJlZGlyZWN0X3VyaXMiOlsidXJuOmlldGY6d2c6b2F1dGg6Mi4wOm9vYiIsImh0dHA6Ly9sb2NhbGhvc3QiXX19")
	}

	// fmt.Println(string(data))

	err = json.Unmarshal(data, &secrets)

	if err != nil {
		fmt.Println("Error: Cannot unmarshal JSON: ", err)
		fmt.Println("File:", filename)
		return secrets, err
	}

	// config.ProjectId = secrets.Installed.ProjectID

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
	if config.Debug == true {
		fmt.Println("Save Credentials to:", filename)
	}

	j, err := json.MarshalIndent(creds, "", " ")

	if err != nil {
		fmt.Println("Error: Cannot marshall JSON:", err)
		return err
	}

	// err = ioutil.WriteFile(filename+".test", j, 0644)
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

func doRefresh(filename string) (string, bool) {
	endpoint := "https://www.googleapis.com/oauth2/v4/token"

	if config.UrlFetch != "" {
		endpoint = config.UrlFetch + endpoint
	}

	creds, err := loadUserCredentials(filename)

	if err != nil {
		fmt.Println(err)
		return "", false
	}

	// We want an access token that is good for a while.
	// Brand new tokens are valid for 3600 seconds
	// For testing require 15 minutes or 900 seconds

	var t time.Time = time.Unix(creds.ExpiresAt-(15*60), 0)
	// fmt.Println(t)

	// fmt.Println(time.Now())

	if time.Now().Before(t) {
		if config.Debug == true {
			fmt.Println("Saved credentials (Access Token) have not expired")
		}

		return creds.AccessToken, true
	}

	if config.Debug == true {
		fmt.Println("Must Refresh Token")
	}

	//************************************************************
	// Load the Google Client Secrets
	//************************************************************

	secrets, err := loadClientSecrets(config.ClientSecretsFile)

	if err != nil {
		fmt.Println(err)
		return "", false
	}

	//************************************************************
	// Build the authenticate URL
	//************************************************************

	content := "client_id=" + secrets.Installed.ClientID + "&"
	content += "client_secret=" + secrets.Installed.ClientSecret + "&"
	content += "grant_type=refresh_token&"
	content += "refresh_token=" + creds.RefreshToken

	req := HttpRequest.NewRequest()

	req.SetHeaders(map[string]string{"Content-Type": "application/x-www-form-urlencoded"})

	res, err := req.Post(endpoint, content)

	if err != nil {
		fmt.Println("Error: ", err)
		return "", false
	}

	body, err := res.Body()

	if err != nil {
		fmt.Println("Error: ", err)
		return "", false
	}

	var tokens OAuthTokens

	err = json.Unmarshal(body, &tokens)

	if err != nil {
		fmt.Println("Error: Cannot unmarshal JSON: ", err)
		return "", false
	}

	var expires int64 = int64(time.Now().UTC().Unix()) + int64(tokens.ExpiresIn)

	creds.AccessToken = tokens.AccessToken
	// creds.IDToken = tokens.IDToken
	creds.ExpiresAt = expires

	// email, err := get_email_address(tokens.AccessToken)

	// if err == nil {
	// 	if config.Debug == true {
	// 		fmt.Println("Email:", email)
	// 	}

	// 	creds.Email = email
	// }

	err = saveUserCredentials(filename, creds)

	if err != nil {
		fmt.Println("Error: Cannot save user credentials: ", err)
		return "", false
	}

	return creds.AccessToken, true
}

func debug_displayAccessToken(accessToken string) {
	endpoint := "https://www.googleapis.com/oauth2/v3/tokeninfo"

	if config.UrlFetch != "" {
		endpoint = config.UrlFetch + endpoint
	}

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

	fmt.Println("Token Info: ")
	fmt.Println(string(body))
}

func debug_displayUserInfo(accessToken string) {
	endpoint := "https://www.googleapis.com/oauth2/v3/userinfo"

	if config.UrlFetch != "" {
		endpoint = config.UrlFetch + endpoint
	}

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

// func debug_displayIDToken(accessToken, idToken string) {
// 	endpoint := "https://www.googleapis.com/oauth2/v3/tokeninfo"

// 	endpoint += "?id_token=" + idToken

// 	if config.UrlFetch != "" {
// 		endpoint = config.UrlFetch + endpoint
// 	}

// 	req := HttpRequest.NewRequest()

// 	req.SetHeaders(map[string]string{"Authorization": "Bearer " + accessToken})

// 	res, err := req.Get(endpoint)

// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 		return
// 	}

// 	body, err := res.Body()

// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 		return
// 	}

// 	fmt.Println(string(body))
// }

// func get_email_address(accessToken string) (string, error) {
// 	type Access_Token struct {
// 		Azp            string `json:"azp"`
// 		Aud            string `json:"aud"`
// 		Sub            string `json:"sub"`
// 		Scope          string `json:"scope"`
// 		Exp            string `json:"exp"`
// 		Expires_in     string `json:"expires_in"`
// 		Email          string `json:"email"`
// 		Email_verified string `json:"email_verified"`
// 		Access_type    string `json:"access_type"`
// 	}

// 	//************************************************************
// 	//
// 	//************************************************************

// 	endpoint := "https://www.googleapis.com/oauth2/v3/tokeninfo"

// 	if config.UrlFetch != "" {
// 		endpoint = config.UrlFetch + endpoint
// 	}

// 	req := HttpRequest.NewRequest()

// 	req.SetHeaders(map[string]string{"Authorization": "Bearer " + accessToken})

// 	//************************************************************
// 	//
// 	//************************************************************

// 	res, err := req.Get(endpoint)

// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 		return "", err
// 	}

// 	body, err := res.Body()

// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 		return "", err
// 	}

// 	//************************************************************
// 	//
// 	//************************************************************

// 	var tokens Access_Token

// 	err = json.Unmarshal(body, &tokens)

// 	if err != nil {
// 		fmt.Println("Error: Cannot unmarshal JSON: ", err)
// 		return "", err
// 	}

// 	return tokens.Email, nil
// }

func get_tokens() (string, error) {
	//************************************************************
	// Note: Application Default Credentials only work on Compute
	// Engine when interfacing with Cloud Shell.
	//************************************************************

	if config.Flags.Adc == true {
		return get_sa_tokens()
	}

	// fmt.Println("Auth:", config.Flags.Auth)
	// fmt.Println("Login:", config.Flags.Login)

	if config.Flags.Auth == false {
		if fileExists(SavedUserCredentials) {
			accessToken, valid := doRefresh(SavedUserCredentials)

			if valid == true {
				// fmt.Println("Access Token: ", accessToken)
				// fmt.Println("ID Token:     ", idToken)

				// debug_displayAccessToken(accessToken)
				// debug_displayUserInfo(accessToken)
				// debug_displayIDToken(accessToken, idToken)

				return accessToken, nil
			}
		}
	}

	//************************************************************
	// Load the Google Client Secrets
	//************************************************************

	secrets, err := loadClientSecrets(config.ClientSecretsFile)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	//************************************************************
	// Build the authenticate URL
	//************************************************************

	url := ENDPOINT
	url += "?client_id=" + secrets.Installed.ClientID
	url += "&response_type=code"
	url += "&scope=" + SCOPE
	url += "&access_type=offline"
	if len(config.Flags.Login) != 0 {
		url += "&login_hint=" + config.Flags.Login
	}

	url += "&redirect_uri=urn:ietf:wg:oauth:2.0:oob"

	token, err := manualAuthentication(secrets, url)

	// Wait for token to take effect
	// time.Sleep(1000 * time.Millisecond)

	return token, err
}

func get_sa_tokens() (string, error) {

	ctx := context.Background()

	creds, err := google.FindDefaultCredentials(ctx, SCOPE)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	token, err := creds.TokenSource.Token()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return token.AccessToken, nil
}

func manualAuthentication(secrets ClientSecrets, url string) (string, error) {

	fmt.Print("Go to the following link in your browser:\n\n" + url + "\n\nEnter verification code: ")

	reader := bufio.NewReader(os.Stdin)

	text, _ := reader.ReadString('\n')

	auth_code := strings.Replace(text, "\n", "", -1)

	return processAuthCode(secrets, auth_code)
}

func processAuthCode(secrets ClientSecrets, auth_code string) (string, error) {
	//************************************************************
	content := "client_id=" + secrets.Installed.ClientID
	content += "&client_secret=" + secrets.Installed.ClientSecret
	content += "&code=" + auth_code
	content += "&grant_type=authorization_code"
	content += "&redirect_uri=urn:ietf:wg:oauth:2.0:oob"
	//************************************************************

	endpoint := "https://www.googleapis.com/oauth2/v4/token"

	req := HttpRequest.NewRequest()

	req.SetHeaders(map[string]string{"Content-Type": "application/x-www-form-urlencoded"})

	res, err := req.Post(endpoint, content)

	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}

	body, err := res.Body()

	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
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
		return "", err
	}

	if config.Debug == true {
		fmt.Println("JSON:", tokens)
	}

	if tokens.Error != "" {
		fmt.Println("Error: Cannot authenticate")
		fmt.Println(tokens.Error)
		fmt.Println(tokens.ErrorDescription)
		return "", errors.New(tokens.ErrorDescription)
	}

	//************************************************************
	//
	//************************************************************

	var expires int64 = int64(time.Now().UTC().Unix()) + int64(tokens.ExpiresIn)

	var creds UserCredentials

	// creds.ClientID = secrets.Installed.ClientID
	// creds.ClientSecret = secrets.Installed.ClientSecret

	creds.RefreshToken = tokens.RefreshToken
	// creds.Scope = tokens.Scope
	// creds.Type = tokens.TokenType

	creds.AccessToken = tokens.AccessToken
	// creds.IDToken = tokens.IDToken
	creds.ExpiresAt = expires

	//************************************************************
	//
	//************************************************************

	// email, err := get_email_address(creds.AccessToken)

	// if err == nil {
	// 	fmt.Println("Email:", email)

	// 	creds.Email = email
	// }

	//************************************************************
	//
	//************************************************************

	err = saveUserCredentials(SavedUserCredentials, creds)

	if err != nil {
		fmt.Println("Error: Cannot save user credentials: ", err)
		return "", err
	}

	//************************************************************
	//
	//************************************************************

	if config.Debug == true {
		debug_displayAccessToken(creds.AccessToken)
		// debug_displayUserInfo(creds.AccessToken)
		// debug_displayIDToken(creds.AccessToken, creds.IDToken)
	}

	return creds.AccessToken, nil
}
