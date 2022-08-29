package scratchy

import (
	"go/ast"
	"go/token"

	"github.com/Nv7-Github/scratch"
)

type Builder struct {
	files []*ast.File
	fset  *token.FileSet
	path  string

	sprites map[string]*Sprite
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(path string, fset *token.FileSet, files []*ast.File) error {
	scratch.Clear()
	b.fset = fset
	b.files = files
	b.path = path

	return b.TypePass()
}
