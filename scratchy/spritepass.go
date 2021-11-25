package scratchy

import (
	"go/ast"

	"github.com/Nv7-Github/scratch"
	"github.com/Nv7-Github/scratch/sprites"
	"github.com/Nv7-Github/scratchy/types"
)

func (p *Program) SpritePass(file *ast.File) error {
	for _, decl := range file.Decls {
		struc, ok := decl.(*ast.GenDecl)
		if ok {
			// Sprite or Variable
			_, ok := struc.Specs[0].(*ast.ValueSpec)
			if ok {
				// Variable!
				spec := struc.Specs[0].(*ast.ValueSpec)
				typ, err := p.ConvType(spec.Type)
				if err != nil {
					return err
				}

				for _, name := range spec.Names {
					v := p.MakeVariable(name.Name, true, typ, scratch.Stage.BasicSprite)
					p.GlobalVariables[v.Name] = v
				}
				continue
			}

			// Sprite, initialize sprite object!
			spec, ok := struc.Specs[0].(*ast.TypeSpec)
			// Not global variable or sprite, skip
			if !ok {
				continue
			}
			sprite := sprites.AddSprite(spec.Name.Name)
			spriteV := &Sprite{
				Name:      spec.Name.Name,
				Functions: make(map[string]*Function),
				Variables: make(map[string]*Variable),
				Sprite:    sprite,
			}
			p.Sprites[spriteV.Name] = spriteV
			p.Scope.Sprite = spriteV

			// Add Fields
			typ := spec.Type.(*ast.StructType)
			for _, field := range typ.Fields.List {
				typ, err := p.ConvType(field.Type)
				if err != nil {
					return err
				}

				// Add fields
				for _, name := range field.Names {
					v := p.MakeVariable(name.Name, false, typ, sprite.BasicSprite)
					spriteV.Variables[v.Name] = v
				}
			}
		}
	}

	return nil
}

func (p *Program) MakeVariable(name string, global bool, typ types.Type, sprite *sprites.BasicSprite) *Variable {
	var val interface{}
	nameV := p.GetVarName(name, global)
	switch typ.(type) {
	case types.BasicType:
		val = sprite.AddVariable(nameV, "")

	case *types.ArrayType:
		val = ArrayValue{sprite.AddList(nameV, []interface{}{})}

	case *types.MapType:
		val = MapValue{sprite.AddList(nameV+"_key", []interface{}{}), sprite.AddList(nameV+"_value", []interface{}{})}
	}

	return &Variable{
		Name:  nameV,
		Type:  typ,
		Value: val,
	}
}
