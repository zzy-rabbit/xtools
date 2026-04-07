package xmutex

import (
	"sync"
	"sync/atomic"
)

type XMutex struct {
	mu sync.RWMutex

	readers int32 // 当前读锁数量
	writer  int32 // 是否有写锁（0/1）
}

func (m *XMutex) RLock() {
	m.mu.RLock()
	atomic.AddInt32(&m.readers, 1)
}

func (m *XMutex) RUnlock() {
	atomic.AddInt32(&m.readers, -1)
	m.mu.RUnlock()
}

func (m *XMutex) Lock() {
	m.mu.Lock()
	atomic.StoreInt32(&m.writer, 1)
}

func (m *XMutex) Unlock() {
	atomic.StoreInt32(&m.writer, 0)
	m.mu.Unlock()
}

func (m *XMutex) Locked() bool {
	return atomic.LoadInt32(&m.writer) == 1 || atomic.LoadInt32(&m.readers) > 0
}

func (m *XMutex) RLocked() bool {
	return atomic.LoadInt32(&m.readers) > 0
}
