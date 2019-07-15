package ioutil

import "os"
import "bytes"

const MYBUFSIZ = 65536 * 2


func readAll(f *os.File, capacity int) ([]byte, error) {
	var b *bytes.Buffer = &bytes.Buffer{
	}
	b.Grow(capacity)
	fd := f.Fd()
	_, err := b.ReadFrom(fd)
	return b.Bytes() ,err
}

func ReadFile(filenameAsString string) ([]byte, error) {
	var err error

	// Currently, there is no way to declare type of other package, so Let it infer
	var f *os.File
	f, err = os.Open(filenameAsString)
	var n int = MYBUFSIZ

	var buf []byte
	buf, err = readAll(f, n)
	return buf, err
}
