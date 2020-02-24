package syscall

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

func Open(path string, flag int, perm int) (int, error) {
	var fd int
	var _p0 *byte
	_p0 = BytePtrFromString(path)
	fd = syscallwrapper(__x64_sys_open, _p0, flag)
	return fd, nil
}

func Write(fd int, b []byte) (int, error) {
	var addr *byte = &b[0]
	var n int
	n = syscallwrapper(__x64_sys_write, fd, addr, len(b))
	return n, nil
}

func Read(fd int, b []byte) (int, error) {
	var ptr *byte
	ptr = &b[0]
	var nread int
	nread = syscallwrapper(__x64_sys_read, fd, ptr, cap(b))
	return nread, nil
}

func Exit(code int) {
	syscallwrapper(__x64_sys_exit, code)
	return
}
