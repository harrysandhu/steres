package main

import "strings"

func NSplit(seq []byte, size int) []string {
	sep := " "
	s := string(seq)
	tks := strings.Split(s, sep)
	tokenizedSequence := []string{}
	for i := 0; i < len(tks); i++ {
		tmp := ""
		for c := 0; c < size; c++ {
			if i+c < len(tks) {
				tmp += tks[i+c] + " "
			} else {
				break
			}

		}
		tokenizedSequence = append(tokenizedSequence, string(tmp))
	}
	return tokenizedSequence
}
