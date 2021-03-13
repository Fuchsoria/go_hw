package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}

	var result strings.Builder
	skipNext := false

	for i, r := range str {
		if skipNext {
			skipNext = false

			continue
		}

		var nextRune rune
		var inOneRune rune

		if len(str)-1 >= i+1 {
			nextRune, _ = utf8.DecodeRuneInString(string(str[i+1]))
		}

		if len(str)-1 >= i+2 {
			inOneRune, _ = utf8.DecodeRuneInString(string(str[i+2]))
		}

		if i == 0 && unicode.IsDigit(r) {
			return "", ErrInvalidString
		}

		if r != 92 && unicode.IsDigit(nextRune) && unicode.IsDigit(inOneRune) {
			return "", ErrInvalidString
		}

		if (!unicode.IsLetter(r) && r != 92 && !unicode.IsDigit(r)) && unicode.IsDigit(nextRune) {
			return "", ErrInvalidString
		}

		if unicode.IsLetter(nextRune) && r == 48 {
			continue
		} else if r == 92 && unicode.IsDigit(nextRune) && unicode.IsDigit(inOneRune) {
			digit, _ := strconv.Atoi(string(inOneRune))

			result.WriteString(strings.Repeat(string(nextRune), digit))

			skipNext = true
		} else if r == 92 && unicode.IsDigit(nextRune) {
			result.WriteRune(nextRune)

			skipNext = true
		} else if r == 92 && nextRune == 92 && unicode.IsDigit(inOneRune) {
			digit, _ := strconv.Atoi(string(inOneRune))

			result.WriteString(strings.Repeat(string(r), digit))

			skipNext = true
		} else if r == 92 && nextRune == 92 {
			result.WriteRune(92)

			skipNext = true
		} else if unicode.IsDigit(nextRune) {
			digit, _ := strconv.Atoi(string(nextRune))

			result.WriteString(strings.Repeat(string(r), digit))
		} else if unicode.IsLetter(r) {
			result.WriteRune(r)
		}
	}

	return result.String(), nil
}
