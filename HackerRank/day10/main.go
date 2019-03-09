package main

import (
	"bufio"
	"os"
	"strconv"
	"io"
	"strings"
	"fmt"
	"sort"
)

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024 * 1024)
	nTemp, _ := strconv.ParseInt(readLine(reader), 10, 64)
	n := int32(nTemp)

	binaryNumber := strconv.FormatInt(int64(n), 2)
	sliceArray := make([]int, 0, 0)
	for _, v := range binaryNumber {
		if v == 49 {
			if len(sliceArray) == 0 {
				sliceArray = append(sliceArray, 1)
			} else {
				el := sliceArray[len(sliceArray) - 1]
				el++
				sliceArray[len(sliceArray) - 1] = el
			}
		} else {
			sliceArray = append(sliceArray, 0)
		}
	}

	sort.Ints(sliceArray)

	fmt.Println(sliceArray[len(sliceArray) - 1])

	fmt.Println(len(binaryNumber))
	fmt.Println(binaryNumber)
}

func readLine(reader *bufio.Reader) string {
	line, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(line), "\r\n")
}
