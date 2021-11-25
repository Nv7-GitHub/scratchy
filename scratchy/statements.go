package scratchy

import (
	"go/ast"

	"github.com/Nv7-Github/scratchy/types"
)

func (p *Program) AddStmt(stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		_, err := p.AddExpr(s.X)
		return err

	default:
		return p.NewError(stmt.Pos(), "unknown statement type: %T", s)
	}
}

func (p *Program) AddExpr(expr ast.Expr) (*types.Value, error) {
	switch e := expr.(type) {
	case *ast.CallExpr:
		return p.AddFuncCall(e)

	case *ast.BasicLit:
		return p.AddConst(e)

	default:
		return nil, p.NewError(expr.Pos(), "unknown expression type: %T", e)
	}
}

func (p *Program) CodePass() error {
	for _, sprite := range p.Sprites {
		p.CurrSprite = sprite
		for _, fn := range sprite.Functions {
			p.CurrStack = fn.ScratchFuntion
			for _, stmt := range fn.Code.List {
				err := p.AddStmt(stmt)
				if err != nil {
					return err
				}
			}
		}

		// Main function?
		fn, exists := sprite.Functions["main"]
		if exists && len(fn.Params) == 0 {
			onStart := p.CurrSprite.Sprite.NewWhenFlagClicked()
			call, err := p.CurrSprite.Sprite.NewFunctionCall(fn.ScratchFuntion)
			if err != nil {
				return err
			}
			onStart.Add(call)
		}
	}

	return nil
}