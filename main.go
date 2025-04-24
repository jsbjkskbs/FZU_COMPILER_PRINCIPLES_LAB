package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"app/lexer"
	"app/utils/log"
	"app/utils/mmap"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: go run main.go <mode>")
		return
	}

	EnvChecker()

	mode := os.Args[1]
	switch mode {
	case "lexer":
		LexerTest()
	default:
		println("Unknown mode:", mode)
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
	files, err := GetDirFiles("tests/lexer")
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

	// Create result directory if it doesn't exist
	err = os.MkdirAll("tests/lexer/result", os.ModePerm)
	if err != nil {
		panic(err)
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
			result, err := os.Create("tests/lexer/result/" + file.info.Name() + ".result")
			if err != nil {
				panic(err)
			}
			defer func(result *os.File) {
				err := result.Close()
				if err != nil {
					panic(err)
				}
			}(result)
			err = StartSingleLexerTest(file.path, result)
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
		if err != nil && !errors.Is(err, io.EOF) {
			_, err2 := fmt.Fprintf(writer, "Error: %s\n", err.Error())
			if err2 != nil {
				return err2
			}
		}
		if token.Type == lexer.EOF || errors.Is(err, io.EOF) {
			break
		}
		if token.Type != 0 {
			_, err = fmt.Fprintf(writer, "(%s, %s)\n", token.Type.ToString(), token.Val)
			if err != nil {
				return err
			}
		}
	}
	_, err = fmt.Fprintln(writer)
	if err != nil {
		return err
	}
	return nil
}
