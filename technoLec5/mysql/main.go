package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

var (
	db *sql.DB
)

func PrintByID(id int64)  {
	var fio string
	var info sql.NullString
	var score int
	row := db.QueryRow("SELECT fio, info, score FROM students WHERE id = ?", id)
	err := row.Scan(&fio, &info, &score)
	PanicOnErr(err)
	fmt.Println("PrintByID:", id, "fio:", fio, "info:", info.String, "score:", score)
}

func main()  {
	var err error

	db, err = sql.Open("mysql", "root@tcp(localhost:3306)/golang_test?charset=utf8")
	PanicOnErr(err)
	fmt.Println("open connections: ", db.Stats().OpenConnections)
	err = db.Ping()
	db.SetMaxOpenConns(10)
	fmt.Println("open connections: ", db.Stats().OpenConnections)
	PanicOnErr(err)

	rows, err := db.Query("SELECT fio, score FROM students")
	PanicOnErr(err)
	for rows.Next(){
		var fio string
		var score int
		err = rows.Scan(&fio, &score)
		PanicOnErr(err)
		fmt.Println("rows.Next fio: ", fio, "score:", score)
	}
	rows.Close()

	var fio string
	row := db.QueryRow("SELECT fio FROM students WHERE id = 1")
	err = row.Scan(&fio)
	PanicOnErr(err)
	fmt.Println("db.QueryRow fio: ", fio)

	result, err := db.Exec("INSERT INTO students (`fio`) VALUES (?)", "Слукин Сергей")
	PanicOnErr(err)
	affected, err := result.RowsAffected()
	PanicOnErr(err)
	lastId, err := result.LastInsertId()
	PanicOnErr(err)

	fmt.Println("Insert - Rows Affected", affected, "Last InsertId: ", lastId)

	PrintByID(lastId)

	result, err = db.Exec("UPDATE students set info = ? where id = ?", "test user", lastId)
	PanicOnErr(err)

	affected, err = result.RowsAffected()
	PanicOnErr(err)
	fmt.Println("Update - Rows Affected", affected)
	PrintByID(lastId)

	stmt, err := db.Prepare("UPDATE students SET info = ?, score = ? WHERE id = ?")
	PanicOnErr(err)
	result, err = stmt.Exec("prepared statements update", lastId, lastId)
	PanicOnErr(err)

	affected, err = result.RowsAffected()
	PanicOnErr(err)
	fmt.Println("Update - RowsAffected", affected)

}

func PanicOnErr(err error)  {
	if err != nil {
		panic(err)
	}
}
