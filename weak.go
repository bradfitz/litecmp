// Copyright 2020 Brad Fitzpatrick. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// From @ianlancetaylor's https://github.com/golang/go/issues/41303#issuecomment-717401656
// with a race fix.

package litecmp

import (
	"runtime"
	"sync/atomic"
)

// strong is a strong pointer to some object.
// To fetch the pointer, use the Get method.
type strong struct {
	p *intermediate
}

// Get returns the value to which s points.
func (s *strong) Get() interface{} {
	return s.p.val
}

// Weak returns a weak reference to the value to which strong points.
func (s *strong) Weak() *weak {
	return &weak{s.p}
}

// clear is a finalizer for s.
func (s *strong) clear() {
	var zeroEfacePtr *interface{}
	s.p.val.Store(zeroEfacePtr)
}

// Create a strong pointer to an object.
func makeStrong(val interface{}) *strong {
	inter := new(intermediate)
	inter.val.Store(&val)
	r := &strong{inter}
	runtime.SetFinalizer(r, func(s *strong) {
		println("weak-finalizer")
		s.clear()
	})
	return r
}

// Weak is a weak reference to some value.
type weak struct {
	p *intermediate
}

// Get returns the value to which w points.
// If there are no remaining strong pointers to the value,
// Get may return nil.
func (w *weak) Get() interface{} {
	ip, _ := w.p.val.Load().(*interface{})
	if ip != nil {
		return *ip
	}
	return nil
}

// intermediate is used to implement weak references.
type intermediate struct {
	val atomic.Value // of *interface{}
}
