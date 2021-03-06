package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidString      = errors.New("invalid string")
	reverseSlash     rune = 92
)

func validateRunes(i int, r rune, prevRune rune, isPrevRuneEscaped bool) (bool, error) {
	if i == 0 && unicode.IsDigit(r) {
		return false, ErrInvalidString
	}

	if unicode.IsDigit(r) && unicode.IsDigit(prevRune) && !isPrevRuneEscaped {
		return false, ErrInvalidString
	}

	return true, nil
}

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}

	var result strings.Builder
	var prevRune rune
	var isPrevRuneEscaped bool

	for i, r := range str {
		_, validatingError := validateRunes(i, r, prevRune, isPrevRuneEscaped)

		if validatingError != nil {
			return "", validatingError
		}

		isLastIteration := len(str) == i+utf8.RuneLen(r)

		switch {
		case !isPrevRuneEscaped && prevRune == reverseSlash:
			isPrevRuneEscaped = true

			result.WriteRune(r)
		case unicode.IsDigit(r) && isPrevRuneEscaped:
			isPrevRuneEscaped = false
			digit, _ := strconv.Atoi(string(r))

			result.WriteString(strings.Repeat(string(prevRune), digit-1))
		case unicode.IsDigit(r) && unicode.IsLetter(prevRune):
			isPrevRuneEscaped = false
			digit, _ := strconv.Atoi(string(r))

			result.WriteString(strings.Repeat(string(prevRune), digit))
		case unicode.IsLetter(r) && unicode.IsLetter(prevRune) && isLastIteration:
			result.WriteRune(prevRune)
			result.WriteRune(r)
		case unicode.IsLetter(r) && isLastIteration:
			result.WriteRune(r)
		case unicode.IsLetter(prevRune):
			isPrevRuneEscaped = false

			result.WriteRune(prevRune)
		default:
			isPrevRuneEscaped = false
		}

		prevRune = r
	}

	return result.String(), nil
}
