package scratchy

import (
	"archive/zip"
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/Nv7-Github/scratch"
	"github.com/Nv7-Github/scratch/assets"
	"github.com/Nv7-Github/scratchy/functions"
)

func TestScratchy(t *testing.T) {
	// Save bindings
	f, err := os.Create("../scratch/bindings.go")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(functions.GenBindings())
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	// Create project
	scratch.Clear()

	fset := token.NewFileSet()
	parsed, err := parser.ParseDir(fset, "../test", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	prog := newProgram(fset)
	// Sprite pass
	for _, pkg := range parsed {
		for _, file := range pkg.Files {
			prog.SpritePass(file)
		}
	}
	// Function pass
	for _, pkg := range parsed {
		for _, file := range pkg.Files {
			prog.FunctionPass(file)
		}
	}

	// Add BGs
	scratch.Stage.AddCostume(assets.CostumeBlank("background1"))
	for _, sprite := range prog.Sprites {
		sprite.Sprite.AddCostume(assets.CostumeScratchCat("cat"))
	}

	// Code Pass
	err = prog.CodePass()
	if err != nil {
		t.Fatal(err)
	}

	// Save
	out, err := os.Create("../Project.sb3")
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	zip := zip.NewWriter(out)
	defer zip.Close()

	err = scratch.Save(zip)
	if err != nil {
		t.Fatal(err)
	}
}
