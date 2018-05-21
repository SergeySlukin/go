package main
import (
	"sync"
	"fmt"
	"math/rand"
	"time"
)

func boring1(wg *sync.WaitGroup, die chan bool) <-chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case c <- fmt.Sprintf("Boring %d", rand.Intn(100)):
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			case <-die:
				fmt.Println("Done")
				wg.Done()
				return
			}
		}
	}()
	return c
}

func main() {
	die := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)
	resp1 := boring1(&wg, die)
	resp2 := boring1(&wg, die)
	for i := 0; i < 5; i++ {
		fmt.Println("resp 1", <-resp1)
		fmt.Println("resp 2", <-resp2)
	}
	die <- true
	die <- true

	wg.Wait()
}
