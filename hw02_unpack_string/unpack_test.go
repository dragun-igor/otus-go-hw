package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "\n3aб2в0", expected: "\n\n\naбб"},
		{input: "a\nб2в3", expected: "a\nббввв"},
		{input: "a\n0б2в3", expected: "aббввв"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackWithAsterisk(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestErrFirstChar(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "0p"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrFirstChar), "actual error %q", err)
		})
	}
}

func TestErrInvalidChar(t *testing.T) {
	invalidStrings := []string{"a-bc", "kf.y", "al/r"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidChar), "actual error %q", err)
		})
	}
}

func TestErrTwoDigit(t *testing.T) {
	invalidStrings := []string{"gj78", "i09", "y67"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrTwoDigit), "actual error %q", err)
		})
	}
}

func TestErrProtLetter(t *testing.T) {
	invalidStrings := []string{`lk\n`, `kj\q`, `uh\r`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrProtLetter), "actual error %q", err)
		})
	}
}

func TestErrLastBackslash(t *testing.T) {
	invalidStrings := []string{`lk\`, `kj\`, `uh\`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrLastBackslash), "actual error %q", err)
		})
	}
}
