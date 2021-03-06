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

	case *ast.ReturnStmt:
		return p.AddReturn(s)

	case *ast.AssignStmt:
		return p.AddAssignStmt(s)

	case *ast.IfStmt:
		return p.AddIf(s)

	case *ast.IncDecStmt:
		return p.IncDecStmt(s)

	case *ast.ForStmt:
		return p.AddLoop(s)

	case *ast.RangeStmt:
		return p.AddRange(s)

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

	case *ast.BinaryExpr:
		return p.AddMath(e)

	case *ast.Ident:
		return p.AddIdent(e)

	case *ast.SelectorExpr:
		return p.AddSelector(e)

	case *ast.IndexExpr:
		return p.AddIndex(e)

	case *ast.ParenExpr:
		return p.AddExpr(e.X)

	default:
		return nil, p.NewError(expr.Pos(), "unknown expression type: %T", e)
	}
}

func (p *Program) AddFuncCode(fn *Function) error {
	p.Scope.Fn = fn
	p.Scope.Stack = fn.ScratchFunction
	p.Scope.Vars = make(map[string]*Variable)

	// Add global params
	for _, par := range p.GlobalVariables {
		p.Scope.Vars[par.Name] = par
	}

	// Add function params
	for i, par := range fn.Params {
		p.Scope.Vars[par.Name] = &Variable{
			Name:  par.Name,
			Type:  par.Type,
			Value: fn.ScratchFunction.Parameters[i],
		}
	}

	// Add stmts
	for _, stmt := range fn.Code.List {
		err := p.AddStmt(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Program) CodePass() error {
	for _, sprite := range p.Sprites {
		p.Scope = &Scope{
			Sprite: sprite,
		}
		for _, fn := range sprite.Functions {
			err := p.AddFuncCode(fn)
			if err != nil {
				return err
			}
		}

		// Main function?
		fn, exists := sprite.Functions["main"]
		if exists && len(fn.Params) == 0 {
			onStart := p.Scope.Sprite.Sprite.NewWhenFlagClicked()
			call, err := p.Scope.Sprite.Sprite.NewFunctionCall(fn.ScratchFunction)
			if err != nil {
				return err
			}
			onStart.Add(call)
		}

		return nil
	}

	return nil
}
