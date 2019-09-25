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
	CMD_NONE = iota
	CMD_INFO
	CMD_PUTTY
	CMD_SSH
	CMD_INLINE_SSH
	CMD_WINSSH
	CMD_EXEC
	CMD_UPLOAD
	CMD_DOWNLOAD
	CMD_BITVISE
	CMD_WINSCP
	CMD_BENCHMARK_DOWNLOAD
	CMD_BENCHMARK_UPLOAD
)

func process_cmdline() {
	config.Command = CMD_NONE

	if len(os.Args) < 2 {
		cmd_help()
		os.Exit(0)
	}

	for _, arg := range os.Args {
		if arg == "-help" || arg == "--help" {
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

			if x == len(os.Args) - 1 {
				fmt.Println("Error: Missing email address to --login")
				os.Exit(1)
			}

			config.Flags.Login = os.Args[x + 1]
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

		// WINSCP args
		if strings.HasPrefix(arg, "/rawsettings") {
			// config.sshFlags = append(config.sshFlags, os.Args[x:]...)
			config.WinscpFlags = os.Args[x:]
			break
		}

		if strings.HasPrefix(arg, "-") {
			fmt.Println("Error: Unknown option: " + arg)
			os.Exit(1)
		}

		args = append(args, arg)
	}

	// fmt.Println("Debug:", config.Debug)
	// fmt.Println("Auth:", config.Flags.Auth)
	// fmt.Println("Login:", config.Flags.Login)

	if len(args) == 0 {
		cmd_help()
		os.Exit(0)
	}

	for x := 0; x < len(args); x++ {
		arg := args[x]

		switch arg {
		case "info":
			config.Command = CMD_INFO

		case "putty":
			if isWindows() == true {
				config.Command = CMD_PUTTY
			} else {
				fmt.Println("Error: This command is only supported on Windows. For Linux use ssh")
				os.Exit(1)
			}

		case "ssh":
			if isWindows() == true {
				config.Command = CMD_INLINE_SSH
			} else {
				config.Command = CMD_SSH
			}

		case "bitvise":
			if isWindows() == true {
				config.Command = CMD_BITVISE
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

		case "winssh":
			if isWindows() == true {
				config.Command = CMD_WINSSH
			} else {
				fmt.Println("Error: This command is only supported on Windows. For Linux use ssh")
				os.Exit(1)
			}

		case "exec":
			if len(args) < 2 {
				fmt.Println("Error: expected a remote command")
				os.Exit(1)
			}

			config.Command = CMD_EXEC
			config.RemoteCommand = args[x + 1]
			x++

		case "download":
			if len(args) < 2 {
				fmt.Println("Error: expected a source file name")
				os.Exit(1)
			}

			config.Command = CMD_DOWNLOAD
			config.SrcFile = strings.ReplaceAll(args[x + 1], "\\", "/")
			x++

			if len(args) >= 3 {
				config.DstFile = strings.ReplaceAll(args[x + 1], "\\", "/")
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

			path, err := filepath.Abs(args[x + 1])

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			config.Command = CMD_UPLOAD
			config.SrcFile = path
			x++

			if len(args) >= 3 {
				config.DstFile = strings.ReplaceAll(args[x + 1], "\\", "/")
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

		case "benchmark":
			if len(args) < 2 {
				fmt.Println("Error: expected download or upload option")
				os.Exit(1)
			}

			x++

			subcmd := args[x]

			switch subcmd {
			case "download":
				config.Command = CMD_BENCHMARK_DOWNLOAD
				config.benchmark_size = 128 * 1024
				config.benchmark_size = 10240 * 1024

			case "upload":
				config.Command = CMD_BENCHMARK_UPLOAD
				config.benchmark_size = 10240 * 1024

			default:
				fmt.Println("Error: expected a sub command (download, upload)")
				os.Exit(1)
			}

		default:
			if config.Command != CMD_NONE {
				fmt.Println("Error: Unknown command line argument: ", arg)
				os.Exit(1)
			}

			if isWindows() == true {
				fmt.Println("Error: expected a command (info, putty, ssh, winssh, exec, upload, download)")
			} else {
				fmt.Println("Error: expected a command (info, ssh, exec, upload, download)")
			}
			os.Exit(1)
		}
	}
}

func cmd_help() {
	fmt.Println("Usage: cloudshell [command]")
	fmt.Println("  cloudshell                            - display cloudshell program help")
	fmt.Println("  cloudshell info                       - display Cloud Shell information")
	if isWindows() == true {
		fmt.Println("  cloudshell putty                      - connect to Cloud Shell with Putty")
	}
	fmt.Println("  cloudshell ssh                        - connect to Cloud Shell with SSH")
	fmt.Println("  cloudshell winssh                     - connect to Cloud Shell with Windows OpenSSH")
	fmt.Println("  cloudshell winscp                     - connect to Cloud Shell with Windows WinSCP")
	fmt.Println("  cloudshell bitvise                    - connect to Cloud Shell with Windows Bitvise")
	fmt.Println("  cloudshell exec \"command\"             - Execute remote command on Cloud Shell")
	fmt.Println("  cloudshell upload src_file dst_file   - Upload local file to Cloud Shell")
	fmt.Println("  cloudshell download src_file dst_file - Download from Cloud Shell to local file")
	fmt.Println("  cloudshell benchmark download         - Benchmark download speed from Cloud Shell")
	fmt.Println("  cloudshell benchmark upload           - Benchmark upload speed from Cloud Shell")
	fmt.Println("")
	fmt.Println("--debug - Turn on debug output")
	fmt.Println("--adc  -  Use Application Default Credentials - Compute Engine only")
	fmt.Println("--auth  - (re)Authenticate ignoring user_credentials.json")
	fmt.Println("--login - Specify an email address as a login hint")
}
