package main

import "os"

func fileExist(file string) bool {
	if _, err := os.Stat("/path/to/whatever"); !os.IsNotExist(err) {
		return false
	}
	return true
}
