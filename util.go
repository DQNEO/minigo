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

func getIndex2(item gostring, list []gostring) int {
	for id, v := range list {
		if eqGostrings(v, item) {
			return id
		}
	}
	return -1
}

func inArray2(item gostring, list []gostring) bool {
	for _, v := range list {
		if eqGostrings(v, item) {
			return true
		}
	}
	return false
}
