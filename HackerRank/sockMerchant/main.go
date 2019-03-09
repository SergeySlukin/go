package main

import (
	"bufio"
	"os"
	"io"
	"strings"
	"strconv"
	"fmt"
)

func sockMerchant(n int32, ar []int32) int32 {
	var result int32
	socks := make(map[int32]int)
	for _, v := range ar {
		socks[v]++
	}
	for _, v := range socks {
		if v > 1 {
			result += int32(v / 2)
		}
	}
	return result
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	nTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)

	n := int32(nTemp)

	arTemp := strings.Split(readLine(reader), " ")
	var ar []int32

	for i := 0; i < int(n); i++ {
		arItemTemp, err := strconv.ParseInt(arTemp[i], 10, 64)
		checkError(err)
		arItem := int32(arItemTemp)
		ar = append(ar, arItem)
	}

	result := sockMerchant(n, ar)
	fmt.Fprintf(os.Stdout, "%d\n", result)
}

func readLine(reader *bufio.Reader) string {
	line, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(line), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
