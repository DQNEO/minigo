// Declarations only. Minigo won't parse this file.
package syscall

// Actual definition is in assembly code
func Syscall(number uintptr, a1 uintptr, a2 uintptr, a3 uintptr) int {
	return 0
}

// Actual definition is in iruntime code
func cstring2string(b *byte) string {
	return ""
}
