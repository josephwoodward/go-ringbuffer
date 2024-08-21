package main

import "testing"

func TestTracksEnd(t *testing.T) {
	ring := NewRingBuffer(20)
	msg := []byte("Hello")

	ring.Write(msg)
	if ring.writePos() != len(msg) {
		t.Errorf("end expected to be %d but got %d", len(msg), ring.r)
	}

	msg = []byte("World")
	if ring.writePos() != len(msg) {
		t.Errorf("end expected to be %d but got %d", len(msg), ring.writePos())
	}
}

func TestReachesEnd(t *testing.T) {
	ring := NewRingBuffer(20)

	ring.Write([]byte("Hello"))
	ring.Write([]byte("World"))
	want := 10
	if ring.writePos() != want {
		t.Errorf("expected write position to be %d but got %d", want, ring.writePos())
	}

	b := make([]byte, 5)
	ring.Read(b)
	want = 5
	if ring.readPost() != want {
		t.Errorf("expected read position to be %d but was %d", want, ring.readPost())
	}
}
