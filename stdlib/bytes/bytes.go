package bytes

import "io"

type Buffer struct {
	bytes []byte
	buf []byte
}

func (b *Buffer) Grow(capacity int) {
	var buf []byte
	buf = makeSlice(capacity, capacity, 24)
	b.buf = buf
}

func (b *Buffer) ReadFrom(br *File) (int, error) {
	var nread int
	var err error

	nread, err = br.Read(b.buf)
	bytes := b.buf[0:nread:nread]
	b.bytes = bytes
	return nread, err
}

func (b *Buffer) Bytes() []byte {
	return b.bytes
}


type File struct {
	fd int
}

func (f *File) Fd() int {
	return f.fd
}

// Read implements io.Reader
func (f *File) Read(p []byte) (int, error) {
	//fd := f.innerFile.fd
	fd := f.Fd()
	var ptr *byte
	ptr = &p[0]
	var nread int
	nread = read(fd, ptr, cap(p))
	return nread, nil
}
