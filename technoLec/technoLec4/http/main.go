package main

import (
	"net/http"
	"log"
	"io/ioutil"
)

func main()  {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handler(w http.ResponseWriter, r *http.Request)  {
	c := http.Client{}
	resp, err := c.Get("http://artii.herokuapp.com/make?text="+ r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	w.WriteHeader(200)
	w.Write(body)

}
