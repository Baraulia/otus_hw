package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var sb strings.Builder
	var currentSymbol string
	var shielding bool

	for _, s := range str {
		if shielding {
			sb.WriteString(string(s))
			currentSymbol = string(s)
			shielding = false
			continue
		}

		if string(s) == "\\" && !shielding {
			shielding = true
			continue
		}

		if unicode.IsDigit(s) {
			var digit int
			digit, err := strconv.Atoi(string(s))
			if err != nil || currentSymbol == "" {
				return "", ErrInvalidString
			}

			if digit == 0 {
				newStr := sb.String()[:len(sb.String())-1]
				sb.Reset()
				sb.WriteString(newStr)
				continue
			}

			sb.WriteString(strings.Repeat(currentSymbol, digit-1))
			currentSymbol = ""
			continue
		}
		sb.WriteString(string(s))
		currentSymbol = string(s)
	}

	return sb.String(), nil
}
