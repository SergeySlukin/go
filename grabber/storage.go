package main

import (
	"os"
)

type Storage struct {
	filename string
	file *os.File
}

func NewStorage(filename string) *Storage {

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	return &Storage{filename:filename, file:file}
}

func (s *Storage) Write(v interface{})  {
	switch x := v.(type) {
	case []byte:
		s.file.Write(x)
	case string:
		s.file.WriteString(x + "\n")
	}
}

func (s *Storage) Close() {
	s.file.Close()
}