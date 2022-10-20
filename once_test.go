// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lazy

import (
	"testing"
)

type one int

func (o *one) Increment() {
	*o++
}

func run(t *testing.T, v *T[int], o *one, c chan bool) {
	v.Get(func() int {
		o.Increment()
		return 0
	})
	if v := *o; v != 1 {
		t.Errorf("once failed inside run: %d is not 1", v)
	}
	c <- true
}

func TestTRace(t *testing.T) {
	o := new(one)
	var v T[int]
	c := make(chan bool)
	const N = 10
	for i := 0; i < N; i++ {
		go run(t, &v, o, c)
	}
	for i := 0; i < N; i++ {
		<-c
	}
	if *o != 1 {
		t.Errorf("once failed outside run: %d is not 1", *o)
	}
}

func TestTPanic(t *testing.T) {
	var v T[int]
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("T.Get did not panic")
			}
		}()
		v.Get(func() int {
			panic("failed")
		})
	}()

	v.Get(func() int {
		t.Fatalf("Once.Do called twice")
		return 0
	})
}

func BenchmarkOnce(b *testing.B) {
	var v T[int]
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			v.Get(func() int {
				return i
			})
		}
	})
}
