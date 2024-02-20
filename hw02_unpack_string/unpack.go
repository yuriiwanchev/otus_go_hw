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
			err := addRepeated(&stringBuilder, input[i], input[i+1])
			if err != nil {
				return "", ErrInvalidString
			}

			i++
		}
	}

	return stringBuilder.String(), nil
}

func addRepeated(stringBuilder *strings.Builder, charToRepeat byte, countChar byte) error {
	count, err := strconv.Atoi(string(countChar))
	if err != nil {
		return ErrInvalidString
	}

	r := strings.Repeat(string(charToRepeat), count)
	stringBuilder.Write([]byte(r))
	return nil
}
