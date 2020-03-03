package syscall

import "unsafe"

func BytePtrFromString(s string) *byte {
	bs := []byte(s)
	var r *byte = &bs[0]
	return r
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
	_p0 = BytePtrFromString(path)
	fd = syscall(__x64_sys_open, _p0, flag)
	return fd, nil
}

func Write(fd int, b []byte) (int, error) {
	var addr *byte = &b[0]
	var n int
	n = syscall(__x64_sys_write, fd, addr, len(b))
	return n, nil
}

func Read(fd int, b []byte) (int, error) {
	var ptr *byte
	ptr = &b[0]
	var nread int
	nread = syscall(__x64_sys_read, fd, ptr, cap(b))
	return nread, nil
}

func ReadDirent(fd int, buf []byte) (int, error) {
	return Getdents(fd, buf)
}

func Getdents(fd int, buf []byte) (int, error) {
	var _p0 unsafe.Pointer
	_p0 = unsafe.Pointer(&buf[0])
	nread := syscall(_x64_getdents64, uintptr(fd), uintptr(_p0), len(buf))
	return nread, nil
}

type linux_dirent struct {
	d_ino    int
	d_off    int
	d_reclen1 uint16
	d_type   byte
	d_name   byte
}

//@TODO DRY: Same func exists in iruntime
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

func ParseDirent(buf []byte, void int, names []string) (int,  []string) {
	p := uintptr(unsafe.Pointer(&buf[0]))
	var dirp *linux_dirent = p
	name := cstring2string(uintptr(unsafe.Pointer(&dirp.d_name)))
	if name != "." && name != ".." {
		names = append(names, name)
	}
	return int(dirp.d_reclen1), names
}

func Exit(code int) {
	syscall(__x64_sys_exit, code)
	return
}
