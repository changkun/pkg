package gls

import (
	"sync"

	"github.com/changkun/gobase/rt"
)

var storage sync.Map // map[uint64]map[interface{}]interface{}

func init() {
	storage = sync.Map{}
}

// Store stores goroutine local data
func Store(k, v interface{}) {
	goid := rt.GoID()
	if _, ok := storage.Load(goid); !ok {
		storage.Store(goid, map[interface{}]interface{}{})
	}
	m, _ := storage.Load(goid)
	m.(map[interface{}]interface{})[k] = v
}

// Get gets goroutine local data
func Get(k interface{}) (interface{}, bool) {
	goid := rt.GoID()
	if _, ok := storage.Load(goid); !ok {
		return nil, false
	}
	m, _ := storage.Load(goid)
	v, ok := m.(map[interface{}]interface{})[k]
	return v, ok
}

// Clear deletes all gls data
func Clear() {
	storage.Delete(rt.GoID())
}
