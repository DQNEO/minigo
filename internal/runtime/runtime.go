package iruntime

var _argv []*byte

const heapSize uintptr = 640485760
var heap [heapSize]byte
var heapHead uintptr
var heapPtr uintptr
const intSize = 8

// https://github.com/torvalds/linux/blob/v5.5/arch/x86/entry/syscalls/syscall_64.tbl
const __x64_sys_write = 1
const __x64_sys_exit = 60

func init() {
	// set head address of heap
	heapHead = uintptr(&heap[0])
	heapPtr = heapHead
}

func cstring2string(b *byte) string {
	var bs []byte
	if b == nil {
		return string(bs)
	}

	var i int
	for {
		if b == nil || *b == 0 {
			break
		}
		bs = append(bs, *b)
		b++
	}
	return string(bs)
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

func malloc(size uintptr) uintptr {
	if heapPtr+ size > heapHead + heapSize {
		panic([]byte("malloc exceeds heap capacity"))
		return 0
	}
	r := heapPtr
	heapPtr += size
	return r
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
	for i>0 {
		mod := i % 10
		tmp = append(tmp, byte('0') + byte(mod))
		i = i /10
	}

	if isMinus {
		r = append(r, '-')
	}

	for j:=len(tmp)-1;j>=0;j--{
		r = append(r, tmp[j])
	}

	if len(r) == 0 {
		return []byte{'0'}
	}
	return r
}

func printstring(b []byte) {
	var addr *byte = &b[0]
	write(1, addr, len(b))
}

func write(fd int, addr *byte, length int) {
	syscall(__x64_sys_write, fd, addr, length)
}

func exit(code int) {
	syscall(__x64_sys_exit, code)
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

func append1(x []byte, elm byte) []byte {
	var z []byte
	xlen := len(x)
	zlen := xlen + 1

	if cap(x) >= zlen {
		z = x[:zlen]
	} else {
		var newcap int
		if xlen == 0 {
			newcap = 1
		} else {
			newcap = xlen * 2
		}
		z = makeSlice(zlen, newcap, 1)
		for i:=0;i<xlen;i++ {
			z[i] = x[i]
		}
	}

	z[xlen] = elm
	return z
}

func append8(x []int, elm int) []int {
	var z []int
	xlen := len(x)
	zlen := xlen + 1

	if cap(x) >= zlen {
		z = x[:zlen]
	} else {
		var newcap int
		if xlen == 0 {
			newcap = 1
		} else {
			newcap = xlen * 2
		}
 		z = makeSlice(zlen, newcap, 8)
		for i:=0;i<xlen;i++ {
			z[i] = x[i]
		}
	}

	z[xlen] = elm
	return z
}


func append24(x []interface{}, elm interface{}) []interface{} {
	//dumpInterface(elm)
	var z []interface{}
	xlen := len(x)
	zlen := xlen + 1

	if cap(x) >= zlen {
		z = x[:zlen]
	} else {
		var newcap int
		if xlen == 0 {
			newcap = 1
		} else {
			newcap = xlen * 2
		}
		z = makeSlice(zlen, newcap, 24)
		for i:=0;i<xlen;i++ {
			z[i] = x[i]
		}
	}

	z[xlen] = elm
	return z
}

func eq(a string, b string) bool {
	if len(a) != len(b) {
		return false
	}

	for i:=0;i<len(a);i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func cmpStrings(a string, b string, flag bool) bool {
	var ret bool
	ret = eq(a,b)
	if flag {
		return ret
	} else {
		return !ret
	}
}

func concat(as string, bs string) string {
	a := []byte(as)
	b := []byte(bs)

	var r []byte
	for i:=0;i<len(a);i++ {
		r = append(r, a[i])
	}
	for i:=0;i<len(b);i++ {
		r = append(r, b[i])
	}
	return string(r)
}

const MiniGo int = 1
