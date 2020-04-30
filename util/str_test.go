package util

import "testing"

func TestMakeRandomAlphaString(t *testing.T) {
	s := MakeRandomAlphaString(32)
	t.Logf("RandomAlphaString: %s\n", s)
}

func TestMakeRandomBase64String(t *testing.T) {
	s := MakeRandomBase64String(32)
	t.Logf("RandomBase64String: %s\n", s)
}
