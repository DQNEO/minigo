package main

func getBaseNameFromImport(path gostring) gostring {
	if strings_Congtains(path, S("/")) {
		words := strings_Split(path, S("/"))
		r := words[len(words)-1]
		return r
	} else {
		return path
	}

}

func getIndex(item gostring, list []gostring) int {
	for id, v := range list {
		if eqGostrings(v, item) {
			return id
		}
	}
	return -1
}

func inArray(item gostring, list []gostring) bool {
	for _, v := range list {
		if eqGostrings(v, item) {
			return true
		}
	}
	return false
}
