package syscall

import "unsafe"

func BytePtrFromString(s string) (*byte, error) {
	bs := []byte(s)
	return &bs[0], nil
}

// 64-bit system call numbers
// https://github.com/torvalds/linux/blob/v5.5/arch/x86/entry/syscalls/syscall_64.tbl#L11
const __x64_sys_read = 0
const __x64_sys_write = 1
const __x64_sys_open = 2
const __x64_sys_exit = 60
const  _x64_getdents64 = 217

func Open(path string, flag int, perm int) (int, error) {
	var fd int
	var _p0 *byte
	_p0,_ = BytePtrFromString(path)
	fd = Syscall(__x64_sys_open, uintptr(unsafe.Pointer(_p0)), uintptr(flag), 0)
	return fd, nil
}

func Write(fd int, b []byte) (int, error) {
	var addr *byte = &b[0]
	var n int
	n = Syscall(__x64_sys_write, uintptr(fd), uintptr(unsafe.Pointer(addr)), uintptr(len(b)))
	return n, nil
}

func Read(fd int, b []byte) (int, error) {
	var ptr *byte
	ptr = &b[0]
	var nread int
	nread = Syscall(__x64_sys_read, uintptr(fd), uintptr(unsafe.Pointer(ptr)), uintptr(cap(b)))
	return nread, nil
}

func ReadDirent(fd int, buf []byte) (int, error) {
	return Getdents(fd, buf)
}

func Getdents(fd int, buf []byte) (int, error) {
	var _p0 unsafe.Pointer
	_p0 = unsafe.Pointer(&buf[0])
	nread := Syscall(_x64_getdents64, uintptr(fd), uintptr(_p0), uintptr(len(buf)))
	return nread, nil
}

// http://man7.org/linux/man-pages/man2/getdents64.2.html#top_of_page
type LinuxDirent64 struct {
	d_ino    int
	d_off    int
	d_reclen1 uint16
	d_type   byte
	d_name   byte
}

//@TODO DRY: Same func exists in iruntime
func cstring2string(b *byte) string {
	var buf []byte
	for {
		if b == nil || *b == 0 {
			break
		}
		buf = append(buf, *b)
		p := uintptr(unsafe.Pointer(b)) + 1
		b = (*byte)(unsafe.Pointer(p))
	}
	return string(buf)
}

func ParseDirent(buf []byte, void int, names []string) (int,  []string) {
	p := uintptr(unsafe.Pointer(&buf[0]))
	var dirp *LinuxDirent64 = p
	name := cstring2string(uintptr(unsafe.Pointer(&dirp.d_name)))
	if name != "." && name != ".." {
		names = append(names, name)
	}
	return int(dirp.d_reclen1), names
}

func Exit(code int) {
	Syscall(__x64_sys_exit, uintptr(code), 0 ,0)
	return
}
