package main

import "strings"

func getBaseNameFromImport(path bytes) bytes {
	if strings_Congtains(path, S("/")) {
		words := strings.Split(string(path), "/")
		r := words[len(words)-1]
		return bytes(r)
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
