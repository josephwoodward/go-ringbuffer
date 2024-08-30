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

	// Write is still within buffer capacity
	if len(b)+r.w <= r.size {
		// TODO: We can improve this by always filling remaining
		n := copy(r.buffer[r.w:], b)
		r.w += n
		r.full = false
		return n, nil
	}

	// If write will exceed buffer then we need to wrap
	// Start by filling remaining capacity of buffer before wrapping back to beginning
	remaining := r.size - r.w
	n = copy(r.buffer[r.w:r.w+remaining], b[:remaining])
	r.w += n
	if r.w >= r.size {
		r.full = true
		r.w = 0
	}

	n += copy(r.buffer[r.w:len(b)-n], b[remaining:])
	return n, nil
}

func (r *RingBuffer) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}

	// ensure reader can never overtake writer
	if r.r > r.w {
		// TODO: Add option to block on read
		return 0, BufferOverflowErr
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
