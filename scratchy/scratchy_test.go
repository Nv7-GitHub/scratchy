package scratchy

import (
	"archive/zip"
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/Nv7-Github/scratch"
	"github.com/Nv7-Github/scratch/assets"
)

func TestScratchy(t *testing.T) {
	scratch.Clear()

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
