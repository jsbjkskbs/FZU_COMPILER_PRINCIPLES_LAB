package utils

import "strconv"

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

func AppendEscape(r rune) string {
	switch r {
	case 'n':
		return "\n"
	case 't':
		return "\t"
	case 'r':
		return "\r"
	case 'b':
		return "\b"
	case 'f':
		return "\f"
	case 'a':
		return "\a"
	case 'v':
		return "\v"
	case '"':
		return "\""
	case '\'':
		return "'"
	default:
		return string(r)
	}
}
