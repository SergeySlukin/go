package main

import (
	"os"
	"fmt"
	"strconv"
)

func main()  {
	s, sep := "", ""
	for k, arg := range os.Args[0:] {
		s += sep + strconv.Itoa(k) + "=" + arg
		sep = " "
	}
	fmt.Println(s)
}
