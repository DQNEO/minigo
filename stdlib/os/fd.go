package os

import "syscall"

var sfd1 PollFD = PollFD{
	Sysfd: 1,
}

var sfd2 PollFD = PollFD{
	Sysfd: 2,
}

type PollFD struct {
	Sysfd int
}

func openFileNolog(name string, flag int, perm int) (*File, error) {
	fid, err := syscall.Open(name, flag, perm)
	pfd := &PollFD{
		Sysfd: fid,
	}
	f := &File{
		pfd: pfd,
	}
	return f, err
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
