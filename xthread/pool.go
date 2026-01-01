package xthread

import (
	"context"
	"github.com/zzy-rabbit/xtools/xcontext"
	"github.com/zzy-rabbit/xtools/xerror"
)

type pool struct {
	poolCancel context.CancelFunc
	coreCancel context.CancelFunc
	coreChan   chan *core
	taskChan   chan task
}

type IPool interface {
	Do(ctx context.Context, f func()) (<-chan struct{}, xerror.IError)
	Close(ctx context.Context)
}

func New(_ context.Context, coreSize int, taskSize int) IPool {
	poolCtx, poolCancel := context.WithCancel(xcontext.Background())
	coreCtx, coreCancel := context.WithCancel(xcontext.Background())
	p := &pool{
		poolCancel: poolCancel,
		coreCancel: coreCancel,
		coreChan:   make(chan *core, coreSize),
		taskChan:   make(chan task, taskSize),
	}
	p.start(poolCtx)

	for i := 0; i < coreSize; i++ {
		p.coreChan <- newCore(coreCtx)
	}
	return p
}

func (p *pool) start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				func() {
					for {
						select {
						case t := <-p.taskChan:
							close(t.finish)
						default:
							return
						}
					}
				}()
				close(p.taskChan)
				close(p.coreChan)
				p.coreCancel()
				return
			case t := <-p.taskChan:
				func() {
					select {
					case <-ctx.Done():
						close(t.finish)
						return
					case c := <-p.coreChan:
						func() {
							defer func() {
								select {
								case <-ctx.Done():
									return
								default:
									p.coreChan <- c
								}
							}()
							c.do(ctx, t)
						}()
					}
				}()
			}
		}
	}()
}

func (p *pool) Do(ctx context.Context, f func()) (<-chan struct{}, xerror.IError) {
	finish := make(chan struct{})
	t := task{
		f:      f,
		finish: finish,
	}
	select {
	case <-ctx.Done():
		close(finish)
		return finish, xerror.ErrThreadPoolExited
	case p.taskChan <- t:
		return finish, nil
	default:
		close(finish)
		return finish, xerror.ErrTaskQueueFull
	}
}

func (p *pool) Close(ctx context.Context) {
	p.poolCancel()
}
