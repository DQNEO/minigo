package strconv

func Atoi(s string) (int , error) {
	var i int
	i = atoi(s)
	return i, nil
}

func Itoa(i int) string {
	var r []byte
	var tmp []byte
	var isMinus bool
	if i < 0 {
		i = i * -1
		isMinus = true
	}
	for i>0 {
		mod := i % 10
		tmp = append(tmp, byte('0') + byte(mod))
		i = i /10
	}

	if isMinus {
		r = append(r, '-')
	}

	for j:=len(tmp)-1;j>=0;j--{
		r = append(r, tmp[j])
	}

	if len(r) == 0 {
		r = []byte{'0'}
	}
	return string(r)
}
