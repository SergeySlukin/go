package main

import (
	"bufio"
	"strings"
	"os"
	"strconv"
	"fmt"
)

func main() {

	phone := make(map[string]string)

	reader := bufio.NewReader(os.Stdin)

	lineN := readLine(reader)
	n, err := strconv.Atoi(lineN)
	checkError(err)

	for i := 0; i < n; i++ {
		line := readLine(reader)
		kv := strings.Split(line, " ")
		if len(kv) > 1 {
			phone[kv[0]] = kv[1]
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if name, ok := phone[line]; ok {
			fmt.Printf("%s=%s\n", line, name)
		} else {
			fmt.Println("Not Found")
		}
	}

}

func readLine(reader *bufio.Reader) string {
	line, _,  err := reader.ReadLine()
	checkError(err)

	return strings.TrimRight(string(line), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
