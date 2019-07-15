package os

import "syscall"

const O_RDONLY = 0

var sstdout file = file{
	id :1,
}

var sstderr file = file{
	id :2,
}

var Stdout *File = &File{
	innerFile: &sstdout,
}

var Stderr *File = &File{
	innerFile: &sstderr,
}

type file struct {
	id int
}

// File represents an open file descriptor.
type File struct {
	innerFile *file
}

func openFileNolog(name string, flag int, perm int) (*File, error) {
	fd, err := syscall.Open(name, flag, perm)
	fl := &file{
		id:fd,
	}
	f := &File{
		innerFile: fl,
	}

	return f, err
}

func OpenFile(name string, flag int, perm int) (*File, error) {
	return openFileNolog(name, flag, perm)
}

func Open(name string) (*File, error) {
	return OpenFile(name, O_RDONLY, 0)
}

func (f *File) write(b []byte) (int, error) {
	var fid int = f.innerFile.id
	var n int
	var addr *byte = &b[0]
	n = write(fid, addr, len(b))
	return n,nil
}

// Write writes len(b) bytes to the File.
func (f *File) Write(b []byte) (int, error) {
	n, err := f.write(b)
	return n, err
}

func Exit(i int) {
}

var Args []string

func init() {
	Args = runtime_args()
}

