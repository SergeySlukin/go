package main

import (
	"time"
	"os"
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"strings"
)

func main() {

	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}

	for range os.Args[1:] {
		fmt.Print(<-ch)
	}

	fmt.Printf("%.2fs elapsed ", time.Since(start).Seconds())

}

func fetch(url string, ch chan<- string)  {
	start := time.Now()
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		ch<- fmt.Sprint(err)
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch<- fmt.Sprintf("%s error = %v\n", url, err)
		return
	}

	ch<- fmt.Sprintf("%.2fs %7d %s\n", time.Since(start).Seconds(), nbytes, url)

}




