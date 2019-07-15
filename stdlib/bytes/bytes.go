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

func (b *Buffer) ReadFrom(r *BufferFile) (int, error) {
	var nread int
	var err error
	nread, err = r.Read(b.buf)
	bytes := b.buf[0:nread:nread]
	b.bytes = bytes
	return nread, err
}

func (b *Buffer) Bytes() []byte {
	return b.bytes
}


type BufferFile struct {
	fd int
}

func (f *BufferFile) Fd() int {
	return f.fd
}

// Read implements io.Reader
func (f *BufferFile) Read(p []byte) (int, error) {
	//fd := f.innerFile.fd
	fd := f.Fd()
	var ptr *byte
	ptr = &p[0]
	var nread int
	nread = read(fd, ptr, cap(p))
	return nread, nil
}
