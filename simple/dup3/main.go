package main

import (
	"os"
	"io/ioutil"
	"fmt"
	"strings"
)

func main()  {
	counts := make(map[string]int)
	for _, filename := range os.Args[1:] {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			for _, world := range strings.Split(line, " ") {

				counts[world]++
			}
		}
	}

	for line, n := range counts {
		fmt.Printf("%s\t%d\n", line, n)
	}
}
