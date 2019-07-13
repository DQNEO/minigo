package os

var Stdout *File = &File{
	id: 1,
}

var Stderr *File = &File{
	id: 2,
}

// File represents an open file descriptor.
type File struct {
	id int
}

func (f *File) Write(b []byte) (int, error) {
	var fid int = f.id
	var n int
	var addr *byte = &b[0]
	n = write(fid, addr, len(b))
	return n,nil
}

func Exit(i int) {
}

var Args []string

func getArgs() []string {
	var r []string
	for _, a := range libcArgs {
		// we can regard *byte as string
		var s string = string(a)
		r = append(r, s)
	}
	return r
}

func init() {
	Args = getArgs()
}

