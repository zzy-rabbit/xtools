package xmutex

import (
	"context"
	"sync"
	"sync/atomic"
)

type Mutex struct {
	mu sync.RWMutex

	readers int32 // 当前读锁数量
	writer  int32 // 是否有写锁（0/1）
}

func (m *Mutex) RLock(ctx context.Context) {
	m.mu.RLock()
	atomic.AddInt32(&m.readers, 1)
}

func (m *Mutex) RUnlock(ctx context.Context) {
	atomic.AddInt32(&m.readers, -1)
	m.mu.RUnlock()
}

func (m *Mutex) Lock(ctx context.Context) {
	m.mu.Lock()
	atomic.StoreInt32(&m.writer, 1)
}

func (m *Mutex) Unlock(ctx context.Context) {
	atomic.StoreInt32(&m.writer, 0)
	m.mu.Unlock()
}

func (m *Mutex) Locked(ctx context.Context) bool {
	return atomic.LoadInt32(&m.writer) == 1 || atomic.LoadInt32(&m.readers) > 0
}

func (m *Mutex) RLocked(ctx context.Context) bool {
	return atomic.LoadInt32(&m.readers) > 0
}
