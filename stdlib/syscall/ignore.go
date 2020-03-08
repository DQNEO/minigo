// Declarations only. Minigo won't parse this file.
package syscall

func Syscall(number uintptr, a1 uintptr, a2 uintptr, a3 uintptr) int {
	return 0
}

func cstring2string(b *byte) string {
	return ""
}
