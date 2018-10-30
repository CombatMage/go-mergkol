package main

import "os"

func fileExist(file string) bool {
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		return true
	}
	return false
}
