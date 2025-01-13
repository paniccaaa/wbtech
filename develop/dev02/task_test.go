package dev02

import (
	"errors"
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		input string
		want  string
		err   error
	}{
		{"a4bc2d5e", "aaaabccddddde", nil},
		{"abcd", "abcd", nil},
		{"45", "", errors.New("некорректная строка")},
		{"", "", nil},
		{"qwe\\4\\5", "qwe45", nil},
		{"qwe\\45", "qwe44444", nil},
		{"qwe\\\\5", "qwe\\\\\\\\\\", nil},
		{"qwe10", "qweeeeeeeeee", nil},
		{"qwe\\1\\z2\\n2\\o10", "qwe1zznnoooooooooo", nil},
	}

	for _, test := range tests {
		got, err := Unpack(test.input)
		if test.err != nil && err == nil || test.err == nil && err != nil {
			t.Errorf("UnpackString(%q) returned error %v, want %v", test.input, err, test.err)
		}
		if got != test.want {
			t.Errorf("UnpackString(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}
