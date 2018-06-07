package main

import (
	"net/http"
	"html/template"
	"fmt"
	"time"
	"math/rand"
)

type User struct {
	Auth     bool
	Username string
}

var sessions = map[string]*User{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")

		if err == http.ErrNoCookie {
			tmpl, _ := template.New("auth.html").ParseFiles("auth.html")
			tmpl.Execute(w, nil)
			fmt.Fprint(w, err.Error(), http.StatusInternalServerError)
			return
		} else if err != nil {
			PanicOnErr(err)
		}
		username, ok := sessions[sessionID.Value]
		if ok {
			tmpl, _ := template.New("index.html").ParseFiles("index.html")
			tmpl.Execute(w, username)
		} else {
			tmpl, _ := template.New("auth.html").ParseFiles("auth.html")
			tmpl.Execute(w, username)
		}

	})

	http.HandleFunc("/get_cookie", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		inputLogin := r.Form["login"][0]
		expiration := time.Now().Add(365 * 24 * time.Hour)
		sessionID := RandStringRunes(32)
		sessions[sessionID] = &User{
			Auth: true,
			Username: inputLogin,
		}
		cookie := http.Cookie{Name: "session_id", Value: sessionID, Expires: expiration}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	})
	http.ListenAndServe(":3000", nil)
}

func PanicOnErr(err error)  {
	if err != nil {
		panic(err)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(length int) string  {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
