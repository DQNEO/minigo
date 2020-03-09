package os

import "syscall"

const O_RDONLY = 0

var Stdout *File = &File{
	pfd: &sfd1,
}

var Stderr *File = &File{
	pfd: &sfd2,
}

// File represents an open file descriptor.
type File struct {
	pfd *PollFD
}

func (f *File) Fd() int {
	return f.pfd.Sysfd
}

func OpenFile(name string, flag int, perm int) (*File, error) {
	return openFileNolog(name, flag, perm)
}

func Open(name string) (*File, error) {
	return OpenFile(name, O_RDONLY, 0)
}

func (f *File) write(b []byte) (int, error) {
	fd := f.pfd
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

func (f *File) read(p []byte) (int, error) {
	fd := f.pfd
	var n int
	var err error
	n, err = fd.Read(p)
	return n, err
}

// Read reads up to len(b) bytes from the File.
func (f *File) Read(p []byte) (int, error) {
	return f.read(p)
}

func (f *File) Readdirnames(n int) ([]string, error) {
	return f.readdirnames(n)
}

var void int

func (f *File) readdirnames(n int) ([]string, error) {
	var names []string
	var buf [4096]byte
	for {
		nbuf, _ := f.pfd.ReadDirent(buf[:])
		if nbuf == -1 {
			panic("getdents failed")
		}
		if nbuf == 0 {
			break
		}
		var bufp int
		for ; bufp < nbuf; {
			var reclen int
			reclen, void /* _ causes segmentation fault */, names = syscall.ParseDirent(buf[bufp:], 0,  names)
			bufp = bufp + reclen
		}
	}
	return names, nil
}

func Exit(code int) {
	syscall.Exit(code)
}

var Args []string

func init() {
	Args = runtime_args()
}

