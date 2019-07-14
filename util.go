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

func getIndex(item bytes, list []bytes) int {
	for id, v := range list {
		if eq(v, item) {
			return id
		}
	}
	return -1
}

func inArray(item bytes, list []bytes) bool {
	for _, v := range list {
		if eq(v, item) {
			return true
		}
	}
	return false
}
