package config

import (
	"flag"
)

var Config = struct {
	Target string

	Lexer struct {
		UsingNoBufferedReader bool
	}

	Path   string
	Silent bool
}{}

func ReadFlag() {
	t := flag.String("t", "lexer", "Target to run: lexer or parser")
	lnb := flag.Bool("lexer--no-buffered", false, "Use no buffered reader for lexer")
	b := flag.Bool("b", false, "Enable benchmark mode")
	s := flag.Bool("s", false, "Stop writing results to file")
	flag.Parse()

	Config.Target = *t
	Config.Lexer.UsingNoBufferedReader = *lnb
	if *b {
		Config.Path = "tests/benchmark/"
		println("Benchmark mode enabled")
	} else {
		Config.Path = "tests/"
	}
	Config.Silent = *s
}
