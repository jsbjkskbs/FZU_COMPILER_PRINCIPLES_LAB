package utils

import (
	"unicode"
)

func IsDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func IsLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func IsHex(r rune) bool {
	return unicode.Is(unicode.Hex_Digit, r)
}

func IsOctal(r rune) bool {
	return r >= '0' && r <= '7'
}
