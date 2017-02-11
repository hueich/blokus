package blokusService

import (
	"testing"
)

func TestNewRouterNoPrefix(t *testing.T) {
	if r := NewRouter(""); r == nil {
		t.Error("NewRouter(''): got nil, want router")
	}
}

func TestNewRouterWithPrefix(t *testing.T) {
	if r := NewRouter("/foo"); r == nil {
		t.Error("NewRouter('/foo'): got nil, want router")
	}
}
