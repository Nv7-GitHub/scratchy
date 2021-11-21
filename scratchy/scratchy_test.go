package scratchy

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestScratchy(t *testing.T) {
	fset := token.NewFileSet()
	parsed, err := parser.ParseDir(fset, "../test", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	prog := newProgram()
	// Sprite pass
	for _, pkg := range parsed {
		for _, file := range pkg.Files {
			prog.SpritePass(file)
		}
	}

	// Print
	for _, v := range prog.GlobalVariables {
		spew.Config.MaxDepth = 2
		spew.Dump(v)
	}
	for _, sprite := range prog.Sprites {
		spew.Config.MaxDepth = 3
		spew.Dump(sprite)
	}
}
