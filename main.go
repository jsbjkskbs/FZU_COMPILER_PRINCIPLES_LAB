package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	. "app/config"
	"app/lexer"
	"app/parser"
	"app/utils/log"
	"app/utils/mmap"
)

func main() {
	// EnvChecker()

	ReadFlag()

	switch Config.Target {
	case "lexer":
		LexerTest()
	case "parser":
		ParserTest()
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

type FileInfo struct {
	info os.FileInfo
	path string
}

func GetDirFiles(dir string) ([]FileInfo, error) {
	var files []FileInfo
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		files = append(files, FileInfo{info: info, path: filepath.Join(dir, entry.Name())})
	}

	return files, nil
}

func Divider() log.Argument {
	return log.Argument{
		Highlight: true,
		Format:    "====================\n",
	}
}

// LexerTest runs the lexer test on all files in the tests/lexer directory
func LexerTest() {
	files, err := GetDirFiles(Config.Path + "lexer")
	if err != nil {
		panic(err)
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
			log.Argument{FrontColor: log.Green, Highlight: true, Format: "%s\n", Args: []any{file.path}},
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
				f, err = os.Create(Config.Path + "lexer/result/" + file.info.Name() + ".result")
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
			err = StartSingleLexerTest(file.path, writer)
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
				log.Sprintf(log.Argument{FrontColor: log.Green, Highlight: true, Format: "%s", Args: []any{file.path}}),
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
			log.Argument{FrontColor: log.Green, Highlight: true, Format: "%s\n", Args: []any{file.path}},
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
			result, err := os.Create(Config.Path + "parser/result/" + file.info.Name() + ".result")
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
			err = StartSingleParserTest(file.path, writer)
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
				log.Sprintf(log.Argument{FrontColor: log.Green, Highlight: true, Format: "%s", Args: []any{file.path}}),
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
