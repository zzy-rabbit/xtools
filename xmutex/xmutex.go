package xmutex

import (
	"context"
	"sync"
	"time"
)

type XMutex struct {
	mutex     sync.RWMutex
	readers   int
	writers   int
	onDeleted func()
}

func (m *XMutex) SetOnDeletedCallback(ctx context.Context, callback func()) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.onDeleted = callback
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
		if m.onDeleted != nil {
			m.onDeleted()
		}
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
		if m.onDeleted != nil {
			m.onDeleted()
		}
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
