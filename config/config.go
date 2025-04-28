package config

import (
	"flag"
	"strings"
)

var Config = struct {
	Target string

	Lexer struct {
		UsingNoBufferedReader bool
	}

	Path   string
	Files  []string
	Silent bool
}{}

func ReadFlag() {
	t := flag.String("t", "lexer", "Target to run: lexer or parser")
	lnb := flag.Bool("lexer--no-buffered", false, "Use no buffered reader for lexer")
	b := flag.Bool("b", false, "Enable benchmark mode")
	s := flag.Bool("s", false, "Stop writing results to file")
	f := flag.String("f", "", "File to run tests on in the folder, split by |, eg. 1.in|2.in|3.in")
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
	if *f != "" {
		Config.Files = strings.Split(*f, "|")
	}
}
