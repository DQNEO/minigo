package main

import (
	"./stdlib/strings"
)

// "foo/bar/buz" => "foo/bar"
func getDir(path string) string {
	found := strings.LastIndexByte(path, '/')
	if found == -1 {
		// not found
		return path
	}

	buf := []byte(path)
	return string(buf[0:found])
}

// "foo/bar/buz" => "buz"
func getBaseNameFromImport(path string) string {
	if path[len(path) - 1] == '/' {
		buf := []byte(path)
		path = string(buf[0:len(buf) - 1])
	}
	found := strings.LastIndexByte(path, '/')
	if found == -1 {
		// not found
		return path
	}
	buf := []byte(path)
	return string(buf[found+1:])
}

func getIndex(item string, list []string) int {
	for id, v := range list {
		if v == item {
			return id
		}
	}
	return -1
}

func inArray(item string, list []string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
