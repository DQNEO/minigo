package main

import "testing"

func TestTokenize(t *testing.T) {
	filename := "token.txt"
	tokens := tokenizeFromFile(filename)
	expected := string(readFile(filename))

	var actual string
	for _, tok := range tokens {
		actual += tok.render()
	}

	if expected != actual {
		t.Errorf("%s expected but got %s", expected, actual)
	}
}
