package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
)

func isWindows() bool {
	if runtime.GOOS == "windows" {
		return true;
	}

	return false
}

func check_os() error {
	if config.Debug == true {
		fmt.Println("Runtime Environment:")
		fmt.Printf("    OS: %s\n", runtime.GOOS)
		fmt.Printf("    ARCH: %s\n", runtime.GOARCH)
		fmt.Println("")
	}

	if runtime.GOOS == "windows" {
		// Supported
		return nil
	}

	if runtime.GOOS == "linux" {
		// err := errors.New("This operating system is not supported: " + runtime.GOOS)
		// return err
		// Supported
		return nil
	}

	if runtime.GOOS == "darwin" {
		err := errors.New("This operating system is not supported: " + runtime.GOOS)
		return err
	}

	err := errors.New("This operating system is not supported: " + runtime.GOOS)
	return err
}

func get_home_directory() (string, error) {
	path := ""

	if runtime.GOOS == "windows" {
		home_drive := os.Getenv("HOMEDRIVE")
		home_path := os.Getenv("HOMEPATH")

		if home_drive == "" {
			err := errors.New("Error: Missing Environment Variable: HOMEDRIVE")
			return path, err
		}

		if home_path == "" {
			err := errors.New("Error: Missing Environment Variable: HOMEPATH")
			return path, err
		}

		path = home_drive + home_path

		return path, nil
	} else if runtime.GOOS == "linux" {
		home_dir := os.Getenv("HOME")

		if home_dir == "" {
			err := errors.New("Error: Missing Environment Variable: HOME")
			return path, err
		}

		return home_dir, nil
	} else {
		err := errors.New("This operating system is not supported: " + runtime.GOOS)
		return path, err
	}
}
