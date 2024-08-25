package main

import (
	"errors"
	"testing"
)

func Test_ReadWritePointersTracked(t *testing.T) {
	t.Parallel()

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

func Test_Read(t *testing.T) {
	t.Run("can read into larger buffer", func(t *testing.T) {
		t.Parallel()

		ring := NewRingBuffer(10)
		ring.Write([]byte("0123456789"))

		b := make([]byte, 12)
		n, err := ring.Read(b)
		if err != nil {
			t.Errorf("unexpected error reading from buffer: %s", err)
		}

		if n != 10 {
			t.Errorf("number of bytes written should be 0 but was %d", n)
		}
		if string(b[:n]) != "0123456789" {
			t.Errorf("expected string from buffer but was %s", string(b))
		}
		if string(b[n:12]) != "\x00\x00" {
			t.Errorf("expected string from buffer but was %s", string(b))
		}
	})

	t.Run("can read segmenet into smaller buffer", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ring := NewRingBuffer(25)
		ring.Write([]byte("Hello"))
		ring.Write([]byte("World"))

		want := 10
		if ring.writePos() != want {
			t.Errorf("expected write position to be %d but got %d", want, ring.writePos())
		}

		// Act - Read "Hello"
		var n int
		var err error
		want = 5
		b := make([]byte, want)
		if n, err = ring.Read(b); err != nil || n != want {
			t.Errorf("unexpected error reading from buffer: %s", err)
		} else if n != want {
			t.Errorf("number of bytes written should be %d but was %d", want, n)
		}

		if string(b) != "Hello" {
			t.Error("expected initial input string to be read but was not")
		}
		if ring.readPos() != want {
			t.Errorf("expected read position to be %d but was %d", want, ring.readPos())
		}

		// Act - Read "World"
		if n, err = ring.Read(b); err != nil || n != want {
			t.Errorf("unexpected error reading from buffer: %s", err)
		} else if n != want {
			t.Errorf("number of bytes written should be %d but was %d", want, n)
		}

		if string(b) != "World" {
			t.Error("expected following string read to be read but was not")
		}
		if ring.readPos() != 10 {
			t.Errorf("expected read position to be %d but was %d", 10, ring.readPos())
		}
	})

	t.Run("can wrap around when reading", func(t *testing.T) {
		t.Parallel()
		// TODO
	})

	t.Run("can't read past write pointer - incrementally", func(t *testing.T) {
		t.Parallel()

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

func Test_ReadsToWritePosition(t *testing.T) {
}

func Test_CannotWriteMoreThanBufferSize(t *testing.T) {
	t.Run("cannot exceed in one go", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

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
