package main

import (
	"sync"
)

type FilePool struct {
	sync.Mutex
	size       int
	wg         sync.WaitGroup
	fileString chan string
	kill       chan bool
}

func NewFilePool(size int) *FilePool {
	p := &FilePool{
		fileString: make(chan string, 100),
		kill:       make(chan bool),
	}

	p.Resize(size)
	return p
}

func (p *FilePool) Resize(size int) {
	p.Lock()
	defer p.Unlock()
	for p.size < size {
		p.size++
		p.wg.Add(1)
		go p.worker()
	}
	for p.size > size {
		p.size--
		p.kill <- true
	}
}

func (p *FilePool) worker() {
	defer p.wg.Done()
	for {
		select {
		case s, ok := <-p.fileString:
			if !ok {
				return
			}
			statistics.Views++
			stringHandler(s)
		case <-p.kill:
			return
		}
	}
}

func (p *FilePool) Close() {
	close(p.fileString)
}

func (p *FilePool) Wait() {
	p.wg.Wait()
}
