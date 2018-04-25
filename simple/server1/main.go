package main

import (
	"net/http"
	"log"
	"fmt"
)

func main()  {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Url path = %q\n", r.URL.Path)
}

