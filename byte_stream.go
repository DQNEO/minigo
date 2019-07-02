package main

import (
	"errors"
	"io/ioutil"
)

type ByteStream struct {
	filename  gostring
	source    []byte
	nextIndex int
	line      int
	column    int
}

func NewByteStreamFromString(name gostring, contents gostring) *ByteStream {
	return &ByteStream{
		filename:  name,
		source:    []byte(contents),
		nextIndex: 0,
		line:      1,
		column:    0,
	}
}

func NewByteStreamFromFile(path gostring) *ByteStream {
	s := readFile(path)
	return &ByteStream{
		filename:  path,
		source:    s,
		nextIndex: 0,
		line:      1,
		column:    0,
	}
}

func readFile(filename gostring) []byte {
	bytes, err := ioutil.ReadFile(string(filename))
	if err != nil {
		panic(err)
	}
	return bytes
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
