package main

func plus(a int, b int) int {
	return a + b
}

func fvoid() {

}

func fvoidsemicolon() {

}

func fvoidreturn() {
	return
}

func f1() {
	fmtPrintf(S("%d\n"), plus(0, 1))
	fmtPrintf(S("%d\n"), plus(1, 1))
	fmtPrintf(S("%d\n"), plus(2, 1))
	fmtPrintf(S("%d\n"), plus(2, 2))
}

func f2() {
	var sum int
	sum = plus(2, 3)
	fmtPrintf(S("%d\n"), sum)
}

func returnNil() []int {
	return nil
}

var ary = [2]int{8, 9}

func returnSlice() []int {
	s := ary[:]
	return s
}

func receiveNilSlice() {
	nilSlice := returnNil()
	fmtPrintf(S("%d\n"), len(nilSlice)+6) // 6
}

func receiveIntSlice() {
	intSlice := returnSlice()
	fmtPrintf(S("%d\n"), len(intSlice)+5) // 7
	fmtPrintf(S("%d\n"), intSlice[0])     // 8
}

func returnSliceLiteral() []int {
	return []int{10}
}

func receiveSliceLiteral() {
	intSlice := returnSliceLiteral()
	fmtPrintf(S("%d\n"), len(intSlice)+8) // 9
	fmtPrintf(S("%d\n"), intSlice[0])     // 10
}

func returnStringSliceLiteral() []bytes {
	return []bytes{S("11")}
}

func receiveStringSliceLiteral() {
	slice := returnStringSliceLiteral()
	fmtPrintf(S("%s\n"), slice[0])
}

func main() {
	f1()
	f2()
	fvoid()
	fvoidsemicolon()
	fvoidreturn()
	receiveNilSlice()
	receiveIntSlice()
	receiveSliceLiteral()
	receiveStringSliceLiteral()
}
