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

import "sync"

type T[V any] struct {
	once sync.Once
	v    V
}

// Get calls fn if and only if it is being called for the first
// time and caches the result.
//
// Subsequent invocations of Get return the same result.
//
// It has the same semantics as sync.Once because it's a wrapper
// around sync.Once.
func (t *T[V]) Get(fn func() V) V {
	t.once.Do(func() { t.v = fn() })
	return t.v
}
