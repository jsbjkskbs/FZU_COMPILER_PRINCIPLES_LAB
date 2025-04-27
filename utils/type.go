package utils

import (
	"strconv"
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

func HexToRune(hex string) rune {
	if len(hex) == 4 {
		i, _ := strconv.ParseInt(hex, 16, 32)
		return rune(i)
	} else if len(hex) == 8 {
		i, _ := strconv.ParseInt(hex, 16, 64)
		return rune(i)
	} else {
		return 0
	}
}

func OctalToRune(octal string) rune {
	if len(octal) == 2 {
		i, _ := strconv.ParseInt(octal, 8, 32)
		return rune(i)
	}
	return 0
}

func RemoveLeadingZeros(s string) string {
	for i, r := range s {
		if r != '0' {
			return s[i:]
		}
	}
	return "0"
}
