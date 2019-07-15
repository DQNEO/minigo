package ioutil

import "os"

const MYBUFSIZ = 65536 * 2

func ReadFile(filenameAsString string) ([]byte, error) {
	var fd int
	var buf []byte
	buf = makeSlice(MYBUFSIZ, MYBUFSIZ, 24)
	var err error

	// Currently, there is no way to declare type of other package, so Let it infer
	f := os.AnyFile
	f = nil

	f, err = os.Open(filenameAsString)

	var n int = MYBUFSIZ

	var nread int
	nread = read(f.id, buf, n)
	var buf2 []byte
	buf2 = buf[0:nread:nread]
	return buf2,nil
}
