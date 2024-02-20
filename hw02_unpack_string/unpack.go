package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var stringBuilder strings.Builder

	for i := 0; i < len(input); i++ {
		if unicode.IsDigit(rune(input[i])) {
			return "", ErrInvalidString
		}

		if input[i] == '\\' {
			i++
		}

		if i+1 >= len(input) {
			_, err := stringBuilder.Write([]byte(string(input[i])))
			if err != nil {
				return "", ErrInvalidString
			}
			break
		}

		isNextRuneIsNotDigit := !unicode.IsDigit(rune(input[i+1]))

		if isNextRuneIsNotDigit {
			_, err := stringBuilder.Write([]byte(string(input[i])))
			if err != nil {
				return "", ErrInvalidString
			}
		} else {
			count, err := strconv.Atoi(string(input[i+1]))
			if err != nil {
				return "", ErrInvalidString
			}

			r := strings.Repeat(string(input[i]), count)
			stringBuilder.Write([]byte(r))

			i++
		}
	}

	return stringBuilder.String(), nil
}
