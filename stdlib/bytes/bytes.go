package bytes

import "os"

type Buffer struct {
	bytes []byte
	buf   []byte
}

func (b *Buffer) Grow(capacity int) {
	var buf []byte
	buf = make([]byte, capacity, capacity)
	b.buf = buf
}

func (b *Buffer) ReadFrom(br *os.File) (int, error) {
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
