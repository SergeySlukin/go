package main

import (
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"strconv"
)

type Todo struct {
	Name string `json:"name"`
	Done bool   `json:"done"`
}

var todos = []Todo{
	{"Learn Go", false},
}

func main()  {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	fileContents, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.Write(fileContents)
	})

	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request", r.URL.Path)
		defer r.Body.Close()

		switch r.Method {
		case http.MethodGet:
			productsJson, _ := json.Marshal(todos)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(productsJson)
		case http.MethodPost:
			decoder := json.NewDecoder(r.Body)
			todo := Todo{}
			err := decoder.Decode(&todo)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			todos = append(todos, todo)
		case http.MethodPut:
			id := r.URL.Path[len("/todos/"):]
			index, _ := strconv.ParseInt(id, 10, 0)
			todos[index].Done = true
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}