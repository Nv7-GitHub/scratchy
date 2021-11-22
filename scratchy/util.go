package scratchy

import (
	"go/ast"
	"strconv"

	"github.com/Nv7-Github/scratchy/types"
)

func (p *Program) GetVarName(name string, global bool) string {
	var dat map[string]*Variable
	if global {
		dat = p.GlobalVariables
	} else {
		dat = p.CurrSprite.Variables
	}

	_, exists := dat[name]
	if exists {
		i := 0
		var nameV string
		for exists {
			nameV = name + strconv.Itoa(i)
			_, exists = dat[nameV]
			i++
		}
		return nameV
	}
	return name
}

func (p *Program) ConvType(typ ast.Expr) (types.Type, error) {
	switch t := typ.(type) {
	case *ast.Ident:
		name := t.Name
		bTyp, ok := basicTypeConv[name]
		if !ok {
			return nil, p.NewError(typ.Pos(), "unknown type: %s", name)
		}
		return bTyp, nil

	case *ast.ArrayType:
		valTyp := t.Elt.(*ast.Ident).Name
		valType, ok := basicTypeConv[valTyp]
		if !ok {
			return nil, p.NewError(typ.Pos(), "unknown array value type: %s", valTyp)
		}
		return types.NewArrayType(valType), nil

	case *ast.MapType:
		keyTyp := t.Key.(*ast.Ident).Name
		valTyp := t.Value.(*ast.Ident).Name
		keyType, ok := basicTypeConv[keyTyp]
		if !ok {
			return nil, p.NewError(typ.Pos(), "unknown map key type: %s", keyTyp)
		}
		valType, ok := basicTypeConv[valTyp]
		if !ok {
			return nil, p.NewError(typ.Pos(), "unknown map value type: %s", valTyp)
		}
		return types.NewMapType(keyType, valType), nil

	default:
		return nil, p.NewError(typ.Pos(), "unknown type: %v", typ)
	}
}

var basicTypeConv = map[string]types.BasicType{
	"bool":    types.BOOL,
	"float64": types.NUMBER,
	"string":  types.STRING,
}
