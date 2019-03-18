package ioutil

const MYBUFSIZ = 1024
const O_RDONLY = 0

func ReadFile(filename string) ([]byte, error) {
	var fd int
	var buf []byte
	buf = makeSlice(MYBUFSIZ, MYBUFSIZ)
	fd = open(filename, O_RDONLY)
	read(fd, buf, MYBUFSIZ)
	// @TODO set len of buf
	return buf,nil
}
