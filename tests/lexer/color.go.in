package log

import (
	"fmt"
)

type Color string

const (
	Black   Color = "Black"
	Red     Color = "Red"
	Green   Color = "Green"
	Yellow  Color = "Yellow"
	Blue    Color = "Blue"
	Magenta Color = "Magenta"
	Cyan    Color = "Cyan"
	White   Color = "White"
)

var frontColorToCode = map[Color]int{
	Black:   30,
	Red:     31,
	Green:   32,
	Yellow:  33,
	Blue:    34,
	Magenta: 35,
	Cyan:    36,
	White:   37,
}

var backColorToCode = map[Color]int{
	Black:   40,
	Red:     41,
	Green:   42,
	Yellow:  43,
	Blue:    44,
	Magenta: 45,
	Cyan:    46,
	White:   47,
}

type Argument struct {
	Format     string
	Args       []any
	FrontColor Color
	BackColor  Color
	Highlight  bool
	Underline  bool
}

func Sprintf(arguments ...Argument) string {
	s := ""
	for _, arg := range arguments {
		frontColorCode, ok := frontColorToCode[arg.FrontColor]
		if !ok {
			frontColorCode = 0
		}
		backColorCode, ok := backColorToCode[arg.BackColor]
		if !ok {
			backColorCode = 0
		}
		code := 0
		if arg.Highlight {
			code = 1
		} else if arg.Underline {
			code = 4
		}
		if frontColorCode == 0 && backColorCode == 0 {
			s += fmt.Sprintf("\033[%dm%s\033[0m", code, fmt.Sprintf(arg.Format, arg.Args...))
			continue
		} else if frontColorCode == 0 {
			s += fmt.Sprintf("\033[%d;%dm%s\033[0m", code, backColorCode, fmt.Sprintf(arg.Format, arg.Args...))
			continue
		} else if backColorCode == 0 {
			s += fmt.Sprintf("\033[%d;%dm%s\033[0m", code, frontColorCode, fmt.Sprintf(arg.Format, arg.Args...))
			continue
		}
		s += fmt.Sprintf("\033[%d;%d;%dm%s\033[0m", code, frontColorCode, backColorCode, fmt.Sprintf(arg.Format, arg.Args...))
	}
	return s
}
