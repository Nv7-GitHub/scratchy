package scratchy

import (
	"go/ast"
	"go/token"

	"github.com/Nv7-Github/scratch/blocks"
	"github.com/Nv7-Github/scratch/values"
	"github.com/Nv7-Github/scratchy/types"
)

var mathOps = map[token.Token]blocks.MathOperation{
	token.ADD: blocks.MathOperationAdd,
	token.SUB: blocks.MathOperationSubtract,
	token.MUL: blocks.MathOperationMultiply,
	token.QUO: blocks.MathOperationDivide,
}

var compareOps = map[token.Token]blocks.CompareOperand{
	token.EQL: blocks.CompareEqual,
	token.LSS: blocks.CompareLessThan,
	token.GTR: blocks.CompareGreaterThan,
}

func (p *Program) AddMath(expr *ast.BinaryExpr) (*types.Value, error) {
	v1, err := p.AddExpr(expr.X)
	if err != nil {
		return nil, err
	}
	v2, err := p.AddExpr(expr.Y)
	if err != nil {
		return nil, err
	}
	// Math op?
	mathOp, exists := mathOps[expr.Op]
	if exists {
		if !v1.Type.Equal(types.NUMBER) {
			return nil, p.NewError(expr.X.Pos(), "math only accepts numbers")
		}
		if !v2.Type.Equal(types.NUMBER) {
			return nil, p.NewError(expr.Y.Pos(), "math only accepts numbers")
		}

		blk := p.Scope.Sprite.Sprite.NewMath(v1.Value, v2.Value, mathOp)
		p.Scope.Stack.Add(blk)
		return &types.Value{
			Type:  types.NUMBER,
			Value: values.NewBlockValue(blk),
		}, nil
	}

	// Compare op?
	boolOp, exists := compareOps[expr.Op]
	if exists {
		if !v1.Type.Equal(types.NUMBER) {
			return nil, p.NewError(expr.X.Pos(), "comparison operations only accept numbers")
		}
		if !v2.Type.Equal(types.NUMBER) {
			return nil, p.NewError(expr.Y.Pos(), "comparison operations only accept numbers")
		}

		blk := p.Scope.Sprite.Sprite.NewCompare(v1.Value, v2.Value, boolOp)
		p.Scope.Stack.Add(blk)
		return &types.Value{
			Type:  types.BOOL,
			Value: values.NewBlockValue(blk),
		}, nil
	}

	return nil, p.NewError(expr.OpPos, "unknown operation: %s", expr.Op.String())
}
