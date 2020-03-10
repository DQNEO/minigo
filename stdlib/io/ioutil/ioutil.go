package ioutil

import (
	"bytes"
	"io"
	"os"
)

const MYBUFSIZ = 65536 * 2

func readAll(f io.Reader, capacity int) ([]byte, error) {
	var buf *bytes.Buffer = &bytes.Buffer{}
	buf.Grow(capacity)
	_, err := buf.ReadFrom(f)
	return buf.Bytes(), err
}

func ReadFile(filename string) ([]byte, error) {
	var err error

	var f *os.File
	f, err = os.Open(filename)
	var n int = MYBUFSIZ

	var buf []byte
	buf, err = readAll(f, n)
	return buf, err
}
