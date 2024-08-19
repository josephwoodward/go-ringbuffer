package main

import "testing"

func TestTracksEnd(t *testing.T) {
	buf := NewRingBuffer(20)

	buf.Add([]byte("Hello"))
	want := 5
	if buf.end != want {
		t.Errorf("end expected to be %d but got %d", want, buf.end)
	}

	buf.Add([]byte("World"))
	want = 10
	if buf.end != want {
		t.Errorf("end expected to be %d but got %d", want, buf.end)
	}
}

func TestReachesEnd(t *testing.T) {
	buf := NewRingBuffer(20)

	buf.Add([]byte("Hello"))
	want := 5
	if buf.end != want {
		t.Errorf("end expected to be %d but got %d", want, buf.end)
	}

	buf.Add([]byte("World"))
	want = 10
	if buf.end != want {
		t.Errorf("end expected to be %d but got %d", want, buf.end)
	}
}
