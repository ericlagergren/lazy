package lazy

import "testing"

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
