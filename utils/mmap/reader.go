package mmap

import (
	"io"
)

type Reader struct {
	io.Reader
	io.Closer

	_r  *_MemoryMap
	pos int
}

func NewMMapReader(filepath string) (*Reader, error) {
	f, err := Open(filepath)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, io.EOF
	}

	return &Reader{
		_r: f,
	}, nil
}

func (m *Reader) Read(p []byte) (n int, err error) {
	if m._r == nil {
		return 0, io.EOF
	}

	n, err = m._r.ReadAt(p, int64(m.pos))
	if err != nil {
		return n, err
	}

	m.pos += n

	if n == 0 {
		return n, io.EOF
	}
	if n < len(p) {
		return n, io.EOF
	}

	return n, nil
}

func (m *Reader) Close() error {
	if m._r == nil {
		return nil
	}

	err := m._r.Close()
	if err != nil {
		return err
	}

	m._r = nil

	return nil
}
