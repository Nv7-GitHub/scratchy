package scratchy

import (
	"fmt"
	"go/importer"
	"go/types"
	"log"
)

type Sprite struct {
	Name string
}

// Searches for sprites and global variables
func (b *Builder) TypePass() error {
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check(b.path, b.fset, b.files, nil)
	if err != nil {
		log.Fatal(err)
	}
	s := pkg.Scope()

	fmt.Println(s.Lookup("Sprite").Type().Underlying().(*types.Struct))
	fmt.Println(types.NewMethodSet(types.NewPointer(s.Lookup("Sprite").Type())).At(0)) // Lookup methods for a pointer to the sprite
	fmt.Println(s.Child(0).Child(0).Child(0).Names())

	return nil
}
