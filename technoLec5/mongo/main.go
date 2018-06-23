package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

var session *mgo.Session

type student struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Fio   string        `json:"fio" bson:"fio"`
	Info  string        `json:"info" bson:"info"`
	Score int           `json:"score" bson:"score"`
}

func main() {
	var err error
	session, err := mgo.Dial("mongodb://localhost")

	//если коллекции не будет, то она создасться автоматически
	collection := session.DB("test").C("students")

	index := mgo.Index{
		Key: []string{"fio"},
	}
	err = collection.EnsureIndex(index)
	PanicOnError(err)

	if n, _ := collection.Count(); n == 0 {
		firstStudent := &student{
			bson.NewObjectId(),
			"Sergey Slukin",
			"work: mail.ru group",
			10,
		}
		err = collection.Insert(firstStudent)
		PanicOnError(err)
	}

	var allStudents []student
	// bson.M{} - условие для поиска
	err = collection.Find(bson.M{}).All(&allStudents)
	PanicOnError(err)
	for k, v := range allStudents {
		fmt.Printf("student[%d]: %+v\n", k, v)
	}

	id := bson.NewObjectId()
	// bson.M{"_id": id} - задание условия поиска
	var nonExistenStudent student
	err = collection.Find(bson.M{"_id":id}).One(&nonExistenStudent)
	if err == mgo.ErrNotFound {
		fmt.Println("student ", err)
	} else if err != nil {
		PanicOnError(err)
	}

	secondStudent := &student{id, "Иван Иванов", "", 0}
	err = collection.Insert(secondStudent)
	PanicOnError(err)

	err = collection.Find(bson.M{"_id":id}).One(&nonExistenStudent)
	if err == mgo.ErrNotFound {
		fmt.Println("student ", err)
	} else if err != nil {
		PanicOnError(err)
	}
	fmt.Printf("student %+v\n", nonExistenStudent)

	secondStudent.Info = "all records"
	collection.UpdateAll(
		bson.M{"fio": "Иван Иванов"},
		bson.M{
			"$set": bson.M{"info": "all Иван info"},
		},
	)

	secondStudent.Info = "single record"
	collection.Update(bson.M{"_id": secondStudent.ID}, &secondStudent)

	err = collection.Find(bson.M{"_id": secondStudent.ID}).One(&nonExistenStudent)
	PanicOnError(err)
	fmt.Printf("Second Student after update: %+v\n", nonExistenStudent)



}

func PanicOnError(err error)  {
	if err != nil {
		panic(err)
	}
}
