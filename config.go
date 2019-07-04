package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ConfigJson struct {
	ClientSecretsFile	string   `json:"client_secrets_file"`
}

// Global Flags
type FlagsStruct struct {
	Adc		bool
	Auth		bool
	Login		string
	Info		bool
}

type Config struct {
	// Global debug flag
	Debug			bool

	UseAdcCredentials	bool

	ProjectId		string

	ClientSecretsFile	string

	// Command to execute
	Command			int

	// Command "exec"
	RemoteCommand		string

	// Commands "download" and "upload"
	SrcFile			string
	DstFile			string

	// Command line global options
	Flags			FlagsStruct
}

var config Config

func init_config() error {
	in, err := os.Open("config.json")

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer in.Close()

	data, err := ioutil.ReadAll(in)

	var configJson	ConfigJson

	err = json.Unmarshal(data, &configJson)

	if err != nil {
		fmt.Println("Error: Cannot unmarshal JSON: ", err)
		return err
	}

	config.ClientSecretsFile = configJson.ClientSecretsFile

	// fmt.Println("Client Secrets File:", config.ClientSecretsFile)

	process_cmdline()

	return nil
}
