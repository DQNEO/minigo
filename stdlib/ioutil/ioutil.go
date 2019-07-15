package ioutil

import "os"

const MYBUFSIZ = 65536 * 2

type Buffer struct {
	bytes []byte
	buf []byte
}

func (b *Buffer) Grow(capacity int) {
	var buf []byte
	buf = makeSlice(capacity, capacity, 24)
	b.buf = buf
}

func (b *Buffer) ReadFrom(f *os.File) (int, error) {
	fd := f.innerFile.fd.Sysfd
	var nread int
	nread = read(fd, b.buf, cap(b.buf))
	bytes := b.buf[0:nread:nread]
	b.bytes = bytes
	return nread, nil
}

func (b *Buffer) Bytes() []byte {
	return b.bytes
}

func readAll(f *os.File, capacity int) ([]byte, error) {
	var b *Buffer = &Buffer{
	}
	b.Grow(capacity)
	_, err := b.ReadFrom(f)
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
