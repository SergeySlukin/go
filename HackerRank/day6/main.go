package main

import (
	"bufio"
	"os"
	"fmt"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	for scanner.Scan() {
		var evenString string
		var oddString string
		for index, char := range scanner.Text() {
			if index % 2 == 0 {
				evenString += string(char)
			} else {
				oddString += string(char)
			}
		}
		fmt.Printf("%v %v\n", evenString, oddString)
	}
}