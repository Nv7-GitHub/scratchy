package scratchy

import (
	"go/ast"

	"github.com/Nv7-Github/scratch/blocks"
)

func (p *Program) AddIf(stmt *ast.IfStmt) error {
	cond, err := p.AddExpr(stmt.Cond)
	if err != nil {
		return err
	}

	oldScope := *p.Scope

	// New scope
	var blk blocks.Stack
	if stmt.Else != nil {
		blk = p.Scope.Sprite.Sprite.NewIfElse(cond.Value)
		oldScope.Stack.Add(blk.(blocks.Block))
	} else {
		blk = p.Scope.Sprite.Sprite.NewIf(cond.Value)
		oldScope.Stack.Add(blk.(blocks.Block))
	}

	newVars := make(map[string]*Variable, len(oldScope.Vars))
	for k, v := range oldScope.Vars {
		newVars[k] = v
	}
	p.Scope = &Scope{
		Sprite: oldScope.Sprite,
		Fn:     oldScope.Fn,
		Vars:   newVars,
		Stack:  blk,
	}

	// Add body
	for _, stmt := range stmt.Body.List {
		err := p.AddStmt(stmt)
		if err != nil {
			return err
		}
	}

	// Add else if its there
	if stmt.Else != nil {
		// New scope again
		newVars := make(map[string]*Variable, len(oldScope.Vars))
		for k, v := range oldScope.Vars {
			newVars[k] = v
		}
		p.Scope = &Scope{
			Sprite: oldScope.Sprite,
			Fn:     oldScope.Fn,
			Vars:   newVars,
			Stack:  blk.(*blocks.IfElseBlock).Else,
		}

		// Add body
		for _, stmt := range stmt.Else.(*ast.BlockStmt).List {
			err := p.AddStmt(stmt)
			if err != nil {
				return err
			}
		}
	}

	// Old scope
	p.Scope = &oldScope
	return nil
}
