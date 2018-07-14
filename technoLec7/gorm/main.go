package main

import "github.com/jinzhu/gorm"
import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
)

type TestUser struct {
	ID       int
	Username string
}

func main() {
	db, err := gorm.Open("mysql", "root@/golang_test?charset=utf8&parseTime=True&loc=Local")
	db.LogMode(true)
	if err != nil {
		panic(err)
	}

	db.CreateTable(&TestUser{})
	fmt.Println("Table created ", db.HasTable(&TestUser{}))
	db.DropTable(&TestUser{})
	defer db.Close()

	complex(db)
}

func complex(db *gorm.DB) {
	if !db.HasTable(&User{}) {
		db.CreateTable(&User{})
	}
	fill(db)
	query(db)
}

func fill(db *gorm.DB) {
	user := &User{
		FirstName: "Sergey",
		LastName:  "Slukin",
		Username:  "slukin",
		Salary:    5000,
	}

	db.Create(user)

	var users = []User{
		User{Username: "foobar", FirstName: "foo", LastName: "bar", Salary: 200},
		User{Username: "helloword", FirstName: "Hello", LastName: "World", Salary: 300},
		User{Username: "john", FirstName: "John", Salary: 200},
	}

	for _, user := range users {
		db.Create(&user)
	}
}

func query(db *gorm.DB) {
	u := User{FirstName: "Sergey"}
	db.Where(&u).First(&u)

	//db.Find(&u)
	fmt.Println(u)

	users := []User{}
	db.Where(&User{Salary: 200}).Find(&users)
}

type User struct {
	// auto-populate columns: id, created_at, updated_at, deleted_at
	gorm.Model

	Username string `sql:"type:VARCHAR(255);not null;unique"`

	// Set default value
	LastName string `sql:"DEFAULT: ''"`

	// Custom column name instead of default snake_case format
	FirstName string `gorm:"column:FirstName"`

	Role string `sql:"-"`

	Salary int64
}

func (u *User) TableName() string {
	return "users"
}

/*func (u *User) BeforeAction() (err error)  {
	if u.Role != "admin" {
		err = errors.New("Permission denied.")
	}
	return
}*/