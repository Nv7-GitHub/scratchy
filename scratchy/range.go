package scratchy

import (
	"go/ast"

	"github.com/Nv7-Github/scratch/blocks"
	styps "github.com/Nv7-Github/scratch/types"
	"github.com/Nv7-Github/scratch/values"
	"github.com/Nv7-Github/scratchy/types"
)

func (p *Program) AddRange(stmt *ast.RangeStmt) error {
	key := stmt.Key.(*ast.Ident).Obj.Decl.(*ast.AssignStmt)
	val, err := p.GetBasicVariable(key.Rhs[0].(*ast.UnaryExpr).X)
	if err != nil {
		return err
	}
	// Check type of rhs
	if !val.Type.BasicType().Equal(types.MAP) && !val.Type.BasicType().Equal(types.ARRAY) && !val.Type.Equal(types.STRING) {
		return p.NewError(key.Rhs[0].Pos(), "range must have composite value or string")
	}

	// Make decls
	var iterV, valDecl *styps.Variable
	var keyDeclName, valDeclName *string
	iterV = p.Scope.Sprite.Sprite.AddVariable(p.GetVarName("_rangeiter", false), 1)
	p.Scope.Stack.Add(p.Scope.Sprite.Sprite.NewSetVariable(iterV, values.NewIntValue(1)))

	for i, decl := range key.Lhs {
		d := decl.(*ast.Ident)
		if d.Name != "_" {
			if i == 0 {
				keyDeclName = &d.Name
			} else {
				name := p.GetVarName(d.Name, false)
				valDecl = p.Scope.Sprite.Sprite.AddVariable(name, "")
				valDeclName = &d.Name
			}
		}
	}
	// Decl types
	var keyDecl *styps.Variable
	var keyDeclType, valDeclType types.Type
	switch t := val.Type.(type) {
	case *types.BasicType:
		if keyDeclName != nil {
			keyDeclType = types.NUMBER
			keyDecl = iterV
		}
		if valDecl != nil {
			valDeclType = types.STRING
		}

	case *types.ArrayType:
		if keyDeclName != nil {
			keyDeclType = types.NUMBER
			keyDecl = iterV
		}
		if valDecl != nil {
			valDeclType = t.ValueType
		}

	case *types.MapType:
		if keyDeclName != nil {
			keyDeclType = t.KeyType
			name := p.GetVarName(*keyDeclName, false)
			keyDecl = p.Scope.Sprite.Sprite.AddVariable(name, "")
		}
		if valDecl != nil {
			valDeclType = t.ValType
		}
	}

	var blk types.Stack
	switch val.Type.BasicType() {
	case types.STRING:
		length := p.Scope.Sprite.Sprite.NewStringLength(values.NewVariableValue(val.Value.(*styps.Variable)))
		p.Scope.Stack.Add(length)
		blk = p.Scope.Sprite.Sprite.NewRepeat(values.NewBlockValue(length))

	case types.ARRAY:
		length := p.Scope.Sprite.Sprite.NewLengthOfList(val.Value.(ArrayValue).Val)
		p.Scope.Stack.Add(length)
		blk = p.Scope.Sprite.Sprite.NewRepeat(values.NewBlockValue(length))

	case types.MAP:
		length := p.Scope.Sprite.Sprite.NewLengthOfList(val.Value.(MapValue).Key)
		p.Scope.Stack.Add(length)
		blk = p.Scope.Sprite.Sprite.NewRepeat(values.NewBlockValue(length))
	}
	p.Scope.Stack.Add(blk.(blocks.Block))

	// New scope
	oldScope := *p.Scope
	newVars := make(map[string]*Variable, len(oldScope.Vars))
	for k, v := range oldScope.Vars {
		newVars[k] = v
	}
	if keyDeclName != nil {
		newVars[*keyDeclName] = &Variable{
			Name:  *keyDeclName,
			Type:  keyDeclType,
			Value: keyDecl,
		}
	}
	if valDecl != nil {
		newVars[*valDeclName] = &Variable{
			Name:  *valDeclName,
			Type:  valDeclType,
			Value: valDecl,
		}
	}
	p.Scope = &Scope{
		Sprite: oldScope.Sprite,
		Fn:     oldScope.Fn,
		Vars:   newVars,
		Stack:  blk,
	}

	// Get value for valDecl
	switch val.Type.BasicType() {
	case types.STRING:
		if valDeclName != nil {
			ch := p.Scope.Sprite.Sprite.NewStringIndex(values.NewVariableValue(val.Value.(*styps.Variable)), values.NewVariableValue(iterV))
			p.Scope.Stack.Add(ch)
			blk := p.Scope.Sprite.Sprite.NewSetVariable(valDecl, values.NewBlockValue(ch))
			p.Scope.Stack.Add(blk)
		}

	case types.ARRAY:
		if valDeclName != nil {
			item := p.Scope.Sprite.Sprite.NewItemOfList(val.Value.(ArrayValue).Val, values.NewVariableValue(iterV))
			p.Scope.Stack.Add(item)
			blk := p.Scope.Sprite.Sprite.NewSetVariable(valDecl, values.NewBlockValue(item))
			p.Scope.Stack.Add(blk)
		}

	case types.MAP:
		if keyDeclName != nil {
			keyItem := p.Scope.Sprite.Sprite.NewItemOfList(val.Value.(MapValue).Key, values.NewVariableValue(iterV))
			p.Scope.Stack.Add(keyItem)
			blk := p.Scope.Sprite.Sprite.NewSetVariable(keyDecl, values.NewBlockValue(keyItem))
			p.Scope.Stack.Add(blk)
		}
		if valDeclName != nil {
			valItem := p.Scope.Sprite.Sprite.NewItemOfList(val.Value.(MapValue).Val, values.NewVariableValue(iterV))
			p.Scope.Stack.Add(valItem)
			blk := p.Scope.Sprite.Sprite.NewSetVariable(valDecl, values.NewBlockValue(valItem))
			p.Scope.Stack.Add(blk)
		}
	}

	// Add body
	for _, stmt := range stmt.Body.List {
		err := p.AddStmt(stmt)
		if err != nil {
			return err
		}
	}

	// Add increment
	change := p.Scope.Sprite.Sprite.NewChangeVariable(iterV, values.NewIntValue(1))
	p.Scope.Stack.Add(change)

	// Go back to old scope
	p.Scope = &oldScope
	return nil
}
