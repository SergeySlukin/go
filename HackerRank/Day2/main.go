package main

import (
	"strconv"
	"bufio"
	"os"
	"fmt"
)

func main()  {

	var _ = strconv.Itoa

	var i uint64 = 4
	var d float64 = 4.0
	var s string = "HackerRank "

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		if n, err := strconv.ParseUint(scanner.Text(), 10, 64); err == nil {
			fmt.Println(i + (uint64(n)))
		} else if n, err := strconv.ParseFloat(scanner.Text(), 64); err == nil {
			fmt.Printf("%.1f \n",d + (float64(n)))
		} else {
			fmt.Println(s + scanner.Text())
		}
	}
	
}
