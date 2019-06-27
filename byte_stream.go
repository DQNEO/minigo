package main

import (
	"errors"
	"fmt"
	"io/ioutil"
)

type ByteStream struct {
	filename  gostring
	source    []byte
	nextIndex int
	line      int
	column    int
}

func NewByteStreamFromString(name string, contents string) *ByteStream {
	return &ByteStream{
		filename:  gostring(name),
		source:    []byte(contents),
		nextIndex: 0,
		line:      1,
		column:    0,
	}
}

func NewByteStreamFromFile(path string) *ByteStream {
	s := readFile(path)
	return &ByteStream{
		filename:  gostring(path),
		source:    s,
		nextIndex: 0,
		line:      1,
		column:    0,
	}
}

func readFile(filename string) []byte {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (bs *ByteStream) location() string {
	return fmt.Sprintf("%s:%d:%d", bs.filename, bs.line, bs.column)
}

func (bs *ByteStream) get() (byte, error) {
	if bs.nextIndex >= len(bs.source) {
		return 0, errors.New("EOF")
	}
	r := bs.source[bs.nextIndex]
	if r == '\n' {
		bs.line++
		bs.column = 1
	}
	bs.nextIndex++
	bs.column++
	return r, nil
}

func (bs *ByteStream) unget() {
	bs.nextIndex--
	r := bs.source[bs.nextIndex]
	if r == '\n' {
		bs.line--
	}
}
