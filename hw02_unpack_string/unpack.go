package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func validateRunes(i int, r rune, nextRune rune, inOneRune rune) (bool, error) {
	if i == 0 && unicode.IsDigit(r) {
		return false, ErrInvalidString
	}

	if r != 92 && unicode.IsDigit(nextRune) && unicode.IsDigit(inOneRune) {
		return false, ErrInvalidString
	}

	if (!unicode.IsLetter(r) && r != 92 && !unicode.IsDigit(r)) && unicode.IsDigit(nextRune) {
		return false, ErrInvalidString
	}

	return true, nil
}

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

		var nextRune, inOneRune rune

		if len(str)-1 >= i+1 {
			nextRune, _ = utf8.DecodeRuneInString(string(str[i+1]))
		}

		if len(str)-1 >= i+2 {
			inOneRune, _ = utf8.DecodeRuneInString(string(str[i+2]))
		}

		_, validatingError := validateRunes(i, r, nextRune, inOneRune)

		if validatingError != nil {
			return "", validatingError
		}

		switch {
		case r == 92 && unicode.IsDigit(nextRune) && unicode.IsDigit(inOneRune):
			digit, _ := strconv.Atoi(string(inOneRune))

			result.WriteString(strings.Repeat(string(nextRune), digit))

			skipNext = true
		case r == 92 && unicode.IsDigit(nextRune):
			result.WriteRune(nextRune)

			skipNext = true
		case r == 92 && nextRune == 92 && unicode.IsDigit(inOneRune):
			digit, _ := strconv.Atoi(string(inOneRune))

			result.WriteString(strings.Repeat(string(r), digit))

			skipNext = true
		case r == 92 && nextRune == 92:
			result.WriteRune(92)

			skipNext = true
		case unicode.IsDigit(nextRune):
			digit, _ := strconv.Atoi(string(nextRune))

			result.WriteString(strings.Repeat(string(r), digit))
		case unicode.IsLetter(r):
			result.WriteRune(r)
		}
	}

	return result.String(), nil
}
