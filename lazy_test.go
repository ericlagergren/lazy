package lazy

import (
	"errors"
	"testing"
)

func TestT(t *testing.T) {
	var v T[int]
	for i := 0; i < 100; i++ {
		x := v.Get(func() int {
			return i + 42
		})
		if x != 42 {
			t.Fatalf("expected 42, got %d", x)
		}
	}
}

func TestE(t *testing.T) {
	var v E[int]
	err42 := errors.New("42!")
	for i := 0; i < 100; i++ {
		x, err := v.Get(func() (int, error) {
			return i + 42, err42
		})
		if x != 42 {
			t.Fatalf("expected 42, got %d", x)
		}
		if err != err42 {
			t.Fatalf("expected %v, got %v", err42, err)
		}
	}
}
