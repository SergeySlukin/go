package main

import (
	"go/token"
	"go/parser"
	"log"
	"go/ast"
	"strings"
	"fmt"
)

func main()  {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, `.\10_ast\e.go`, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	// Decls содержит в себе список всех объявленных переменных и функций
	for _, decl := range f.Decls {
		// Если мы нашли функцию
		fdecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		// Проверяем, может быть это пример
		if isExample(fdecl) {
			//Проанализируем сигнатуру
			output, found := findExampleOutput(fdecl.Body, f.Comments)
			if found {
				fmt.Printf("%s need output '%s'\n", fdecl.Name.Name, output)
			}
		}
	}
}

func findExampleOutput(block *ast.BlockStmt, comments []*ast.CommentGroup) (string, bool)  {
	var last *ast.CommentGroup
	for _, group := range comments {
		if (block.Pos() < group.Pos()) && (block.End() > group.End()) {
			last = group
		}
	}
	if last != nil {
		text := last.Text()
		marker := "Output: "
		if strings.HasPrefix(text, marker) {
			return strings.TrimRight(text[len(marker):], "\n"), true
		}
	}
	return "", false
}

func isExample(fdecl *ast.FuncDecl) bool  {
	return strings.HasPrefix(fdecl.Name.Name, "Example")
}

