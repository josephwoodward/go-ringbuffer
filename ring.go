package main

import "sync"

type RingBuffer struct {
	w, r   int
	size   int
	buffer []byte
	rw     sync.Mutex
	full   bool
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buffer: make([]byte, size),
		w:      0,
		r:      0,
	}
}

func (r *RingBuffer) Write(b []byte) (n int, er error) {
	if len(b) == 0 {
		return 0, nil
	}

	copy(r.buffer[r.w:], b)
	r.w += len(b)

	return 0, nil
}

func (r *RingBuffer) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}

	r.r += copy(b, r.buffer[r.r:r.r+len(b)])

	return r.r, nil
}

func (r *RingBuffer) writePos() int {
	return r.w
}

func (r *RingBuffer) readPost() int {
	return r.r
}
