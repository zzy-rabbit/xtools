package xmutex

import (
	"context"
	"sync"
	"sync/atomic"
)

type XMutex struct {
	mu sync.RWMutex

	readers int32 // 当前读锁数量
	writer  int32 // 是否有写锁（0/1）
}

func (m *XMutex) RLock(ctx context.Context) {
	m.mu.RLock()
	atomic.AddInt32(&m.readers, 1)
}

func (m *XMutex) RUnlock(ctx context.Context) {
	atomic.AddInt32(&m.readers, -1)
	m.mu.RUnlock()
}

func (m *XMutex) Lock(ctx context.Context) {
	m.mu.Lock()
	atomic.StoreInt32(&m.writer, 1)
}

func (m *XMutex) Unlock(ctx context.Context) {
	atomic.StoreInt32(&m.writer, 0)
	m.mu.Unlock()
}

func (m *XMutex) Locked(ctx context.Context) bool {
	return atomic.LoadInt32(&m.writer) == 1 || atomic.LoadInt32(&m.readers) > 0
}

func (m *XMutex) RLocked(ctx context.Context) bool {
	return atomic.LoadInt32(&m.readers) > 0
}
