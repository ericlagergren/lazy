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

func TestTValue(t *testing.T) {
	var v T[int]
	got, ok := v.Value()
	if got != 0 || ok {
		t.Fatalf("expected (%d, %t), got (%d, %t)", 0, false, got, ok)
	}
	v.Get(func() int { return 42 })
	got, ok = v.Value()
	if got != 42 || !ok {
		t.Fatalf("expected (%d, %t), got (%d, %t)", 42, true, got, ok)
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

func TestEValue(t *testing.T) {
	var v E[int]
	got, ok := v.Value()
	if got != 0 || ok {
		t.Fatalf("expected (%d, %t), got (%d, %t)", 0, false, got, ok)
	}
	v.Get(func() (int, error) {
		return 42, nil
	})
	got, ok = v.Value()
	if got != 42 || !ok {
		t.Fatalf("expected (%d, %t), got (%d, %t)", 42, true, got, ok)
	}
}
