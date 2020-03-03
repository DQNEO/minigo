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

func (fd *PollFD) ReadDirent(buf []byte) (int, error) {
	nread, _ := syscall.ReadDirent(fd.Sysfd, buf[:])
	return nread, nil
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

func ParseDirent(buf []byte, names []string) (int, []string) {
	p := uintptr(unsafe.Pointer(&buf[0]))
	var dirp *linux_dirent = p
	name := cstring2string(uintptr(unsafe.Pointer(&dirp.d_name)))
	names = append(names, name)
	return int(dirp.d_reclen1), names
}

func (f *File) readdirnames(n int) ([]string, error) {
	var names []string
	var buf [4096]byte
	for {
		nbuf, _ := f.innerFile.fd.ReadDirent(buf[:])
		if nbuf == -1 {
			panic("getdents failed")
		}
		if nbuf == 0 {
			break
		}
		var bufp int
		for ; bufp < nbuf; 1 {
			var reclen int
			reclen, names = ParseDirent(buf[bufp:], names)
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

