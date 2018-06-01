package main //4
import (
	"net/http"
	"log"
	"strconv"
	"html/template"
)

type Todo struct {
	Name string
	Done bool
}

var todos = []Todo{
	{"Learn Go", false},
	{"Go to web", false},
	{"...", false},
	{"Profit", false},
}

func main()  {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func IsNotDone(todo Todo) bool  {
	return !todo.Done
}

func handler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == http.MethodPost {
		param := r.FormValue("id")
		index, _ := strconv.ParseInt(param, 10, 0)
		todos[index].Done = true
	}

	templ, err := template.New("template.html").Funcs(template.FuncMap{"IsNotDone":IsNotDone}).ParseFiles("template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templ.Execute(w, todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}



}
