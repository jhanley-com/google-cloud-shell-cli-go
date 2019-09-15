package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ConfigJson struct {
	ClientSecretsFile string `json:"oauth_json_file"`
	SSHFlags          string `json:"ssh_flags"`
	Debug             bool   `json:"debug"`
	Proxy             string `json:"proxy"`
	UrlFetch          string `json:"urlfetch"`
}

// Global Flags
type FlagsStruct struct {
	Adc   bool
	Auth  bool
	Login string
	Info  bool
}

type Config struct {
	// Global debug flag
	Debug bool

	UseAdcCredentials bool

	ProjectId string

	ClientSecretsFile string

	// Command to execute
	Command int

	// Command "exec"
	RemoteCommand string

	// Commands "download" and "upload"
	SrcFile string
	DstFile string

	// Command line global options
	Flags FlagsStruct

	// Command line ssh options
	sshFlags []string

	// ABS Path
	AbsPath string

	// proxy
	Proxy string

	// api proxy
	UrlFetch string
}

var config Config

func init_config() error {

	// path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	path, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}

	path = filepath.Dir(path)

	config.AbsPath = path

	// config.json
	in, err := os.Open(config.AbsPath + "/config.json")

	if err != nil {
		// fmt.Println(err)
	} else {

		defer in.Close()

		data, _ := ioutil.ReadAll(in)

		var configJson ConfigJson

		err = json.Unmarshal(data, &configJson)

		if err != nil {
			fmt.Println("Error: Cannot unmarshal JSON: ", err)
			return err
		}

		config.ClientSecretsFile = configJson.ClientSecretsFile

		if configJson.SSHFlags != "" && configJson.SSHFlags != "default" {
			config.sshFlags = []string{configJson.SSHFlags}
		}

		config.Debug = configJson.Debug

		if configJson.Proxy != "" {
			config.Proxy = configJson.Proxy
		}

		config.UrlFetch = configJson.UrlFetch

	}

	process_cmdline()

	return nil
}
