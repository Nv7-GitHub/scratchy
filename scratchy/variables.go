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
		p.Scope.Vars[name] = &Variable{
			Name:  vName,
			Type:  rhs.Type,
			Value: rhs,
		}
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
