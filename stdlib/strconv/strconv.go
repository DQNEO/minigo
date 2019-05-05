package strconv

func Atoi(s string) (int , error) {
	var i int
	i = atoi(s)
	return i, nil
}

