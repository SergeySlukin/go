package main

import (
	"flag"
	"time"
	"os"
	"os/signal"
	"fmt"
)

func main() {

	workers := flag.Int("w", 2, "количество потоков")
	reportDuration := flag.Int("r", 10, "Частота отчетов (сек)")
	dup := flag.Int("d", 500, "кол-во дубликатов для остановки")
	hashFileName := flag.String("hf", "has.bin", "файл ъешей")
	quoteFileName := flag.String("qf", "quotes.txt", "файл записей")



	flag.Parse()

	quoteFileStorage := NewStorage(*quoteFileName)
	hashFileStorage := NewStorage(*hashFileName)

	defer quoteFileStorage.Close()
	defer hashFileStorage.Close()

	p := NewPool(*workers)

	ticker := time.NewTicker(time.Duration(*reportDuration) * time.Second)
	defer ticker.Stop()

	interruptChannel := make(chan os.Signal)
	signal.Notify(interruptChannel, os.Interrupt)

	hashStorage := NewHashStorage()
	hashStorage.SetAlgorithm("md5")

	hashFileStorage.Read("byte", hashStorage)

	for {
		select {
		case text := <-p.receive:
			hash, ok := hashStorage.Add(text)
			if !ok {

				if dupCount++; dupCount == *dup {
					go func() {
						interruptChannel <- os.Interrupt
					}()
					fmt.Println("Достигнут предел повтоов, завершаю работу. Всего записей: ", len(hashStorage.hashMap))
					p.Close()
				}
			} else {
				quoteCount++
				hashFileStorage.Write(hash)
				quoteFileStorage.Write(text)
				fmt.Println(text)
			}
		case <-interruptChannel:
			fmt.Println("CTRL-C: Завершаю работу. Всего записей: ", len(hashStorage.hashMap))
			ticker.Stop()
			return
		case <-ticker.C:
			fmt.Printf("Всего %d / Повторов %d (%d записей/сек)\n", len(hashStorage.hashMap), dupCount, quoteCount / *reportDuration)
			quoteCount = 0
		}
	}

	p.Wait()

}
