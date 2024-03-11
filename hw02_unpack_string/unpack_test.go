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

		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},

		{input: "я9", expected: "яяяяяяяяя"},
		{input: "สวัสดี", expected: "สวัสดี"},
		{input: "สวัส4ดี", expected: "สวัสสสสดี"},
		{input: "🙃0", expected: ""},
		{input: "🙂9", expected: "🙂🙂🙂🙂🙂🙂🙂🙂🙂"},
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

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestUnpackAdditionaly(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
		expectedError  error
	}{
		{"", "", nil},
		{"abc", "abc", nil},
		{"123", "", ErrInvalidString},
		{"a2b3", "aabbb", nil},
		{"a4bc2d5e", "aaaabccddddde", nil},
		{"3abc", "", ErrInvalidString},
		{"aaa10b", "", ErrInvalidString},
		{"aaa0b", "aab", nil},
		{"d\n5abc", "d\n\n\n\n\nabc", nil},
		{"สวัสดี", "สวัสดี", nil},
	}

	for _, test := range tests {
		output, err := Unpack(test.input)
		if output != test.expectedOutput {
			t.Errorf("For input %q, expected output %q, but got %q", test.input, test.expectedOutput, output)
		}
		if !errors.Is(err, test.expectedError) {
			t.Errorf("For input %q, expected error %v, but got %v", test.input, test.expectedError, err)
		}
	}
}
