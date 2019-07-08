package universe

// builtin functions
// https://golang.org/ref/spec#Predeclared_identifiers

// Functions:
//	append cap close complex copy delete imag len
//	make new panic print println real recover

func make(x interface{}) interface{} {
}

func panic(s string) {
	printf("panic:%s\n", s)
	exit(1)
}

func print(x interface{}) {
	printf(x)
}

func recover() interface{} {
	return nil
}

type error interface {
	Error() string
}


