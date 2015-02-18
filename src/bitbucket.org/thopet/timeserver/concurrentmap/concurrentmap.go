package concurrentmap

import (
	"sync"
)

type CMap struct {
	values map[string]string
	lock   *sync.Mutex
}

// New creates a new CMap and returns a pointer.
func New() *CMap {
	return &CMap{
		values: make(map[string]string),
		lock:   new(sync.Mutex),
	}
}

// Get retrieves values given a key. Safe for concurrent use.
func (cm *CMap) Get(key string) (value string, ok bool) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	value, ok = cm.values[key]
	return value, ok
}

// Set sets an appropriate key-value in the backing map. If the key already
// exists it will be overridden.
func (cm *CMap) Set(key string, value string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	cm.values[key] = value
}

// Delete removes a key-value from the map. If the key doesn't exist,
// Delete is a no-op.
func (cm *CMap) Del(key string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	delete(cm.values, key)
}

// Creates a copy of the CMap and returns a pointer to the copy.
func (cm *CMap) Copy() *CMap {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	copy := New()
	for key, value := range cm.values {
		copy.values[key] = value
	}
	return copy
}

// TODO
func (cm *CMap) Equals(other *CMap) bool {
	for key, otherValue := range other.values {
		value, exists := cm.values[key]
		if !exists || value != otherValue  {
			return false
		}
	}
	for key, value := range cm.values {
		otherValue, exists := other.values[key]
		if !exists || value != otherValue  {
			return false
		}
	}
	return true
}
