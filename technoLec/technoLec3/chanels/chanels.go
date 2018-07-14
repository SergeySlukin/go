package main //2
import "fmt"

var c chan int

func main() {
	c := make(chan string)

	go greet(c)

	for i := 0; i < 5; i++ {
		fmt.Println(<-c, ',', <-c)
	}

	stuff := make(chan int, 7)
	for j := 0; j < 19; j = j + 3 {
		stuff <- j
	}

	close(stuff)
	fmt.Println("Res", process(stuff))

}

func greet(c chan <- string) {
	for {
		c <- fmt.Sprintf("Текст 1")
		c <- fmt.Sprintf("Text 2")
	}
}

func process(input <-chan int) (res int)  {
	for r := range input {
		res += r
	}
	return
}