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
	if r.w+len(b) <= r.size {
		r.full = false
		n := copy(r.buffer[r.w:], b)
		r.w += n
		return n, nil
	}

	// TODO: Start wrapping around to the start
	if r.w+len(b) > r.size {
		// fill remainder of buffer before wrapping back to beginning
		remainder := r.size - r.w
		n = copy(r.buffer[r.w:r.w+remainder], b[:remainder])
		r.w += n
		if r.w >= r.size {
			r.full = true
			r.w = 0
		}

		_ = copy(r.buffer[r.w:len(b)-n], b[remainder:])
	}

	// if len(b)+r.w > r.size {
	// 	return 0, BufferOverflowErr
	// }

	return 0, nil
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
