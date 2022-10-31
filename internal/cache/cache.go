package cache

import (
	"errors"
	"sync"
)

const (
	Nil = CacheError("cache: nil")
)

type CacheError string

func (e CacheError) Error() string { return string(e) }

type Cache struct {
	mu   sync.RWMutex
	Data map[string]string
}

func (cch *Cache) Put(uid string, order string) {
	cch.mu.Lock()
	defer cch.mu.Unlock()
	cch.Data[uid] = order
}

func (cch *Cache) Get(uid string) (string, error) {
	cch.mu.RLock()
	defer cch.mu.RUnlock()
	var order string

	if order, found := cch.Data[uid]; found {
		return order, nil
	}

	return order, errors.New("not found")
}

func (cch *Cache) Del(uid string) {
	cch.mu.Lock()
	defer cch.mu.Unlock()
	delete(cch.Data, uid)
}
