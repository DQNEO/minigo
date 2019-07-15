package os

import "syscall"

const O_RDONLY = 0

var sfd1 PollFD = PollFD{
	id:1,
}
var sfd2 PollFD = PollFD{
	id:2,
}

var sstdout file = file{
	fd : &sfd1,
}

var sstderr file = file{
	fd : &sfd2,
}

var Stdout *File = &File{
	innerFile: &sstdout,
}

var Stderr *File = &File{
	innerFile: &sstderr,
}

type file struct {
	fd *PollFD
}

type PollFD struct {
	id int
}

// File represents an open file descriptor.
type File struct {
	innerFile *file
}

func openFileNolog(name string, flag int, perm int) (*File, error) {
	fid, err := syscall.Open(name, flag, perm)
	fl := &file{
		fd:&PollFD{
			id: fid,
		},
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

func (fd *PollFD) Write(b []byte) (int, error) {
	var fid int = fd.id
	var n int
	var addr *byte = &b[0]
	n = write(fid, addr, len(b))
	return n,nil
}

func (f *File) write(b []byte) (int, error) {
	fd := f.innerFile.fd
	var n int
	var err error
	n, err = fd.Write(b)
	return n, err
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

