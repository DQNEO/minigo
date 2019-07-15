package os

const O_RDONLY = 0

var Stdout *File = &File{
	id: 1,
}

var Stderr *File = &File{
	id: 2,
}

// File represents an open file descriptor.
type File struct {
	id int
}

func Open(name string) (int, error) {
	var fd int
	var filename *byte = name
	fd = open(filename, O_RDONLY)
	return fd, nil
}

func (f *File) Write(b []byte) (int, error) {
	var fid int = f.id
	var n int
	var addr *byte = &b[0]
	n = write(fid, addr, len(b))
	return n,nil
}

func Exit(i int) {
}

var Args []string

func init() {
	Args = runtime_args()
}

