// Package lazy implements lazily initialized variables.
//
// To cut down on binary size and speed up init time, a fairly
// common pattern in Go is something similar to the following
//
//    var fooOnce sync.Once
//    var foo T
//
//    func loadFoo() foo {
//        fooOnce.Do(func() {
//            foo = ...
//        })
//        return foo
//    }
//
// (Almost every use of 'sync.Once' in the stdlib is to implement
// this pattern.)
//
// I'm lazy and we have generics, so why not ab^H^Huse them?
package lazy

import (
	"sync"
	"sync/atomic"
)

type T[V any] struct {
	done uint32     // set to 1 after Get is called
	m    sync.Mutex // guards v while writing
	v    V          // underlying value
}

// Get calls fn if and only if it is being called for the first
// time and caches the result.
//
// Subsequent invocations of Get return the same result.
//
// It has the same semantics as [sync.Once].
func (t *T[V]) Get(fn func() V) V {
	if atomic.LoadUint32(&t.done) == 0 {
		t.m.Lock()
		defer t.m.Unlock()

		if t.done == 0 {
			defer atomic.StoreUint32(&t.done, 1)
			t.v = fn()
		}
	}
	return t.v
}

// Value returns the cached value, if set.
func (t *T[V]) Value() (V, bool) {
	if atomic.LoadUint32(&t.done) != 0 {
		return t.v, true
	}

	t.m.Lock()
	defer t.m.Unlock()

	return t.v, t.done != 0
}

type E[V any] struct {
	done uint32     // set to 1 after Get is called
	m    sync.Mutex // guards v, err while writing
	v    V          // underlying value
	err  error      // underlying value
}

// Get calls fn if and only if it is being called for the first
// time and caches the result.
//
// Subsequent invocations of Get return the same result.
//
// It has the same semantics as [sync.Once].
func (t *E[V]) Get(fn func() (V, error)) (V, error) {
	if atomic.LoadUint32(&t.done) == 0 {
		t.m.Lock()
		defer t.m.Unlock()

		if t.done == 0 {
			defer atomic.StoreUint32(&t.done, 1)
			t.v, t.err = fn()
		}
	}
	return t.v, t.err
}

// Value returns the cached value, if set.
func (t *E[V]) Value() (V, bool) {
	if atomic.LoadUint32(&t.done) != 0 {
		return t.v, true
	}

	t.m.Lock()
	defer t.m.Unlock()

	return t.v, t.done != 0
}

// Must returns panics if err != nil and returns t otherwise.
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
