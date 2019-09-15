package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type ConfigJson struct {
	ClientSecretsFile string `json:"oauth_json_file"`
	SSHFlags          string `json:"ssh_flags"`
	Debug             bool   `json:"debug"`
	Proxy             string `json:"proxy"`
	UrlFetch          string `json:"urlfetch"`
	WinscpFlags       string `json:"winscp_flags"`
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

	// Command line winscp options
	winscpFlags []string

	// Path
	AbsPath     string
	PluginsPath string

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
	config.PluginsPath = path + "/plugins"

	// config.json
	in, err := os.Open(user_config_path("/config.json"))
	if err != nil {
		in, err = os.Open(config.AbsPath + "/config.json")
	}

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

		if configJson.WinscpFlags != "" {
			config.winscpFlags = append([]string{"/rawsettings"}, strings.Fields(configJson.WinscpFlags)...)
		}

	}

	process_cmdline()

	return nil
}

func user_config_path(filename string) string {

	// UserCredentials
	user, err := user.Current()
	if err == nil {
		configPath := user.HomeDir + "/.config/cloudshell/"

		_, err := os.Stat(configPath)
		if err != nil {
			err = os.Mkdir(configPath, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		return configPath + filename
	}
	return filename
}
