//go:build windows && !using_traditional_io && using_mmap_io

package mmap

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

var _ = func() struct{} {
	println("mmap: using windows implementation")
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
	return syscall.UnmapViewOfFile(uintptr(unsafe.Pointer(&d[0])))
}

func (m *_MemoryMap) ReadAt(p []byte, off int64) (int, error) {
	if m.data == nil {
		return 0, fmt.Errorf("mmap was closed")
	}
	if off < 0 || int64(len(m.data)) < off {
		return 0, fmt.Errorf("read at offset %d is invalid", off)
	}
	n := copy(p, m.data[off:])
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
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
		return nil, fmt.Errorf("file size is negative: %d", s)

	}

	if s != int64(int32(s)) {
		return nil, fmt.Errorf("file size is too large: %d", s)
	}

	l, h := uint32(s), uint32(s>>32)
	fileMapping, err := syscall.CreateFileMapping(syscall.Handle(f.Fd()), nil, syscall.PAGE_READONLY, h, l, nil)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(fileMapping)
	p, err := syscall.MapViewOfFile(fileMapping, syscall.FILE_MAP_READ, 0, 0, uintptr(s))
	if err != nil {
		return nil, err
	}
	data := unsafe.Slice((*byte)(unsafe.Pointer(p)), s)
	reader := &_MemoryMap{data: data}
	runtime.SetFinalizer(reader, func(r *_MemoryMap) {
		if r.data != nil {
			r.Close()
		}
	})
	return reader, nil
}
