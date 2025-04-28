package entrypoint

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"sync"
	"time"

	. "app/config"
	"app/lexer"
	. "app/utils"
	"app/utils/log"
	"app/utils/mmap"
)

// LexerTest runs the lexer test on all files in the tests/lexer directory
func LexerTest() {
	files, err := GetDirFiles(Config.Path + "lexer")
	if err != nil {
		panic(err)
	}
	if len(Config.Files) > 0 {
		fs := []FileInfo{}
		for _, file := range files {
			if slices.Contains(Config.Files, file.Info.Name()) {
				fs = append(fs, file)
			}
		}
		files = fs
	}
	fmt.Print(log.Sprintf(
		Divider(),
		log.Argument{Highlight: true, Format: "***  Lexer Test  ***\n", Args: []any{}},
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

	if !Config.Silent {
		// Create result directory if it doesn't exist
		err = os.MkdirAll(Config.Path+"lexer/result", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	fmt.Print(log.Sprintf(
		log.Argument{FrontColor: log.Red, Highlight: true, Format: "!!! Starting tests... !!!\n", Args: []any{}},
		Divider(),
	))

	wg := sync.WaitGroup{}
	wg.Add(len(files))
	for _, file := range files {
		go func(file FileInfo) {
			st := time.Now()
			defer wg.Done()
			var f *os.File
			var err error
			if !Config.Silent {
				f, err = os.Create(Config.Path + "lexer/result/" + file.Info.Name() + ".result")
				if err != nil {
					panic(err)
				}
				defer func(f *os.File) {
					err := f.Close()
					if err != nil {
						panic(err)
					}
				}(f)
			}
			writer := bufio.NewWriter(f)
			err = StartSingleLexerTest(file.Path, writer)
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

// StartSingleLexerTest runs the lexer test on a single file
func StartSingleLexerTest(filename string, writer io.Writer) error {
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

	for {
		token, err := l.NextToken()
		if !Config.Silent && err != nil && !errors.Is(err, io.EOF) {
			_, err2 := fmt.Fprintf(writer, "Error: %s\n", err.Error())
			if err2 != nil {
				return err2
			}
		}
		if token.Type == lexer.EOF || errors.Is(err, io.EOF) {
			break
		}
		if !Config.Silent && token.Type != 0 {
			_, err = fmt.Fprintf(writer, "(%s, %s)\n", token.Type.ToString(), token.Val)
			if err != nil {
				return err
			}
		}
	}
	if !Config.Silent {
		_, err = fmt.Fprintln(writer)
		if err != nil {
			return err
		}
	}
	return nil
}
