package bytes

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
	var nread int
	nread = read(fd, b.buf, cap(b.buf))
	bytes := b.buf[0:nread:nread]
	b.bytes = bytes
	return nread, nil
}

func (b *Buffer) Bytes() []byte {
	return b.bytes
}
