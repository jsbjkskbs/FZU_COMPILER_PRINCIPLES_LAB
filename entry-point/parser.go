package entrypoint

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	. "app/config"
	"app/lexer"
	"app/parser"
	. "app/utils"
	"app/utils/log"
	"app/utils/mmap"
)

var p *parser.Parser

func ParserTest() {
	files, err := GetDirFiles(Config.Path + "parser")
	if err != nil {
		panic(err)
	}
	fmt.Print(log.Sprintf(
		Divider(),
		log.Argument{Highlight: true, Format: "*** Parser Test ***\n", Args: []any{}},
		log.Argument{Highlight: true, Format: "*** Got ", Args: []any{}},
		log.Argument{FrontColor: log.Magenta, Highlight: true, Format: "%d ", Args: []any{len(files)}},
		log.Argument{Highlight: true, Format: "Files ***\n", Args: []any{}},
		Divider(),
	))

	for _, file := range files {
		fmt.Print(log.Sprintf(
			log.Argument{Highlight: true, Format: "<< ", Args: []any{}},
			log.Argument{FrontColor: log.Green, Highlight: true, Format: "%s\n", Args: []any{file.Path}},
		))
	}

	fmt.Print(log.Sprintf(
		Divider(),
	))

	// Create result directory if it doesn't exist
	err = os.MkdirAll(Config.Path+"parser/result", os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Print(log.Sprintf(
		log.Argument{FrontColor: log.Red, Highlight: true, Format: "!!! Starting tests... !!!\n", Args: []any{}},
		Divider(),
		log.Argument{FrontColor: log.Red, Highlight: true, Format: "!!! This may take a while to prepare the parser !!!\n", Args: []any{}},
	))

	st := time.Now()

	p = parser.NewParser()
	p.EnsureTable()

	fmt.Print(log.Sprintf(
		log.Argument{FrontColor: log.Green, Highlight: true, Format: "!!! Parser prepared, consume", Args: []any{}},
		log.Argument{FrontColor: log.Green, Highlight: true, Format: " %d ms", Args: []any{time.Since(st).Milliseconds()}},
		log.Argument{FrontColor: log.Green, Highlight: true, Format: "!!!\n", Args: []any{}},
	))

	wg := sync.WaitGroup{}
	wg.Add(len(files))
	for _, file := range files {
		go func(file FileInfo) {
			st := time.Now()
			defer wg.Done()
			result, err := os.Create(Config.Path + "parser/result/" + file.Info.Name() + ".result")
			if err != nil {
				panic(err)
			}
			defer func(result *os.File) {
				err := result.Close()
				if err != nil {
					panic(err)
				}
			}(result)
			writer := bufio.NewWriter(result)
			err = StartSingleParserTest(file.Path, writer)
			if err != nil {
				fmt.Println(
					log.Sprintf(log.Argument{FrontColor: log.Red, Highlight: true, Format: "!!! System Error: %s", Args: []any{err.Error()}}),
				)
			}

			err = writer.Flush()
			if err != nil {
				fmt.Println(
					log.Sprintf(log.Argument{FrontColor: log.Red, Highlight: true, Format: "!!! System Error: %s", Args: []any{err.Error()}}),
				)
			}

			fmt.Println(
				log.Sprintf(log.Argument{Highlight: true, Format: ">> Test for", Args: []any{}}),
				log.Sprintf(log.Argument{FrontColor: log.Green, Highlight: true, Format: "%s", Args: []any{file.Path}}),
				log.Sprintf(log.Argument{Highlight: true, Format: "finished, consume", Args: []any{}}),
				log.Sprintf(log.Argument{FrontColor: log.Green, Highlight: true, Format: "%d ms", Args: []any{time.Since(st).Milliseconds()}}),
			)
		}(file)
	}
	wg.Wait()

	fmt.Print(log.Sprintf(
		Divider(),
		log.Argument{FrontColor: log.Red, Highlight: true, Format: "!!! All tests finished !!!\n", Args: []any{}},
		Divider(),
	))
}

func StartSingleParserTest(filename string, writer io.Writer) error {
	file, err := mmap.NewMMapReader(filename)
	if err != nil {
		panic(err)
	}
	defer func(file *mmap.Reader) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	l := lexer.NewLexer(file)

	p.Parse(l, func(s string) {
		_, _ = fmt.Fprint(writer, s)
	})
	_, err = fmt.Fprintln(writer)
	if err != nil {
		return err
	}
	return nil
}
