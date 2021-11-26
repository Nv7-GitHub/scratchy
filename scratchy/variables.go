package scratchy

import (
	"go/ast"
	"go/token"

	"github.com/Nv7-Github/scratch/blocks"
	styps "github.com/Nv7-Github/scratch/types"
	"github.com/Nv7-Github/scratch/values"
	"github.com/Nv7-Github/scratchy/types"
)

func (p *Program) AddAssignStmt(stmt *ast.AssignStmt) error {
	rhs, err := p.AddExpr(stmt.Rhs[0])
	if err != nil {
		return err
	}

	if stmt.Tok == token.DEFINE {
		name := stmt.Lhs[0].(*ast.Ident).Name
		vName := p.GetVarName(name, false)

		v := p.Scope.Sprite.Sprite.AddVariable(vName, "")
		set := p.Scope.Sprite.Sprite.NewSetVariable(v, rhs.Value)
		p.Scope.Stack.Add(set)

		val := &Variable{
			Name:  vName,
			Type:  rhs.Type,
			Value: v,
		}
		p.Scope.Vars[name] = val
		p.Scope.Sprite.Variables[vName] = val
		return nil
	}

	if stmt.Tok != token.ASSIGN {
		return p.NewError(stmt.TokPos, "unknown assign token: %s", stmt.Tok.String())
	}

	// Get variable, assign to it
	l := stmt.Lhs[0]
	switch e := l.(type) {
	case *ast.Ident, *ast.SelectorExpr: // Normal type
		v, err := p.GetBasicVariable(l)
		if err != nil {
			return err
		}
		_, ok := v.Type.(types.BasicType)
		if !ok {
			return p.NewError(l.Pos(), "cannot assign composite value directly")
		}
		if !rhs.Type.Equal(v.Type) {
			return p.NewError(stmt.Pos(), "cannot assign type %s to variable type %s", rhs.Type.String(), v.Type.String())
		}

		blk := p.Scope.Sprite.Sprite.NewSetVariable(v.Value.(*styps.Variable), rhs.Value)
		p.Scope.Stack.Add(blk)

	case *ast.IndexExpr: // Array or Map
		v, err := p.GetBasicVariable(e.X)
		if err != nil {
			return err
		}
		switch typ := v.Type.(type) {
		case *types.ArrayType:
			if !typ.ValueType.Equal(rhs.Type) {
				return p.NewError(stmt.Rhs[0].Pos(), "cannot assign type %s to array element type of %s", rhs.Type.String(), typ.ValueType.String())
			}
			ind, err := p.AddExpr(e.Index)
			if err != nil {
				return err
			}
			if !ind.Type.Equal(types.NUMBER) {
				return p.NewError(e.Index.Pos(), "expected index to be type %s, got type %s", types.NUMBER.String(), ind.Type.String())
			}
			blk := p.Scope.Sprite.Sprite.NewReplaceInList(v.Value.(ArrayValue).Val, ind.Value, rhs.Value)
			p.Scope.Stack.Add(blk)

		case *types.MapType:
			k, err := p.AddExpr(e.Index)
			if err != nil {
				return err
			}
			if !typ.KeyType.Equal(k.Type) {
				return p.NewError(e.Index.Pos(), "expected key to be type %s, got type %s", typ.KeyType.String(), k.Type.String())
			}
			if !typ.ValType.Equal(rhs.Type) {
				return p.NewError(stmt.Rhs[0].Pos(), "expected value to be type %s, got type %s", typ.ValType.String(), rhs.Type.String())
			}

			// TODO: Map logic (get key index in array1, if 0 then add key to array1, val to array2, else set that index in array2 to val)
			s := p.Scope.Sprite.Sprite
			st := p.Scope.Stack
			v := v.Value.(MapValue)

			// Get index & add if
			ind := s.NewFindInList(v.Key, k.Value)
			st.Add(ind)
			comp := s.NewCompare(values.NewBlockValue(ind), values.NewIntValue(0), blocks.CompareEqual)
			st.Add(comp)
			ife := s.NewIfElse(values.NewBlockValue(comp))
			st.Add(ife)

			// If not exists
			addKey := s.NewAddToList(v.Key, k.Value)
			ife.Add(addKey)
			addVal := s.NewAddToList(v.Val, rhs.Value)
			ife.Add(addVal)

			// If exsits
			setKey := s.NewReplaceInList(v.Val, values.NewBlockValue(ind), rhs.Value)
			ife.Else.Add(setKey)
		default:
			return p.NewError(e.X.Pos(), "cannot index non-composite type")
		}
	}

	return nil
}

