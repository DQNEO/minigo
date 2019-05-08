package main

import "strings"

func getBaseNameFromImport(path string) string {
	var r string
	if strings.Contains(path, "/") {
		words := strings.Split(path, "/")
		r = words[len(words)-1]
	} else {
		r =  path
	}

	return r
}

func get_index(item string, list []string) int {
	for id, v := range list {
		if v == item {
			return id
		}
	}
	return -1
}

func in_array(item string, list []string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
