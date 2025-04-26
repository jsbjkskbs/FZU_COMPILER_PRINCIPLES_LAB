//go:build (linux || darwin) && !using_traditional_io && using_mmap_io

package mmap

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"syscall"
)

var _ = func() struct{} {
	println("mmap: using linux/unix implementation")
	return struct{}{}
}()

type _MemoryMap struct {
	data []byte
}

func (m *_MemoryMap) Close() error {
	if m.data == nil || (m.data != nil && len(m.data) == 0) {
		m.data = nil
		return nil
	}

	d := m.data
	m.data = nil
	runtime.SetFinalizer(m, nil)
	return syscall.Munmap(d)
}

func (m *_MemoryMap) ReadAt(p []byte, off int64) (int, error) {
	if m.data == nil {
		return 0, fmt.Errorf("mmap was closed")
	}
	if off < 0 || int64(len(m.data)) < off {
		return 0, fmt.Errorf("read at offset %d is invalid", off)
	}

	c := copy(p, m.data[off:])
	if c < len(p) {
		return c, io.EOF
	}
	return c, nil
}

func Open(filename string) (*_MemoryMap, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	filestat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	s := filestat.Size()
	if s == 0 {
		return &_MemoryMap{
			data: make([]byte, 0),
		}, nil
	} else if s < 0 {
		return nil, fmt.Errorf("file size is negative")
	}
	if s != int64(int(s)) {
		return nil, fmt.Errorf("file size is too large")
	}

	d, err := syscall.Mmap(int(f.Fd()), 0, int(s), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	m := &_MemoryMap{data: d}
	runtime.SetFinalizer(m, (*_MemoryMap).Close)
	return m, nil
}
