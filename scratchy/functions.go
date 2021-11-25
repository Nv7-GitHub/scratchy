package scratchy

import (
	"go/ast"

	"github.com/Nv7-Github/scratch/blocks"
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
		return nil, p.NewError(expr.Fun.Pos(), "calling sprite functions is currently unsupported")
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
