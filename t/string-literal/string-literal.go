package main


func f1() {
	var  format gostring = S("lea \\varname+\\offset(%%rip), %%rax")
	s := Sprintf(format)
	writeln(s)
}

func main() {
	f1()
}
