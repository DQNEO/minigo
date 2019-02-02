package main

//
// This file is to declare internal design of this project.
//
type Slice struct {
	pointer int
	len     int
	cap     int
}

type Map struct {
	pointer int
	len     int
	cap     int
}

type MapData struct {
	elements []Element
}

type Element struct {
	key   *interface{}
	value *interface{}
}
