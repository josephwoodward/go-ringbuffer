package main

type RingBuffer struct {
	start, end int
	size       int
	buffer     []byte
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buffer: make([]byte, size),
		start:  0,
		end:    0,
	}
}

func (b *RingBuffer) Add(s []byte) {
	for _, v := range s {
		b.buffer[b.end] = v
		b.end++
	}
}
