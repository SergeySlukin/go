package main //4
import (
	"sync"
	"fmt"
)

type Task interface {
	Execute()
}

type Pool struct {
	mu    sync.Mutex
	size  int
	tasks chan Task
	kill  chan bool
	wg    sync.WaitGroup
}

func NewPool(size int) *Pool {
	p := &Pool{
		tasks: make(chan Task, 128),
		kill:  make(chan bool),
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
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			task.Execute()
		case <-p.kill:
			return
		}
	}
}

func (p *Pool) Exec(task Task) {
	p.tasks <- task
}

func (p *Pool) Close() {
	close(p.tasks)
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

type ExampleTask string

func (e ExampleTask) Execute() {
	fmt.Println("executing ", string(e))
}

func main() {

	pool := NewPool(5)

	pool.Exec(ExampleTask("foo"))
	pool.Exec(ExampleTask("bar"))

	pool.Resize(3)
	pool.Resize(6)

	for i := 0; i < 20; i++ {
		pool.Exec(ExampleTask(fmt.Sprintf("additional_%d", i + 1)))
	}

	pool.Close()
	pool.Wait()
}
