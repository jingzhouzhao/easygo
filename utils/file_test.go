package utils

import (
	"testing"
)

func TestFileExists(t *testing.T) {
	if !FileExists("/") {
		t.Fatal("[FileExists] unexpected results")
	}
	if FileExists("/123") {
		t.Fatal("[FileExists] unexpected results")
	}
	t.Log("expect!")
}
