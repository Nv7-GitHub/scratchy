package scratchy

import (
	"go/ast"

	"github.com/Nv7-Github/scratch/blocks"
	styps "github.com/Nv7-Github/scratch/types"
	"github.com/Nv7-Github/scratch/values"
	"github.com/Nv7-Github/scratchy/functions"
	"github.com/Nv7-Github/scratchy/types"
)

func (p *Program) AddFuncCall(expr *ast.CallExpr) (*types.Value, error) {
	// Calc args
	args, err := p.CalcArgs(expr.Args)
	if err != nil {
		return nil, err
	}

	name, ok := expr.Fun.(*ast.Ident)
	if !ok {
		spriteName := expr.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Name
		funcName := expr.Fun.(*ast.SelectorExpr).Sel.Name
		if spriteName != p.Scope.Fn.SpriteName {
			return nil, p.NewError(expr.Pos(), "unknown sprite: %s", spriteName)
		}
		fn, exists := p.Scope.Sprite.Functions[funcName]
		if !exists {
			return nil, p.NewError(expr.Pos(), "unknown sprite function: %s", funcName)
		}

		// Check args
		if len(args) != len(fn.Params) {
			return nil, p.NewError(expr.Pos(), "expected %d arguments, got %d", len(fn.Params), len(args))
		}
		for i, arg := range args {
			if !arg.Type.Equal(fn.Params[i].Type) {
				return nil, p.NewError(expr.Pos(), "expected argument %d to %s to be type %s, got type %s", i, fn.Name, fn.Params[i].Type.String(), arg.Type.String())
			}
		}

		// Call
		pars := make([]styps.Value, len(args))
		for i, arg := range args {
			pars[i] = arg.Value
		}
		call, err := p.Scope.Sprite.Sprite.NewFunctionCall(fn.ScratchFunction, pars...)
		if err != nil {
			return nil, err
		}
		p.Scope.Stack.Add(call)
		if fn.RetType != nil {
			return &types.Value{
				Type:  fn.RetType,
				Value: values.NewVariableValue(fn.ReturnVal),
			}, nil
		}
		return nil, nil
	}
	_, exists := p.GlobalFunctions[name.Name]
	if exists {
		return nil, p.NewError(expr.Fun.Pos(), "calling global functions is currently unsupported")
	}

	// builtin func
	res, err := functions.Call(p.Scope.Sprite.Sprite, p.Scope.Stack, name.Name, args)
	if err != nil {
		return nil, p.NewError(expr.Pos(), err.Error())
	}
	return res, nil
}

func (p *Program) AddReturn(stmt *ast.ReturnStmt) error {
	if p.Scope.Fn.RetType != nil {
		// Return value
		v, err := p.AddExpr(stmt.Results[0])
		if err != nil {
			return err
		}
		if !v.Type.Equal(p.Scope.Fn.RetType) {
			return p.NewError(stmt.Results[0].Pos(), "expected return type %s, got type %s", p.Scope.Fn.RetType.String(), v.Type.String())
		}

		blk := p.Scope.Sprite.Sprite.NewSetVariable(p.Scope.Fn.ReturnVal, v.Value)
		p.Scope.Stack.Add(blk)
	}

	// Stop
	blk := p.Scope.Sprite.Sprite.NewStop(blocks.StopOptionThisScript)
	p.Scope.Stack.Add(blk)
	return nil
}
