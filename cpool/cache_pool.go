package cpool

import "context"

const (
	DefaultConcurrency   = 100
	DefaultChannelLength = 10000
)

type CallbackFunc func(interface{})

type CachePool struct {
	ctx   context.Context
	Close context.CancelFunc

	cacheChan chan interface{}
}

func New(ctx context.Context, concurrency int, callback CallbackFunc) *CachePool {
	cp := &CachePool{
		cacheChan: make(chan interface{}, DefaultChannelLength),
	}
	cp.ctx, cp.Close = context.WithCancel(ctx)
	for i := 0; i < concurrency; i++ {
		go cp.run(callback)
	}
	return cp
}

func (cp *CachePool) Add(data interface{}) bool {
	select {
	case cp.cacheChan <- data:
		return true
	default:
		return false
	}
}

func (cp *CachePool) run(callback CallbackFunc) {
	for {
		select {
		case <-cp.ctx.Done():
			return
		case data := <-cp.cacheChan:
			callback(data)
		}
	}
}
