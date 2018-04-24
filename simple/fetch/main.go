package main

import (
	"os"
	"net/http"
	"fmt"
	"io"
	"strings"
)

func main()  {
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fetch: %v\n", err)
			os.Exit(1)
		}
		b, err := io.Copy(os.Stdout, resp.Body)
		status := resp.StatusCode
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Url reading error, url: %s %v\n", url, err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", b)
		fmt.Printf("%d\n", status)
	}
}