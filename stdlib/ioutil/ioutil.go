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
	var nbtyes int
	nbtyes = read(f.id, buf, MYBUFSIZ)
	var buf2 []byte
	buf2 = buf[0:nbtyes:nbtyes]
	// @TODO set len of buf
	return buf2,nil
}
