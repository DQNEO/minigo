package main


func f1() {
	var  format string = "lea \\varname+\\offset(%%rip), %%rax"
	s := Sprintf(format)
	writeln([]byte(s))
}

func main() {
	f1()
}
