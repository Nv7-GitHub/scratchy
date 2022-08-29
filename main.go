package main

import (
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
	err = b.Build(fset, f["main"].Files)
	if err != nil {
		panic(err)
	}
}
