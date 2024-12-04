package workpool

import (
	"runtime"
	"sync"
)

var DefaultSize = 2 * runtime.NumCPU()

type WorkPool struct {
	wg  sync.WaitGroup
	sem chan int
}

func New(size int) WorkPool {
	return WorkPool{
		sync.WaitGroup{},
		make(chan int, size),
	}
}

func (pool *WorkPool) Spawn(task func()) {
	pool.wg.Add(1)
	pool.sem <- 0
	go func() {
		task()
		<-pool.sem
		pool.wg.Done()
	}()
}

func (pool *WorkPool) Wait() {
	pool.wg.Wait()
}

func (pool *WorkPool) Close() {
	close(pool.sem)
}
