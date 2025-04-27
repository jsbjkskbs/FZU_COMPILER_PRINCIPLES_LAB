package main

import (
	"fmt"
	"runtime"

	. "app/config"
	entrypoint "app/entry-point"
	"app/utils/log"
)

func main() {
	// EnvChecker()

	ReadFlag()

	switch Config.Target {
	case "lexer":
		entrypoint.LexerTest()
	case "parser":
		entrypoint.ParserTest()
	default:
		println("Unknown mode:", Config.Target)
	}

}

func EnvChecker() {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" && runtime.GOOS != "unix" {
		fmt.Println(
			log.Sprintf(log.Argument{FrontColor: log.Yellow, Highlight: true, Format: "====================", Args: []any{}}),
			log.Sprintf(log.Argument{FrontColor: log.Yellow, Highlight: true, Format: "\n+++++++ Warn +++++++", Args: []any{}}),
			log.Sprintf(log.Argument{FrontColor: log.Yellow, Highlight: true, Format: "\nThis program may not", Args: []any{}}),
			log.Sprintf(log.Argument{FrontColor: log.Yellow, Highlight: true, Format: "\nperform correctly on", Args: []any{}}),
			log.Sprintf(log.Argument{FrontColor: log.Yellow, Highlight: true, Format: "\nnon-linux/unix OS.", Args: []any{}}),
			log.Sprintf(log.Argument{FrontColor: log.Yellow, Highlight: true, Format: "\n+++++++ Warn +++++++", Args: []any{}}),
			log.Sprintf(log.Argument{FrontColor: log.Yellow, Highlight: true, Format: "\n====================", Args: []any{}}),
		)
	}
}
