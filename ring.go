package main

type RingBuffer struct {
	w, r   int
	size   int
	buffer []byte
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buffer: make([]byte, size),
		w:      0,
		r:      0,
	}
}

func (b *RingBuffer) Write(s []byte) (n int, er error) {
	for _, v := range s {
		b.buffer[b.w] = v
		b.w++
	}

	return 0, nil
}

func (r *RingBuffer) Read(b []byte) (n int, er error) {
	r.r = copy(b, r.buffer)

	return r.r, nil
}

func (r *RingBuffer) writePos() int {
	return r.w
}

func (r *RingBuffer) readPost() int {
	return r.r
}
