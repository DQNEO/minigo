package main

import "strconv"

func strconv_Atoi(gs string) (int, error) {
	i, e := strconv.Atoi(gs)
	return i,e
}
