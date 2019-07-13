package main

import (
	"errors"
)

type ByteStream struct {
	filename  bytes
	source    []byte
	nextIndex int
	line      int
	column    int
}

func NewByteStreamFromString(name bytes, contents bytes) *ByteStream {
	return &ByteStream{
		filename:  name,
		source:    []byte(contents),
		nextIndex: 0,
		line:      1,
		column:    0,
	}
}

func NewByteStreamFromFile(path bytes) *ByteStream {
	s := readFile(path)
	return &ByteStream{
		filename:  path,
		source:    s,
		nextIndex: 0,
		line:      1,
		column:    0,
	}
}

func readFile(filename bytes) []byte {
	bytes, err := ioutil_ReadFile(filename)
	if err != nil {
		panic(S("Unable to read file"))
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
