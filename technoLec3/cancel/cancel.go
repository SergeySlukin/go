package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(die chan bool) <-chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case c <- fmt.Sprintf("boring %d", rand.Intn(100)):
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			case <-die:
				fmt.Println("Jobs done")
				die <- true
				return
			}
		}
	}()
	return c
}

func main() {
	die := make(chan bool)
	res := boring(die)
	for i := 0; i < 5; i++ {
		fmt.Println(<- res)
	}
	die <- true
	<-die
}
