package server

import "testing"

func TestGenerateReply(t *testing.T) {
	var want string

	want = "Hello matt"
	if got := generateReply("matt"); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}

	want = "Hello world"
	if got := generateReply("world"); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
