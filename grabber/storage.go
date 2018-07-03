package main

import (
	"os"
	"io"
	"encoding/hex"
	"fmt"
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

func (s *Storage) Read(t string, hashStorage *HashStorage)  {
	switch t {
	case "byte":
		readBytes(s, hashStorage)
	}
}

func readBytes(s *Storage, h *HashStorage)  {
	data := make([]byte, 16)
	for {
		n, err := s.file.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if n == 16 {
			h.hashMap[hex.EncodeToString(data)] = true
		}
	}
	fmt.Println("Завершено. Прочитано ъешей: ", len(h.hashMap))
}

func (s *Storage) Close() {
	s.file.Close()
}