package main

import (
	"sessionManager/src/session"
	"net/http"
	"time"
	"log"
	"html/template"
)

var globalSession *session.Manager

func init()  {
	globalSession, _ = session.NewManager("memory", "gosessionid", 3600)
}

func main()  {

	http.HandleFunc("/login", loginHandler)

	httpServer := &http.Server{
		Addr: ":3000",
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(httpServer.ListenAndServe())
}

func loginHandler(w http.ResponseWriter, r *http.Request)  {
	sess := globalSession.SessionStart(w, r)
	r.ParseForm()
	if r.Method == http.MethodGet {
		t, _ := template.ParseFiles("login.html")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, sess.Get("username"))
	} else {
		sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/", 302)
	}
}

func count(w http.ResponseWriter, r *http.Request) {
	sess := globalSession.SessionStart(w, r)
	createtime := sess.Get("createtime")
	if createtime == nil {
		sess.Set("createtime", time.Now().Unix())
	} else if (createtime.(int64) + 360) < (time.Now().Unix()) {
		globalSession.SessionDestroy(w, r)
		sess = globalSession.SessionStart(w, r)
	}
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	t, _ := template.ParseFiles("count.html")
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("counnum"))
}