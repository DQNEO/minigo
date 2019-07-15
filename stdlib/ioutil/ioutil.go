package ioutil

import "os"

const MYBUFSIZ = 65536 * 2

func readAll(f *os.File, n int) ([]byte, error) {
	var buf []byte
	buf = makeSlice(n, n, 24)
	fid := f.innerFile.fd.Sysfd
	var nread int
	nread = read(fid, buf, n)
	var buf2 []byte
	buf2 = buf[0:nread:nread]
	return buf2,nil
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
