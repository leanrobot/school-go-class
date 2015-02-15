package concurrentmap

import (
	"sync"
)

type CMap struct {
	values map[string]string
	lock  *sync.Mutex
}

func New() *CMap {
	return &CMap{
		values: make(map[string]string),
		lock : new(sync.Mutex),
	}
}

func (cm *CMap) Get(key string) (value string, ok bool) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	value, ok = cm.values[key]
	return value, ok
}

func (cm *CMap) Set(key string, value string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	cm.values[key] = value
}

func (cm *CMap) Del(key string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	delete(cm.values, key)
}

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
	for key, _ := range other.values {
		if _, exists := cm.values[key]; !exists {
			return false
		}
	}
	for key, _ := range cm.values {
		if _, exists := other.values[key]; !exists {
			return false
		}
	}
	return true
}