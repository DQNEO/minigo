// GOOS=linux GOARCH=amd64 go tool compile -N -S sum.go > sum.s
package main

func sum0(a int, b int) int {
	return a + b
}

func sum1(a int, b int) int {
	c := a + b
	return c
}

func assignstring(a string) string {
	b := a
	return b
}

func concatestring(a string, b string) string {
	return a + b
}

func main() {
	var i int
	i = sum0(2,3)
	println(i)

	i = sum1(2,3)
	println(i)

	assignstring("hello")

	s := concatestring("foo","bar")
	println(s)
}
