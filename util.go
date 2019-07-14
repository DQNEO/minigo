package main

import "strings"

func getBaseNameFromImport(path string) string {
	if strings.Contains(path, "/") {
		words := strings.Split(path, "/")
		r := words[len(words)-1]
		return r
	} else {
		return path
	}

}

func getIndex2(item string, list []string) int {
	for id, v := range list {
		if eq(bytes(v), bytes(item)) {
			return id
		}
	}
	return -1
}

func inArray2(item string, list []string) bool {
	for _, v := range list {
		if eq(bytes(v), bytes(item)) {
			return true
		}
	}
	return false
}
