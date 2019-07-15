package os

const O_RDONLY = 0

var Stdout *File = &File{
	id: 1,
}

var Stderr *File = &File{
	id: 2,
}

var AnyFile *File = &File{
	id:0,
}

// File represents an open file descriptor.
type File struct {
	id int
}

func OpenFile(name string, mode int, perm int) (*File, error) {
	var fd int
	var pchar *byte = name
	fd = open(pchar, mode)
	f := &File{
		id:fd,
	}

	return f, nil
}

func Open(name string) (*File, error) {
	return OpenFile(name, O_RDONLY, 0)
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

