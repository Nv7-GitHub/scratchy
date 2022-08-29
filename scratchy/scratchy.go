package scratchy

import (
	"go/ast"
	"go/token"

	"github.com/Nv7-Github/scratch"
)

type Builder struct {
	files map[string]*ast.File
	fset  *token.FileSet
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(fset *token.FileSet, files map[string]*ast.File) error {
	scratch.Clear()
	b.fset = fset
	b.files = files

	return nil
}
