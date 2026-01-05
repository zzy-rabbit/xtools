package xthread

import (
	"github.com/zzy-rabbit/xtools/xcontext"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	ctx := xcontext.Background()
	p := New(ctx, 10, 100)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	waitGroup := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		f, err := p.Do(ctx, func() {
			time.Sleep(time.Millisecond * time.Duration(r.Intn(1000)+1000))
			t.Log("task", i)
		})
		if err != nil {
			t.Fatal(err)
		}
		waitGroup.Add(1)
		go func() {
			<-f
			waitGroup.Done()
		}()
	}
	go func() {
		time.Sleep(time.Second * 5)
		p.Close(ctx)
	}()
	waitGroup.Wait()
}
