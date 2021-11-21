package scratchy

import (
	"go/ast"

	"github.com/Nv7-Github/scratch/blocks"
)

func (p *Program) FunctionPass(file *ast.File) error {
	for _, decl := range file.Decls {
		_, ok := decl.(*ast.FuncDecl)
		if ok {
			// Function
			fn := decl.(*ast.FuncDecl)
			gfn, err := p.GetGlobalFunction(fn)
			if err != nil {
				return err
			}
			if fn.Recv == nil {
				// Global function
				p.GlobalFunctions[gfn.Name] = &gfn
				continue
			}

			// Get sprite
			spriteName := fn.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident)
			sprite, ok := p.Sprites[spriteName.Name]
			if !ok {
				return p.NewError(spriteName.Pos(), "unknown sprite: %s", spriteName)
			}

			// Local function
			parTypes := []blocks.FunctionParameter{blocks.NewFunctionParameterLabel(gfn.Name)}
			for _, par := range gfn.Params {
				switch par.Type.(BasicType) {
				case NUMBER, STRING:
					parTypes = append(parTypes, blocks.NewFunctionParameterValue(par.Name, blocks.FunctionParameterBool, ""))

				case BOOL:
					parTypes = append(parTypes, blocks.NewFunctionParameterValue(par.Name, blocks.FunctionParameterBool, false))
				}
			}
			fnVal := sprite.Sprite.NewFunction(parTypes...)
			sprite.Functions[gfn.Name] = &Function{
				GlobalFunction: gfn,
				ScratchFuntion: fnVal,
			}
		}
	}
	return nil
}

func (p *Program) GetGlobalFunction(fn *ast.FuncDecl) (GlobalFunction, error) {
	name := fn.Name.Name
	typ := fn.Type

	params := make([]Param, 0, len(fn.Type.Params.List))
	for _, par := range typ.Params.List {
		typ, err := p.ConvType(par.Type)
		if err != nil {
			return GlobalFunction{}, err
		}
		_, ok := typ.(BasicType)
		if !ok {
			return GlobalFunction{}, p.NewError(par.Type.Pos(), "currently composite values can't be function parameters")
		}

		for _, name := range par.Names {
			params = append(params, Param{
				Name: name.Name,
				Type: typ,
			})
		}
	}

	retTyp, err := p.ConvType(typ.Results.List[0].Type)
	if err != nil {
		return GlobalFunction{}, err
	}
	_, ok := retTyp.(BasicType)
	if !ok {
		return GlobalFunction{}, p.NewError(typ.Results.List[0].Pos(), "currently composite values can't be function returns")
	}
	return GlobalFunction{
		Name:    name,
		Params:  params,
		RetType: retTyp,
	}, nil
}
