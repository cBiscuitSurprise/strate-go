package storage

import (
	"sync"
	"time"
)

type Record[V interface{}] struct {
	value *V
	time  int64
}

type TtlCache[K string, V interface{}] struct {
	cache map[K]*Record[V]
	lock  sync.Mutex
}

func NewTtlCache[K string, V interface{}](initLen int, ttl int) (m *TtlCache[K, V]) {
	m = &TtlCache[K, V]{cache: make(map[K]*Record[V], initLen)}
	go func() {
		for now := range time.Tick(time.Second) {
			m.lock.Lock()
			for k, v := range m.cache {
				if now.Unix()-v.time > int64(ttl) {
					delete(m.cache, k)
				}
			}
			m.lock.Unlock()
		}
	}()
	return
}

func (m *TtlCache[K, V]) Len() int {
	return len(m.cache)
}

func (m *TtlCache[K, V]) ForEach(callback func(key K, value *V)) {
	for k, r := range m.cache {
		callback(k, r.value)
	}
}

func (m *TtlCache[K, V]) Put(id K, value *V) {
	m.lock.Lock()
	now := time.Now().Unix()
	if record, ok := m.cache[id]; ok {
		record.time = time.Now().Unix()
	} else {
		record = &Record[V]{
			value: value,
			time:  now,
		}
		m.cache[id] = record
	}
	m.lock.Unlock()
}

func (m *TtlCache[K, V]) Get(id K) (value *V) {
	m.lock.Lock()
	if it, ok := m.cache[id]; ok {
		value = it.value
		it.time = time.Now().Unix()
	}
	m.lock.Unlock()
	return
}
