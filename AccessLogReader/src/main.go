package main

import (
	"sync"
	"time"
	"os"
	"fmt"
	"io/ioutil"
	"log"
	"bytes"
	"io"
	"encoding/json"
)

var mu sync.Mutex

/**
Структура для статистики
 */
type Statistics struct {
	Views       uint64            `json:"views"`
	Urls        int               `json:"urls"`
	Traffic     uint64            `json:"traffic"`
	Crawlers    map[string]uint64 `json:"crawlers"`
	StatusCodes map[uint16]uint64 `json:"status_codes"`
}

var statistics = &Statistics{
	Crawlers:    make(map[string]uint64),
	StatusCodes: make(map[uint16]uint64),
}

/**
Запускаем пул воркеров
Получаем имя файла, читаем его и передае строки в канал
В конце  выводим json и показываем за сколько времени выполнилось приложение
 */
func main() {

	start := time.Now()

	filePool := NewFilePool(4)

	filename := os.Args[1:]

	file, err := ioutil.ReadFile(filename[0])
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(file)
	for {
		line, err := buf.ReadString('\n')
		if len(line) == 0 {
			if err != nil {
				if err == io.EOF {
					break
				}
				break
			}
		}
		filePool.fileString <- line
	}

	filePool.Close()
	filePool.Wait()

	statistics.Urls = len(uniqueUrls)

	jsonEncoder := json.NewEncoder(os.Stdout)
	jsonEncoder.Encode(&statistics)
	fmt.Println(time.Since(start).Seconds())
}
