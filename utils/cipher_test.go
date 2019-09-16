package utils

import (
	"testing"
)

func TestNumberToLetter(t *testing.T) {
	str := NumberToLetter("0123")
	if str != "abcd" {
		t.Fatal("[NumberToLetter] unexpected results")
	}
	t.Log("expect!")
}

func TestLetterToNumber(t *testing.T) {
	str := LetterToNumber("abcd")
	if str != "0123" {
		t.Fatal("[LetterToNumber] unexpected results")
	}
	t.Log("expect!")
}
