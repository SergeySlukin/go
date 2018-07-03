package main

import (
	"sync"
	"time"
)

type Pool struct {
	mu       sync.Mutex
	size     int
	wg       sync.WaitGroup
	t        time.Ticker
	receive  chan string
	kill     chan bool
}

func NewPool(size int) *Pool {
	p := &Pool{
		t:       *time.NewTicker(100 * time.Millisecond),
		receive: make(chan string, 10),
		kill:    make(chan bool),
	}

	p.Resize(size)
	return p
}

func (p *Pool) Resize(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.size < n {
		p.size++
		p.wg.Add(1)
		go p.worker()
	}
	for p.size > n {
		p.size--
		p.kill <- true
	}
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case <-p.t.C:
			Grab(p.receive)
		case <-p.kill:
			return
		}
	}
}

func (p *Pool) Close() {
	p.t.Stop()
	close(p.receive)
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
