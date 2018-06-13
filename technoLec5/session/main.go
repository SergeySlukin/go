package main
import (
	"net/http"
	"html/template"
	"log"
	"time"
	"math/rand"
)



var letters = []rune("qwertyuiolkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM0123456789")

type User struct {
	Auth bool
	Username string
}

var sessions = map[string]*User{}

func main()  {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/get_cookie", cookieHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request)  {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		authTemplate, err := template.New("auth.html").ParseFiles("auth.html")
		if err != nil {
			log.Println("auth template error: ", err)
			http.Error(w, "auth template error", http.StatusInternalServerError)
		}
		authTemplate.Execute(w, nil)
	} else if err != nil {
		PanicOnError(err)
	}

	user, ok := sessions[session.Value]
	if !ok {
		authTemplate, err := template.New("auth.html").ParseFiles("auth.html")
		if err != nil {
			log.Println("auth template error: ", err)
			http.Error(w, "auth template error", http.StatusInternalServerError)
		}
		authTemplate.Execute(w, nil)
	} else {
		indexTemplate, err := template.New("index.html").ParseFiles("index.html")
		if err != nil {
			log.Println("index template error: ", err)
			http.Error(w, "index template error", http.StatusInternalServerError)
		}
		indexTemplate.Execute(w, user)
	}
}

func cookieHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == http.MethodPost {
		r.ParseForm()
		inputLogin := r.FormValue("login") //r.Form["login"][0]
		expires := time.Now().Add(365 * 24 * time.Hour)
		sessionId := RandomString(32)
		sessions[sessionId] = &User{
			Auth: true,
			Username: inputLogin,
		}
		cookie := http.Cookie{Name: "session_id", Value: sessionId, Expires: expires}
		http.SetCookie(w, &cookie)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func PanicOnError(err error)  {
	panic(err)
}
