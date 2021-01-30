package hello

import "testing"

func TestHello(t *testing.T) {
	want := "Hello world!"
	got := Hello()
	if got != want {
		t.Fatalf("got: %v\nwant: %v\n", got, want)
	}
}
