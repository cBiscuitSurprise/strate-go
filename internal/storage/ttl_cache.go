package storage

import (
	"sync"
	"time"

	"github.com/cBiscuitSurprise/strate-go/internal/game"
)

type gameRecord struct {
	value *game.Game
	time  int64
}

type TtlCache struct {
	cache map[string]*gameRecord
	lock  sync.Mutex
}

func NewTtlCache(initLen int, ttl int) (m *TtlCache) {
	m = &TtlCache{cache: make(map[string]*gameRecord, initLen)}
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

func (m *TtlCache) Len() int {
	return len(m.cache)
}

func (m *TtlCache) Put(id string, value *game.Game) {
	m.lock.Lock()
	now := time.Now().Unix()
	if record, ok := m.cache[id]; ok {
		record.time = time.Now().Unix()
	} else {
		record = &gameRecord{
			value: value,
			time:  now,
		}
		m.cache[id] = record
	}
	m.lock.Unlock()
}

func (m *TtlCache) Get(id string) (value *game.Game) {
	m.lock.Lock()
	if it, ok := m.cache[id]; ok {
		value = it.value
		it.time = time.Now().Unix()
	}
	m.lock.Unlock()
	return
}
