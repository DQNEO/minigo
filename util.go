package main

import "strings"

func getBaseNameFromImport(path gostring) gostring {
	if strings.Contains(string(path), "/") {
		words := strings.Split(string(path), "/")
		r := words[len(words)-1]
		return gostring(r)
	} else {
		return path
	}

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
