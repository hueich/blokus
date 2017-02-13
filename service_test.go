package blokusService

import (
	"testing"
)

func TestNewServiceNoPrefix(t *testing.T) {
	if s := New(""); s == nil {
		t.Error("New(''): got service==nil, want service")
	}
}

func TestNewServiceWithPrefix(t *testing.T) {
	if s := New("/foo"); s == nil {
		t.Error("New('/foo'): got service==nil, want service")
	}
}
