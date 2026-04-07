package xmutex

import (
	"context"
	"sync"
	"time"
)

type XMutex struct {
	mutex   sync.RWMutex
	readers int
	writers int
}

func (m *XMutex) Lock(ctx context.Context) {
	for {
		m.mutex.Lock()
		if m.readers > 0 || m.writers > 0 {
			m.mutex.Unlock()
			time.Sleep(time.Millisecond * 100)
			continue
		}
		m.writers++
		m.mutex.Unlock()
		return
	}
}

func (m *XMutex) Unlock(ctx context.Context) {
	for {
		m.mutex.Lock()
		if m.writers <= 0 {
			m.mutex.Unlock()
			return
		}
		m.writers--
		m.mutex.Unlock()
		return
	}
}

func (m *XMutex) RLock(ctx context.Context) {
	for {
		m.mutex.Lock()
		if m.writers > 0 {
			m.mutex.Unlock()
			time.Sleep(time.Millisecond * 100)
			continue
		}
		m.readers++
		m.mutex.Unlock()
		return
	}
}

func (m *XMutex) RUnlock(ctx context.Context) {
	for {
		m.mutex.Lock()
		if m.writers <= 0 {
			m.mutex.Unlock()
			return
		}
		m.readers--
		m.mutex.Unlock()
		return
	}
}

func (m *XMutex) RLocked(ctx context.Context) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.readers > 0
}

func (m *XMutex) Locked(ctx context.Context) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.writers > 0 || m.readers > 0
}
