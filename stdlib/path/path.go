package path

import "github.com/DQNEO/minigo/stdlib/strings"

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

	return path[:found]
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
		path = path[:len(path) - 1]
	}
	found := strings.LastIndexByte(path, '/')
	if found == -1 {
		// not found
		return path
	}

	return path[found+1:]
}
