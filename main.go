package main

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/Nv7-Github/scratchy/scratchy"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseDir(fset, "test", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	b := scratchy.NewBuilder()
	files := make([]*ast.File, 0, len(f["main"].Files))
	for _, file := range f["main"].Files {
		files = append(files, file)
	}
	err = b.Build("test", fset, files)
	if err != nil {
		panic(err)
	}
}
