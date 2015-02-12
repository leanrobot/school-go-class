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