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

func Test_PerformsConsecutiveReads(t *testing.T) {
	ring := NewRingBuffer(20)

	ring.Write([]byte("Hello"))
	ring.Write([]byte("World"))
	want := 10
	if ring.writePos() != want {
		t.Errorf("expected write position to be %d but got %d", want, ring.writePos())
	}

	// Read "Hello"
	want = 5
	b := make([]byte, want)
	ring.Read(b)
	if string(b) != "Hello" {
		t.Error("expected initial input string to be read but was not")
	}
	if ring.readPost() != want {
		t.Errorf("expected read position to be %d but was %d", want, ring.readPost())
	}

	// Read "World"
	ring.Read(b)
	if string(b) != "World" {
		t.Error("expected following string read to be read but was not")
	}
	if ring.readPost() != 10 {
		t.Errorf("expected read position to be %d but was %d", 10, ring.readPost())
	}
	if ring.writePos() != 10 {
		t.Errorf("expected write position to be %d but was %d", 10, ring.readPost())
	}
}
