package os

var _Stdout = File{
	id:1,
}
var _Stderr = File{
	id:2,
}

var Stderr *File = &_Stderr
var Stdout *File = &_Stdout


// File represents an open file descriptor.
type File struct {
	id int
}

func (f *File) Write(b []byte) (int, error) {
	var fid int = f.id
	var n int
	n = write(fid, string(b), len(b))
	return n,nil
}

func Exit(i int) {
}

func init() {

	// runtime_args is written in assembly code
	Args = runtime_args()
}

//func runtime_args() []string
