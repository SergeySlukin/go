package main

import "fmt"

func main()  {

	fmt.Println("start")

	go process(0)

	go func() {
		fmt.Println("Anonim run")
	}()

	for i := 0; i < 100000; i++ {
		go process(i)
	}

	fmt.Sprintln()
}

func process(i int) {
	fmt.Println(i)
}

