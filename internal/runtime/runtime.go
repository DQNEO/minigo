package runtime

import "unsafe"

var envp uintptr

var Envvars map[string]string

// Inital stack layout is illustrated in this page
// http://asm.sourceforge.net/articles/startup.html#st
func envvarsInit() {
	var p uintptr // **byte
	var envlines []string
	p = envp
	for  {
		var bpp **byte = (**byte)(unsafe.Pointer(p))
		var bp *byte = *bpp
		if bp == nil {
			break
		}
		var s = cstring2string(bp)
		envlines = append(envlines, s)
		p = p + 8
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
