package main

import (
	"log"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"fmt"
	"strconv"
)

var db *sql.DB
var dbConnError error

type Student struct{
	Id int64
	Fio string
	About sql.NullString
	Score int
}

func NewStudent() *Student {
	return &Student{}
}

var studentMap = make(map[int64]*Student)

func init()  {
	db, dbConnError = sql.Open("mysql", "root@tcp(localhost:3306)/simleGoMysql?charset=utf8")
}

func main()  {

	if dbConnError != nil {
		panic(dbConnError)
	}

	db.Ping()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/delete", deleteHandler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request)  {

	rows, err := db.Query("SELECT * from `students`")
	if err != nil {
		log.Println("Error in index page ", err)
	}
	for rows.Next() {

		student := NewStudent()
		err := rows.Scan(&student.Id, &student.Fio, &student.About, &student.Score)
		if err != nil {
			log.Println("error", err)
		}
		studentMap[student.Id] = student
	}
	rows.Close()

	indexTemplate, err := template.New("index.html").ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println(err)
	}

	indexTemplate.Execute(w, studentMap)
}

func addHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == http.MethodGet {
		addTemplate, err := template.New("add.html").ParseFiles("templates/add.html")
		if err != nil {
			fmt.Println(err)
		}
		addTemplate.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		inputFio := r.FormValue("fio")
		inputAbout := toNullString(r.FormValue("about"))
		inputScore := toInt(r.FormValue("score"))
		stmt, err := db.Prepare("INSERT INTO `students` (`fio`, `about`, `score`) VALUES (?, ?, ?)")
		if err != nil {
			fmt.Println(err)
		}
		result, err := stmt.Exec(inputFio, inputAbout, inputScore)
		if err != nil {
			fmt.Println(err)
		}

		affected, err := result.RowsAffected()
		fmt.Println("Insert RowAffected: ", affected)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request)  {

	if r.Method == http.MethodGet {
		idString := r.URL.Query()["id"][0]
		id, _ := strconv.ParseInt(idString, 10, 0)
		student, ok := studentMap[id]
		if !ok {
			student := NewStudent()
			row := db.QueryRow("SELECT * FROM `students` WHERE `id` = ?", id)
			err := row.Scan(student.Id, student.Fio, student.About, student.Score)
			if err != nil {
				fmt.Println(err)

				notFoundTemplate, err := template.New("404.html").ParseFiles("templates/404.html")
				if err != nil {
					fmt.Println(err)
				}
				notFoundTemplate.Execute(w, nil)
			} else {
				studentMap[student.Id] = student
				updateTemplate, err := template.New("update.html").ParseFiles("templates/update.html")
				if err != nil {
					fmt.Println(err)
				}
				updateTemplate.Execute(w, student)
			}
		} else {
			updateTemplate, err := template.New("update.html").ParseFiles("templates/update.html")
			if err != nil {
				fmt.Println(err)
			}
			updateTemplate.Execute(w, student)
		}

	} else if r.Method == http.MethodPost {
		r.ParseForm()
		id := r.FormValue("id")
		inputFio := r.FormValue("fio")
		inputAbout := r.FormValue("about")
		inputScore := r.FormValue("score")
		id64 := toInt64(id)
		student, ok := studentMap[id64]
		if !ok {
			student := NewStudent()
			student.Id, student.Fio, student.About, student.Score = id64, inputFio, toNullString(inputAbout), toInt(inputScore)
		} else {
			student.Fio, student.About, student.Score = inputFio, toNullString(inputAbout), toInt(inputScore)
		}
		studentMap[student.Id] = student
		stmt, err := db.Prepare("UPDATE students SET fio = ?, about = ?, score = ? WHERE id = ?")
		if err != nil {
			fmt.Println(err)
		}
		result, err := stmt.Exec(inputFio, inputAbout, inputScore, id)
		if err != nil {
			fmt.Println(err)
		}

		affected, err := result.RowsAffected()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Update RowAffrected: ", affected)
		http.Redirect(w, r, "/", http.StatusFound)

	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request)  {
		r.ParseForm()
		inputId := r.FormValue("id")
		id := toInt64(inputId)
		stmt, err := db.Prepare("DELETE FROM students WHERE id = ?")
		if err != nil {
			fmt.Println(err)
		}
		result, err := stmt.Exec(id)
		if err != nil {
			fmt.Println(err)
		}
		affected, err := result.RowsAffected()
		if affected == 1 {
			delete(studentMap, id)
		}
	http.Redirect(w, r, "/", http.StatusFound)
}

func toNullString(s string) sql.NullString {
	return sql.NullString{String:s, Valid: s != ""}
}

func toInt(s string) int  {
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
		n = 0
	}
	return n
}

func toInt64(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		fmt.Println(err)
		n = 0
	}
	return n
}