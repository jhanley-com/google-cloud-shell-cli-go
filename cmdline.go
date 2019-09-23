package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Commands for this program
const (
	CMD_INFO = iota
	CMD_WINSSH
	CMD_SSH
	CMD_SSH_VSCODE
	CMD_EXEC
	CMD_UPLOAD
	CMD_DOWNLOAD
	CMD_PUTTY
	CMD_WINSCP
	CREATE_PUBKEY
)

func process_cmdline() {
	if len(os.Args) < 2 {
		cmd_help()
		os.Exit(0)
	}

	for _, arg := range os.Args {
		if arg == "--help" {
			cmd_help()
			os.Exit(0)
		}
	}

	var args []string

	for x := 1; x < len(os.Args); x++ {
		arg := os.Args[x]

		if arg == "-help" || arg == "--help" {
			cmd_help()
			os.Exit(0)
		}

		if arg == "-debug" || arg == "--debug" {
			config.Debug = true
			continue
		}

		if arg == "-adc" || arg == "--adc" {
			config.Flags.Adc = true
			continue
		}

		if arg == "-auth" || arg == "--auth" {
			config.Flags.Auth = true
			continue
		}

		if arg == "-login" || arg == "--login" {
			fmt.Println("index:", x)
			fmt.Println("count:", len(os.Args))

			if x == len(os.Args)-1 {
				fmt.Println("Error: Missing email address to --login")
				os.Exit(1)
			}

			config.Flags.Login = os.Args[x+1]
			config.Flags.Auth = true
			x++
			continue
		}

		if strings.HasPrefix(arg, "-login=") {
			p := arg[7:]
			config.Flags.Login = p
			config.Flags.Auth = true
			continue
		}

		if strings.HasPrefix(arg, "--login=") {
			p := arg[8:]
			config.Flags.Login = p
			config.Flags.Auth = true
			continue
		}

		// SSH args
		if arg == "-o" {
			config.sshFlags = append(config.sshFlags, "-o", os.Args[x+1])
			x++
			continue
		}
		if arg == "-D" {
			// support vs code
			if x == 1 && len(os.Args) == 5 && os.Args[4] == "bash" {

				if os.Args[3] == "cloudshell" {
					config.Debug = true

					if isWindows() == true {
						config.Command = CMD_WINSSH
					} else {
						config.Command = CMD_SSH
					}
				} else {
					config.Command = CMD_SSH_VSCODE
					config.sshFlags = append(config.sshFlags, "-D", os.Args[2], os.Args[3])
					break
				}
			}

			config.sshFlags = append(config.sshFlags, "-D", os.Args[x+1])
			x++
			continue
		}
		if arg == "-V" {
			fmt.Println("OpenSSH_for_Windows_7.7p1, LibreSSL 2.6.5")
			os.Exit(0)
		}
		// WINSCP args
		if strings.HasPrefix(arg, "/rawsettings") {
			// config.sshFlags = append(config.sshFlags, os.Args[x:]...)
			config.winscpFlags = os.Args[x:]
			break
		}
		// Proxy
		if arg == "-proxy" || arg == "--proxy" {
			config.Proxy = os.Args[x+1]
			x++
			continue
		}
		if arg == "-v2" || arg == "--v2ray" {
			config.Proxy = "v2ray"
			continue
		}
		if arg == "-ss" || arg == "--shadowsocks" {
			config.Proxy = "shadowsocks"
			continue
		}
		if arg == "-urlfetch" || arg == "--urlfetch" {
			config.UrlFetch = os.Args[x+1]
			x++
			continue
		}

		args = append(args, arg)
	}

	// fmt.Println("Debug:", config.Debug)
	// fmt.Println("Auth:", config.Flags.Auth)
	// fmt.Println("Login:", config.Flags.Login)

	for x := 0; x < len(args); x++ {
		arg := args[x]

		switch arg {
		case "info":
			config.Command = CMD_INFO

		case "ssh":
			if isWindows() == true {
				config.Command = CMD_WINSSH
			} else {
				config.Command = CMD_SSH
			}

			// argString := strings.Join(args, ",")
			// if strings.Contains(argString, "-D,") {
			// 	config.Flags.BindAddress = args[x + 2]

			// 	// break
			// 	return
			// }
			return

		case "exec":
			if len(args) < 2 {
				fmt.Println("Error: expected a remote command")
				os.Exit(1)
			}

			config.Command = CMD_EXEC
			config.RemoteCommand = args[x+1]
			x++

		case "download":
			if len(args) < 2 {
				fmt.Println("Error: expected a source file name")
				os.Exit(1)
			}

			config.Command = CMD_DOWNLOAD
			config.SrcFile = strings.ReplaceAll(args[x+1], "\\", "/")
			x++

			if len(args) >= 3 {
				config.DstFile = strings.ReplaceAll(args[x+1], "\\", "/")
				x++
			} else {
				_, file := path.Split(config.SrcFile)

				config.DstFile = file
			}

			if config.Debug == true {
				fmt.Println("SrcFile:", config.SrcFile)
				fmt.Println("DstFile:", config.DstFile)
			}

		case "upload":
			if len(args) < 2 {
				fmt.Println("Error: expected a source file name")
				os.Exit(1)
			}

			path, err := filepath.Abs(args[x+1])

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			config.Command = CMD_UPLOAD
			config.SrcFile = path
			x++

			if len(args) >= 3 {
				config.DstFile = strings.ReplaceAll(args[x+1], "\\", "/")
				x++
			} else {
				file := strings.ReplaceAll(config.SrcFile, "\\", "/")

				_, file = filepath.Split(file)

				config.DstFile = file
			}

			if config.Debug == true {
				fmt.Println("SrcFile:", config.SrcFile)
				fmt.Println("DstFile:", config.DstFile)
			}

		case "putty":
			if isWindows() == true {
				config.Command = CMD_PUTTY
			} else {
				fmt.Println("Error: This command is only supported on Windows. For Linux use ssh")
				os.Exit(1)
			}

		case "winscp":
			if isWindows() == true {
				config.Command = CMD_WINSCP
			} else {
				fmt.Println("Error: This command is only supported on Windows. For Linux use ssh")
				os.Exit(1)
			}

		case "push_pubkey":
			config.Command = CREATE_PUBKEY
			return

		default:
			if config.Command != 0 {
				return
			}
			if isWindows() == true {
				fmt.Println("Error: expected a command (info, putty, ssh, exec, upload, download)")
			} else {
				fmt.Println("Error: expected a command (info, ssh, exec, upload, download)")
			}
			os.Exit(1)
		}
	}
}

func cmd_help() {
	fmt.Println("Usage: cloudshell [command]")
	fmt.Println("  cloudshell info                       - display Cloud Shell information")
	if isWindows() == true {
		fmt.Println("  cloudshell putty                      - connect to Cloud Shell with Putty")
		fmt.Println("  cloudshell winscp                     - connect to Cloud Shell with WinSCP")
	}
	fmt.Println("  cloudshell ssh                        - connect to Cloud Shell with SSH")
	fmt.Println("  cloudshell exec \"command\"             - Execute remote command on Cloud Shell")
	fmt.Println("  cloudshell upload src_file dst_file   - Upload local file to Cloud Shell")
	fmt.Println("  cloudshell download src_file dst_file - Download from Cloud Shell to local file")
	fmt.Println("  cloudshell push_pubkey                - Push your public key to Cloud Shell")
	fmt.Println("")
	fmt.Println("--debug - Turn on debug output")
	fmt.Println("--adc  -  Use Application Default Credentials - Compute Engine only")
	fmt.Println("--auth  - (re)Authenticate ignoring user_credentials.json")
	fmt.Println("--login - Specify an email address as a login hint")
}
