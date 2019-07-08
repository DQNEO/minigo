package strconv

func Atoi(s string) (int , error) {
	gs := []byte(s)
	var i int
	i = atoi(gs)
	return i, nil
}
