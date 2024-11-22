package mux

import (
	"sync"
)

type SafeMap struct {
	mu    sync.RWMutex
	items map[uint32]*Conn
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		items: make(map[uint32]*Conn),
	}
}

func (sm *SafeMap) Get(key uint32) (*Conn, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	val, exists := sm.items[key]
	return val, exists
}

func (sm *SafeMap) Set(key uint32, value *Conn) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.items[key] = value
}

func (sm *SafeMap) Delete(key uint32) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.items, key)
}
