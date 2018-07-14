package main //4
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
			c<- <-a
		}
	}()
	go func() {
		for {
			c<- <-b
		}
	}()
	return c
}

func boring(s string) <-chan string  {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c<- fmt.Sprintf("%s %d", s, i)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}()
	return c
}

