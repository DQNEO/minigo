package runtime

// func Caller(skip int) (pc uintptr, file string, line int, ok bool) {
func Caller(n int) (int, string, int ,bool) {
	return 0,"", 0, false
}

func FuncForPC(pc int) *Func {
	return nil
}

type Func struct {
	id int
}

func (f *Func) Name() string {
	return ""
}
