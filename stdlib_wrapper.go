package main

import "io/ioutil"

func ioutil_ReadFile(filename gostring) ([]byte, error) {
	return ioutil.ReadFile(string(filename))
}
