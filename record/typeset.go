package record

import (
	"sync"
)

// TypeSet is a collection of Types,
// used to track types registered in the record system.
type TypeSet struct {
	sync.RWMutex
	types map[string]Type
}

// Type returns the type registered at given key
func (ts *TypeSet) Type(s string) Type {
	ts.RLock()
	defer ts.RUnlock()
	return ts.types[s]
}

// Types returns the types registered (allowed)
func (ts *TypeSet) Types() map[string]Type {
	ts.RLock()
	defer ts.RUnlock()

	out := map[string]Type{}
	for k, v := range ts.types {
		out[k] = v
	}
	return out
}
