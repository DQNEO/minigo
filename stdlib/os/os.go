package os

import "unsafe"
import "syscall"

const O_RDONLY = 0

var sfd1 PollFD = PollFD{
	Sysfd: 1,
}
var sfd2 PollFD = PollFD{
	Sysfd: 2,
}

var sstdout file = file{
	fd: &sfd1,
}

var sstderr file = file{
	fd: &sfd2,
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
	Sysfd int
}

// File represents an open file descriptor.
type File struct {
	innerFile *file
}

func (f *File) Fd() int {
	return f.innerFile.fd.Sysfd
}

func openFileNolog(name string, flag int, perm int) (*File, error) {
	fid, err := syscall.Open(name, flag, perm)
	fl := &file{
		fd: &PollFD{
			Sysfd: fid,
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
	var n int
	var err error
	n, err = syscall.Write(fd.Sysfd, b)
	return n, err
}

func (fd *PollFD) Read(b []byte) (int, error) {
	var n int
	var err error
	n, err = syscall.Read(fd.Sysfd, b)
	return n, err
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

func (f *File) read(p []byte) (int, error) {
	fd := f.innerFile.fd
	var n int
	var err error
	n, err = fd.Read(p)
	return n, err
}

// Read reads up to len(b) bytes from the File.
func (f *File) Read(p []byte) (int, error) {
	//fd := f.innerFile.fd
	return f.read(p)
}

func (f *File) Readdirnames(n int) ([]string, error) {
	return f.readdirnames(n)
}

type linux_dirent struct {
	d_ino    int
	d_off    int
	d_reclen1 uint16
	d_type   byte
	d_name   byte
}

func cstring2string(b *byte) string {
	var bs []byte
	for {
		if b == nil || *b == 0 {
			break
		}
		bs = append(bs, *b)
		b = uintptr(b) + 1
	}
	return string(bs)
}


func (f *File) readdirnames(n int) ([]string, error) {
	fd := f.Fd()
	var r []string
	var counter int
	var buf [1024]byte
	for {
		nread, _ := syscall.ReadDirent(fd, buf[:])
		if nread == -1 {
			panic("getdents failed")
		}
		if nread == 0 {
			break
		}

		for bpos := 0; bpos < nread; 1 {
			var dirp *linux_dirent
			dirp = uintptr(buf) + uintptr(bpos)

			bpos = bpos + int(dirp.d_reclen1) // 24 is wrong
			var bp *byte = uintptr(unsafe.Pointer(&dirp.d_name))
			var s string = cstring2string(bp)
			r = append(r, s)
			counter++
		}
	}
	return r, nil
}

func Exit(code int) {
	syscall.Exit(code)
}

var Args []string

func init() {
	Args = runtime_args()
}

