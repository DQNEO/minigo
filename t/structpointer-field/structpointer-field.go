package main


type IrRoot struct {
	vars           []int
	funcs          []int
	stringLiterals []gostring
}

func f1() {

	root := &IrRoot{}

	var vars []int = []int{3}

	root.vars = vars

	fmtPrintf(S("%d\n"), len(vars))        // 1
	fmtPrintf(S("%d\n"), len(root.vars)+1) // 2
	fmtPrintf(S("%d\n"), root.vars[0])     // 3
}

type IrRoot2 struct {
	x  interface{}
	id int
}

func f2() {
	root := &IrRoot2{}

	var i int = 4
	var x interface{} = i
	root.x = x
	x = root.x
	var i2 int
	var ok bool
	i2, ok = x.(int)
	if !ok {
		fmtPrintf(S("ERROR\n"))
	}
	fmtPrintf(S("%d\n"), i2) // 4
}

func main() {
	f1()
	f2()
}
