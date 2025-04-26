//go:build (!linux && !windows && !darwin) || (using_traditional_io && !using_mmap_io)

package mmap

import (
	"fmt"
	"os"
)

var _ = func() struct{} {
	println("mmap: using default implementation[os.File]")
	return struct{}{}
}()

type _MemoryMap struct {
	f *os.File
}

func (m *_MemoryMap) Close() error {
	return m.f.Close()
}

func (m *_MemoryMap) ReadAt(p []byte, off int64) (int, error) {
	return m.f.ReadAt(p, off)
}

func Open(filename string) (*_MemoryMap, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	filestat, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}

	s := filestat.Size()
	if s == 0 {
		return &_MemoryMap{
			f: f,
		}, nil
	} else if s < 0 {
		f.Close()
		return nil, fmt.Errorf("file size is negative")
	}

	if s != int64(int(s)) {
		f.Close()
		return nil, fmt.Errorf("file size is too large")
	}
	return &_MemoryMap{f: f}, nil
}
