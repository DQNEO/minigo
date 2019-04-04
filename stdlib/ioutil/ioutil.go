package ioutil

const MYBUFSIZ = 65536 * 2
const O_RDONLY = 0

func ReadFile(filename string) ([]byte, error) {
	var fd int
	var buf []byte
	buf = makeSlice(MYBUFSIZ, MYBUFSIZ)
	fd = open(filename, O_RDONLY)
	var nbtyes int
	nbtyes = read(fd, buf, MYBUFSIZ)
	var buf2 []byte
	buf2 = buf[0:nbtyes:nbtyes]
	// @TODO set len of buf
	return buf2,nil
}
