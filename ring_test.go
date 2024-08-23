package main

import (
	"errors"
	"testing"
)

func TestReadWritePointersTracked(t *testing.T) {
	ring := NewRingBuffer(20)
	msg := []byte("Hello")

	ring.Write(msg)
	if ring.writePos() != len(msg) {
		t.Errorf("end expected to be %d but got %d", len(msg), ring.readPos())
	}

	msg = []byte("World")
	if ring.writePos() != len(msg) {
		t.Errorf("end expected to be %d but got %d", len(msg), ring.writePos())
	}
}

func Test_ReadsToWritePosition(t *testing.T) {
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
	if ring.readPos() != want {
		t.Errorf("expected read position to be %d but was %d", want, ring.readPos())
	}

	// Read "World"
	ring.Read(b)
	if string(b) != "World" {
		t.Error("expected following string read to be read but was not")
	}
	if ring.readPos() != 10 {
		t.Errorf("expected read position to be %d but was %d", 10, ring.readPos())
	}
	if ring.writePos() != 10 {
		t.Errorf("expected write position to be %d but was %d", 10, ring.readPos())
	}
}
func Test_CannotReadPastWrite(t *testing.T) {
	t.Run("cannot read past write pointer - in one go", func(t *testing.T) {
		ring := NewRingBuffer(10)
		ring.Write([]byte("0123456789"))
		ring.Write([]byte("0123456789"))

		b := make([]byte, 12)
		n, err := ring.Read(b)
		if n != 0 {
			t.Errorf("number of bytes written should be 0 but was %d", n)
		}
		if !errors.Is(err, BufferOverflowErr) {
			t.Errorf("buffer flow error expected but was %s", err)
		}
	})

	t.Run("errors when reading past write pointer - incrementally", func(t *testing.T) {
		ring := NewRingBuffer(20)
		ring.Write([]byte("0123456789"))
		ring.Write([]byte("0123456789"))

		n, err := ring.Write([]byte("0"))
		if n != 0 {
			t.Errorf("number of bytes written should be 0 but was %d", n)
		}
		if !errors.Is(err, BufferOverflowErr) {
			t.Errorf("buffer flow error expected but was %s", err)
		}
	})
}

func Test_CannotWriteMoreThanBufferSize(t *testing.T) {
	t.Run("cannot exceed in one go", func(t *testing.T) {
		ring := NewRingBuffer(20)

		n, err := ring.Write([]byte("012345678901234567890"))
		if n != 0 {
			t.Errorf("number of bytes written should be 0 but was %d", n)
		}
		if !errors.Is(err, BufferOverflowErr) {
			t.Errorf("buffer flow error expected but was %s", err)
		}
	})

	t.Run("cannot exceed incrementally", func(t *testing.T) {
		ring := NewRingBuffer(20)
		ring.Write([]byte("0123456789"))
		ring.Write([]byte("0123456789"))

		n, err := ring.Write([]byte("0"))
		if n != 0 {
			t.Errorf("number of bytes written should be 0 but was %d", n)
		}
		if !errors.Is(err, BufferOverflowErr) {
			t.Errorf("buffer flow error expected but was %s", err)
		}
	})
}
