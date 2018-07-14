package main

import (
	"go/token"
	"go/parser"
	"log"
	"go/doc"
	"fmt"
)

func main()  {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, `.\10_ast\e.go`,  nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	examples := doc.Examples(f)
	for _, example := range examples {
		fmt.Println(example.Name)
	}
}
