package utils

import "os/user"

func CurrentUser() string {
	u, err := user.Current()
	if err != nil {
		return "root"
	}
	return u.Username
}

func UserHome() string {
	u, err := user.Current()
	if err != nil {
		return "/root"
	}

	return u.HomeDir
}