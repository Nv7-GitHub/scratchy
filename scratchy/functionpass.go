package scratchy

import (
	"go/ast"

	"github.com/Nv7-Github/scratch/blocks"
	"github.com/Nv7-Github/scratchy/types"
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
			nm := fn.Recv.List[0].Names[0].Name
			sprite, ok := p.Sprites[spriteName.Name]
			if !ok {
				return p.NewError(spriteName.Pos(), "unknown sprite: %s", spriteName)
			}

			// Local function
			fnV, err := p.GetSpriteFunction(sprite, gfn, nm)
			if err != nil {
				return err
			}
			sprite.Functions[fnV.Name] = fnV
		}
	}
	return nil
}

func (p *Program) GetSpriteFunction(sprite *Sprite, gfn GlobalFunction, spriteName string) (*Function, error) {
	parTypes := []blocks.FunctionParameter{blocks.NewFunctionParameterLabel(gfn.Name)}
	for _, par := range gfn.Params {
		switch par.Type.(types.BasicType) {
		case types.NUMBER, types.STRING:
			parTypes = append(parTypes, blocks.NewFunctionParameterValue(par.Name, blocks.FunctionParameterString, ""))

		case types.BOOL:
			parTypes = append(parTypes, blocks.NewFunctionParameterValue(par.Name, blocks.FunctionParameterBool, false))
		}
	}
	fnVal := sprite.Sprite.NewFunction(parTypes...)

	fn := &Function{
		GlobalFunction:  gfn,
		ScratchFunction: fnVal,
		SpriteName:      spriteName,
	}
	if gfn.RetType != nil {
		retValName := p.GetVarName("return_"+gfn.Name, false)
		val := sprite.Sprite.AddVariable(retValName, "")
		fn.ReturnVal = val
	}

	return fn, nil
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
		_, ok := typ.(types.BasicType)
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

	var retType types.Type
	var err error
	if typ.Results != nil && len(typ.Results.List) > 0 {
		retType, err = p.ConvType(typ.Results.List[0].Type)
		if err != nil {
			return GlobalFunction{}, err
		}
		_, ok := retType.(types.BasicType)
		if !ok {
			return GlobalFunction{}, p.NewError(typ.Results.List[0].Pos(), "currently composite values can't be function returns")
		}
	}
	return GlobalFunction{
		Name:    name,
		Params:  params,
		RetType: retType,
		Code:    fn.Body,
	}, nil
}
