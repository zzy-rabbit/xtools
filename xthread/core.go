package xthread

import "context"

type task struct {
	f      func()
	finish chan struct{}
}

type core struct {
	taskChan chan task
}

func newCore(ctx context.Context) *core {
	c := &core{
		taskChan: make(chan task),
	}
	c.start(ctx)
	return c
}

func (c *core) start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(c.taskChan)
				return
			case t := <-c.taskChan:
				func() {
					defer close(t.finish)
					t.f()
				}()
			}
		}
	}()
}

func (c *core) do(ctx context.Context, t task) {
	c.taskChan <- t
}
