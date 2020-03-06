// GOOS=linux GOARCH=amd64 go tool compile -N -S sum.go
package main

func sum(a int, b int) int {
	r := a + b
	return r
}

func main() {
	s := sum(2,3)
	println(s)
}
