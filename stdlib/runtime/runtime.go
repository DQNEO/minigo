package runtime

func Caller(n int) (interface{}, interface{},interface{},interface{}) {
	return nil,nil,nil,nil
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
