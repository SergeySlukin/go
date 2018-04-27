package main

import (
	"bufio"
	"os"
	"strconv"
	"fmt"
)

//https://www.hackerrank.com/challenges/30-conditional-statements/problem

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if n , err := strconv.Atoi(scanner.Text()); err == nil {
			isBool := isEven(n)
			if isBool {
				fmt.Print("Not Weird")
			} else {
				fmt.Print("Weird")
			}
		}
	}

}

func isEven(n int) bool {
	if n % 2 != 0 {
		return false
	} else {
		if n >= 2 && n < 5 {
			return true
		} else if n > 5 && n < 20 {
			return false
		} else if n > 20 {
			return true
		}
	}
	return false
}