func (p *Program) GetBasicVariable(expr ast.Expr) (*Variable, error) {
	var v *Variable
	switch l := expr.(type) {
	case *ast.Ident:
		val, exists := p.Scope.Vars[l.Name]
		if !exists {
			return nil, p.NewError(l.Pos(), "undefined variable: %s", l.Name)
		}
		v = val

	case *ast.SelectorExpr:
		spriteName := l.X.(*ast.Ident).Name
		varName := l.Sel.Name

		if p.Scope.Fn.SpriteName != spriteName {
			return nil, p.NewError(l.X.Pos(), "unknown sprite: %s", spriteName)
		}
		var exists bool
		v, exists = p.Scope.Sprite.Variables[varName]
		if !exists {
			return nil, p.NewError(l.Sel.Pos(), "unknown sprite variable: %s", varName)
		}
	}
	return v, nil
}

func (p *Program) AddGetVar(v *Variable, pos token.Pos) (*types.Value, error) {
	var val styps.Value
	_, ok := v.Type.(types.BasicType)
	if ok {
		_, ok := v.Value.(*blocks.ScratchParamValue)
		if ok {
			val = v.Value.(*blocks.ScratchParamValue) // It's a function param
		} else {
			val = values.NewVariableValue(v.Value.(*styps.Variable))
		}
	} else {
		_, ok := v.Type.(*types.ArrayType)
		if ok {
			val = values.NewListValue(v.Value.(ArrayValue).Val)
		} else {
			return nil, p.NewError(pos, "cannot get value of type %s", v.Type.String())
		}
	}

	return &types.Value{
		Type:  v.Type,
		Value: val,
	}, nil
}

func (p *Program) AddIdent(expr *ast.Ident) (*types.Value, error) {
	v, exists := p.Scope.Vars[expr.Name]
	if !exists {
		return nil, p.NewError(expr.Pos(), "undefined variable: %s", expr.Name)
	}
	return p.AddGetVar(v, expr.Pos())
}

func (p *Program) AddSelector(expr *ast.SelectorExpr) (*types.Value, error) {
	spriteName := expr.X.(*ast.Ident).Name
	varName := expr.Sel.Name
	if p.Scope.Fn.SpriteName != spriteName {
		return nil, p.NewError(expr.X.Pos(), "unknown sprite: %s", spriteName)
	}

	v, exists := p.Scope.Sprite.Variables[varName]
	if !exists {
		return nil, p.NewError(expr.Sel.Pos(), "unknown sprite variable: %s", varName)
	}
	return p.AddGetVar(v, expr.Sel.Pos())
}

func (p *Program) AddIndex(expr *ast.IndexExpr) (*types.Value, error) {
	v, err := p.GetBasicVariable(expr.X)
	if err != nil {
		return nil, err
	}
	ind, err := p.AddExpr(expr.Index)
	if err != nil {
		return nil, err
	}

	switch typ := v.Type.(type) {
	case *types.ArrayType:
		if !ind.Type.Equal(types.NUMBER) {
			return nil, p.NewError(expr.Index.Pos(), "expected index to be type %s, got type %s", types.NUMBER.String(), ind.Type.String())
		}

		blk := p.Scope.Sprite.Sprite.NewItemOfList(v.Value.(ArrayValue).Val, ind.Value)
		p.Scope.Stack.Add(blk)
		return &types.Value{
			Type:  typ.ValueType,
			Value: values.NewBlockValue(blk),
		}, nil

	case *types.MapType:
		if !ind.Type.Equal(typ.KeyType) {
			return nil, p.NewError(expr.Index.Pos(), "expected key to be type %s, got type %s", typ.KeyType.String(), ind.Type.String())
		}
		keyInd := p.Scope.Sprite.Sprite.NewFindInList(v.Value.(MapValue).Key, ind.Value)
		p.Scope.Stack.Add(keyInd)

		blk := p.Scope.Sprite.Sprite.NewItemOfList(v.Value.(MapValue).Val, values.NewBlockValue(keyInd))
		p.Scope.Stack.Add(blk)

		return &types.Value{
			Type:  typ.ValType,
			Value: values.NewBlockValue(blk),
		}, nil

	default:
		return nil, p.NewError(expr.Pos(), "cannot index type %s", v.Type.String())
	}
}

func (p *Program) IncDecStmt(stmt *ast.IncDecStmt) error {
	v, err := p.GetBasicVariable(stmt.X)
	if err != nil {
		return err
	}

	if !v.Type.Equal(types.NUMBER) {
		return p.NewError(stmt.Pos(), "cannot increment/decrement non-number type")
	}

	if stmt.Tok == token.INC {
		blk := p.Scope.Sprite.Sprite.NewChangeVariable(v.Value.(*styps.Variable), values.NewIntValue(1))
		p.Scope.Stack.Add(blk)
	} else {
		blk := p.Scope.Sprite.Sprite.NewChangeVariable(v.Value.(*styps.Variable), values.NewIntValue(-1))
		p.Scope.Stack.Add(blk)
	}

	return nil
}
