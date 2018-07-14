package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"html/template"
)

func main()  {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handler(w http.ResponseWriter, r *http.Request)  {

	tmpl := template.New("main")
	tmpl, _ = tmpl.Parse(
		`<div style="display: inline-block; border: 1px solid #aaa; border-radius: 3px padding:30px; margin: 20px;"><pre>{{.}}</pre></div>`)

	path := r.URL.Path

	c := http.Client{}
	resp, err := c.Get("http://artii.herokuapp.com/make?text=" + path)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, string(body))
}
