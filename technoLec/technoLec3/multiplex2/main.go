package main

import (
	"fmt"
	"time"
	"math/rand"
)

func main() {
	c := fanIn(boring("Sergey"), boring("Ann"))
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}
	fmt.Println("Finish")
}

func fanIn(a, b <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case msg := <-a:
				c <- msg
			case msg := <-b:
				c <- msg
			}
		}
	}()
	return c
}

func boring(s string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			c <- s
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}()
	return c
}
