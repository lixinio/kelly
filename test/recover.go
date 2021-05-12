package test

import (
	"errors"
	"testing"
)

func CheckError(t *testing.T, e error) {
	if r := recover(); r == nil {
		t.Errorf("checkError did not panic")
		return
	} else {
		switch x := r.(type) {
		case error:
			if !errors.Is(x, e) {
				t.Errorf("invalid error type %v", x)
			}
		default:
			t.Errorf("invalid error type %v", x)
		}
	}
}
