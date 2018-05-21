package main //4
import (
	"fmt"
	"time"
	"math/rand"
)

func main() {
	c := boring("message")
	for i := 0; i < 5; i++ {
		fmt.Printf("%q \n", <-c)
	}
	fmt.Println("Finish")
}

func boring(s string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c<-fmt.Sprintf("%s %d", s, i)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}()
	return c
}
