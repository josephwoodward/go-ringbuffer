package main

import (
	"errors"
	"sync"
)

var (
	BufferOverflowErr = errors.New("buffer too small")
)

type RingBuffer struct {
	w, r   int
	size   int
	buffer []byte
	rw     sync.Mutex
	full   bool
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		size:   size,
		buffer: make([]byte, size),
		w:      0,
		r:      0,
	}
}

func (r *RingBuffer) Write(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	if len(b)+r.w > r.size {
		return 0, BufferOverflowErr
	}

	copy(r.buffer[r.w:], b)
	r.w += len(b)

	return 0, nil
}

func (r *RingBuffer) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}

	remainder := r.size - r.r
	copied := 0
	if len(b) > remainder {
		copied = copy(b, r.buffer[r.r:remainder])
		r.r += copied
	} else {
		copied = copy(b, r.buffer[r.r:r.r+len(b)])
		r.r += copied
	}

	return copied, nil
}

func (r *RingBuffer) writePos() int {
	return r.w
}

func (r *RingBuffer) readPos() int {
	return r.r
}
