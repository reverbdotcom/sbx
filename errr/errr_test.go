package errr

import (
	"testing"
)

func TestNew(t *testing.T) {
	got := New("test")
	want := "🚫 test"
	if got.Error() != want {
		t.Errorf("got %q want %q", got, want)
	}
}
