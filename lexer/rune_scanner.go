package lexer

import (
	"bufio"
	"fmt"
	"io"
	"unicode/utf8"
)

func NewRuneScannerBuffered(r io.Reader) *bufio.Reader {
	return bufio.NewReader(r)
}

type RuneScannerNoBuffered struct {
	reader io.Reader
	buffer [4]byte
	rune   rune
	size   int
	err    error
}

func (r *RuneScannerNoBuffered) ReadRune() (rune, int, error) {
	n, err := r.reader.Read(r.buffer[:])
	if err != nil {
		return 0, 0, err
	}
	r.rune, r.size = utf8.DecodeRune(r.buffer[:n])
	return r.rune, r.size, nil
}

func (r *RuneScannerNoBuffered) UnreadRune() error {
	if r.size == 0 {
		return fmt.Errorf("no rune to unread")
	}
	r.size = 0
	return nil
}

func NewRuneScannerNoBuffered(r io.Reader) *RuneScannerNoBuffered {
	return &RuneScannerNoBuffered{
		reader: r,
	}
}
