package utils

import "strings"

func NumberToLetter(numbers string) string {
	var strs []string
	for _, c := range numbers {
		if c < 48 || c > 57 {
			c = 48
		}
		strs = append(strs, string(c+49))
	}
	return strings.Join(strs, "")
}

func LetterToNumber(letters string) string {
	var strs []string
	for _, c := range letters {
		if c < 97 || c > 106 {
			c = 97
		}
		strs = append(strs, string(c-49))
	}
	return strings.Join(strs, "")
}
