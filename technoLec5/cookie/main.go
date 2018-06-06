package main //1
import (
	"net/http"
	"html/template"
	"time"
	"log"
)

type User struct {
	Auth     bool
	Username string
}

func NewUser() *User  {
	return &User{
		Auth: false,
	}
}

func main() {

	user := NewUser()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err == nil && err != http.ErrNoCookie {
			user.Auth = true
			user.Username = sessionID.Value
		}

		temp, err := template.New("template.html").ParseFiles("template.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		temp.Execute(w, user)
	})

	http.HandleFunc("/get_cookie", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		log.Println(r.Form)
		inputLogin := r.Form["login"][0]
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "session_id", Value: inputLogin, Expires: expiration}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.ListenAndServe(":3000", nil)

}
