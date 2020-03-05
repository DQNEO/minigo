package path

import "strings"

// "foo/bar/buz" => "foo/bar"
func Dir(path string) string {
	if len(path) == 0 {
		return "."
	}

	if path == "/" {
		return "/"
	}

	found := strings.LastIndexByte(path, '/')
	if found == -1 {
		// not found
		return path
	}

	buf := []byte(path)
	return string(buf[0:found])
}

// "foo/bar/buz" => "buz"
func Base(path string) string {
	if len(path) == 0 {
		return "."
	}

	if path == "/" {
		return "/"
	}
	if path[len(path) - 1] == '/' {
		buf := []byte(path)
		path = string(buf[0:len(buf) - 1])
	}
	found := strings.LastIndexByte(path, '/')
	if found == -1 {
		// not found
		return path
	}
	buf := []byte(path)
	return string(buf[found+1:])
}
