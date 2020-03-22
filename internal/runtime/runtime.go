package runtime

import "unsafe"

var envp uintptr

var Envvars map[string]string

// Inital stack layout is illustrated in this page
// http://asm.sourceforge.net/articles/startup.html#st
func envvarsInit() {
	var p uintptr // **byte
	var envlines []string // []{"FOO=BAR\0", "HOME=/home/...\0", ..}

	for p = envp; true ; p = p + 8{
		var bpp **byte = (**byte)(unsafe.Pointer(p))
		if *bpp == nil {
			break
		}
		envlines = append(envlines, cstring2string(*bpp))
	}

	Envvars = make(map[string]string)

	for _, envline := range envlines {
		var i int
		var c byte
		for i, c = range []byte(envline) {
			if c == '=' {
				break
			}
		}
		key := envline[:i]
		value := envline[i+1:]
		Envvars[key] = value
	}

	// Debug envvars
	/*
	for k, v := range Envvars {
		printstring([]byte("# "))
		printstring([]byte(k))
		printstring([]byte(" = "))
		printstring([]byte(v))
		printstring([]byte("\n"))
	}
	 */
}

func runtime_getenv(key string) (string, bool) {
	v, ok := Envvars[key]
	return v, ok
}

func init() {
	heapInit()
	envvarsInit()
}

func printstring(b []byte) {
	write(2, b)
}

func panic(msg []byte) {
	printstring([]byte("panic: "))
	printstring(msg)
	printstring([]byte("\n"))
	exit(2)
}

const MiniGo int = 1

const stackSizeForThread = 1024*1024

const cloneFlag int = 331520

var task uintptr
var tasks []uintptr

func getTask() uintptr {
	tk := task
	tk0 := tasks[0]
	return tk0
}

func startm(mstart uintptr, tsk uintptr) {
	task = tsk
	tasks = append(tasks, task)
	newm(mstart)
}

func newm(mstart uintptr) {
	newm1(mstart)
}

func newm1(mstart uintptr) {
	stk := malloc(stackSizeForThread)
	newosproc(stk, mstart)
}

func newosproc(stk uintptr, mstart uintptr) int {
	return clone(cloneFlag, stk, mstart)
}
