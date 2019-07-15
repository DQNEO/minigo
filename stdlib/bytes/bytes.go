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

func (b *Buffer) ReadFrom(fd int) (int, error) {
	bytes_fd = fd
	var nread int
	var err error
	nread, err = Read(b.buf)
	bytes := b.buf[0:nread:nread]
	b.bytes = bytes
	return nread, err
}

func (b *Buffer) Bytes() []byte {
	return b.bytes
}

var bytes_fd int
// Read implements io.Reader
func Read(p []byte) (int, error) {
	//fd := f.innerFile.fd
	fd := bytes_fd
	var ptr *byte
	ptr = &p[0]
	var nread int
	nread = read(fd, ptr, cap(p))
	return nread, nil
}
