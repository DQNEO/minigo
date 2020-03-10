package runtime

import "unsafe"

var _argv []*byte

// https://github.com/torvalds/linux/blob/v5.5/arch/x86/entry/syscalls/syscall_64.tbl
const __x64_sys_write = 1
const __x64_sys_brk  = 12
const __x64_sys_exit = 60

func init() {
	heapInit()
}

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

func runtime_args() []string {
	var r []string
	for _, a := range _argv {
		// Convert *byte to string
		var s string = cstring2string(a)
		r = append(r, s)
	}
	return r
}

func brk(addr uintptr) uintptr {
	var ret uintptr = Syscall(__x64_sys_brk, addr, 0, 0)
	return ret
}

// This is a copy from stconv
func itoa(i int) []byte {
	var r []byte
	var tmp []byte
	var isMinus bool

	// open(2) returs  0xffffffff 4294967295 on error.
	// I don't understand this yet.
	if i > 2147483648 {
		i = i - 2147483648*2
	}

	if i < 0 {
		i = i * -1
		isMinus = true
	}
	for i > 0 {
		mod := i % 10
		tmp = append(tmp, byte('0')+byte(mod))
		i = i / 10
	}

	if isMinus {
		r = append(r, '-')
	}

	for j := len(tmp) - 1; j >= 0; j-- {
		r = append(r, tmp[j])
	}

	if len(r) == 0 {
		return []byte{'0'}
	}
	return r
}

func printstring(b []byte) {
	var addr *byte = &b[0]
	write(2, addr, len(b))
}

func write(fd int, addr *byte, length int) {
	Syscall(__x64_sys_write, uintptr(fd), uintptr(unsafe.Pointer(addr)), uintptr(length))
}

func exit(code int) {
	Syscall(__x64_sys_exit, uintptr(code), 0 , 0)
}

func panic(msg []byte) {
	printstring([]byte("panic: "))
	printstring(msg)
	printstring([]byte("\n"))
	exit(2)
}

func reportMemoryUsage() {
	printstring([]byte("# memory-usage: "))
	i := getMemoryUsage()
	s := itoa(int(i))
	printstring([]byte(s))
	printstring([]byte("\n"))
}

func getMemoryUsage() uintptr {
	return heapPtr - heapHead
}

func eq(a string, b string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func cmpStrings(a string, b string, flag bool) bool {
	var ret bool
	ret = eq(a, b)
	if flag {
		return ret
	} else {
		return !ret
	}
}

func concat(as string, bs string) string {
	var r []byte
	for i := 0; i < len(as); i++ {
		r = append(r, as[i])
	}
	for i := 0; i < len(bs); i++ {
		r = append(r, bs[i])
	}
	return string(r)
}

const MiniGo int = 1
