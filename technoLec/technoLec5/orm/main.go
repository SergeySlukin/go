package main

import "github.com/jinzhu/gorm"
import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
)

var db *gorm.DB

type Student struct {
	ID    uint64 `sql:"AUTO_INCREMENT" gorm: "primary_key"`
	Fio   string
	Info  string
	Score int
}

func (u *Student) TableName() string {
	return "students"
}

func (u *Student) BeforeSave() (err error) {
	fmt.Println("trigger on before save")
	return
}

func PrintByID(id uint)  {
	st := Student{}
	err := db.Find(&st, id).Error
	if err == gorm.ErrRecordNotFound {
		fmt.Println("Record not found", id)
	} else {
		PanicOnErr(err)
	}
	fmt.Printf("PrintByID: %+v, data: %+v\n", id, st)
}

func main()  {
	var err error

	db, err = gorm.Open("mysql", "root@tcp(localhost:3306)/golang_test?charset=utf8")
	PanicOnErr(err)
	defer db.Close()
	db.DB()
	db.DB().Ping()

	PrintByID(1)
	PrintByID(100500)

	all := []Student{}
	db.Find(&all)
	for k, v := range all {
		fmt.Printf("students[%d] %+v\n", k, v)
	}

	newStudent := Student{
		Fio: "Ivan Ivanov",
	}
	db.Create(&newStudent)
	fmt.Println(newStudent.ID)
	PrintByID(uint(newStudent.ID))

	newStudent.Info = "update"
	newStudent.Score = 10
	db.Save(newStudent)
	PrintByID(uint(newStudent.ID))
}

func PanicOnErr(err error)  {
	if err != nil {
		panic(err)
	}
}
