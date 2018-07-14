package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"strconv"
)

type App struct {
	Id     string `json:"id"`
	Title  string `json:"title, omitempty"`
	Data   int    `json:"-"`
	Flag   int    `json:"flag,string"`
	hidden string
}

func main() {
	data := []byte(`
	{
		"id": "asdq",
		"title": "My awesome app"
	}
	`)

	var app App

	err := json.Unmarshal(data, &app)
	if err != nil {
		panic(err)
	}

	fmt.Printf("App: %+v\n", app)

	embedding()
	pointers()
	custom()
}

type App1 struct {
	Id string `json:"id"`
}

type Org struct {
	Name string `json:"name"`
}

type AppWithOrg struct {
	App1
	Org
}

func embedding() {
	data := []byte(`
	{
		"id": "asdadasd",
		"name": "test"
	}
	`)

	var app AppWithOrg

	err := json.Unmarshal(data, &app)
	if err != nil {
		panic(err)
	}

	fmt.Printf("AppWithOrg: %+v\n", app)
}

func pointers() {
	var parsed map[string]interface{}

	data := []byte(`
	{
		"id": "asd",
		"age": "28"
	}
	`)

	err := json.Unmarshal(data, &parsed)
	if err != nil {
		panic(err)
	}
	fmt.Println("id:", parsed["id"].(string))
}

func custom() {
	data := []byte(`{"Month": "04/2018"}`)
	var d day

	err := json.Unmarshal(data, &d)
	if err != nil {
		panic(err)
	}
	fmt.Println("Got month:", d.Month)

}

type day struct {
	Month Month
}

type Month struct {
	MonthNumber int64
	YearNumber  int64
}

func (m Month) String() string {
	return fmt.Sprintf("%d/%d", m.MonthNumber, m.YearNumber)
}

func (m Month) MarshalJSON() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *Month) UnmarshalJSON(value []byte) (err error) {
	if len(value) < 2 {
		return fmt.Errorf("bad data")
	}
	if value[0] == '"' {
		value = append(value[:0], value[1:]...)
	}

	if value[len(value)-1] == '"' {
		value = value[0 : len(value)-1]
	}

	parts := strings.Split(string(value), "/")
	m.MonthNumber, err = strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return err
	}
	m.YearNumber, err = strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return err
	}
	return nil
}
