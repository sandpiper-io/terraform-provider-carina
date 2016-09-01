// +build !windows

package carina

import (
	"os"
	"os/user"
)

func userHomeDir() (string, error) {
	if os.Getenv("HOME") != "" {
		return os.Getenv("HOME"), nil
	}
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return currentUser.HomeDir, nil
}
