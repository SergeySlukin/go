package main

import (
	"bufio"
	"os"
	"strings"
	"strconv"
	"io"
	"fmt"
)

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024 * 1024)

	var arr [][]int32
	for i := 0; i < 6; i++ {
		arrRowTemp := strings.Split(readLine(reader), " ")

		var arrRow []int32
		for _, arrRowItem := range arrRowTemp {
			arrItemTemp, err := strconv.ParseInt(arrRowItem, 10, 64)
			checkError(err)
			arrItem := int32(arrItemTemp)
			arrRow = append(arrRow, arrItem)
		}

		if len(arrRow) != int(6) {
			panic("Bad input")
		}

		arr = append(arr, arrRow)
	}

	var a int32
	var b int32
	for i := 0; i < 3; i++ {
		for k, v := range arr[i] {
			if k < 3 {
				a += v
			}
		}
	}
	for i := 3; i < 6; i++ {
		for k, v := range arr[i] {
			if k > 1 {
				b += v
			}

		}
	}

	if a > b {
		fmt.Println(a)
	} else {
		fmt.Println(b)
	}
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

