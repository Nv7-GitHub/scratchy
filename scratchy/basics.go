package scratchy

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/Nv7-Github/scratch/values"
	"github.com/Nv7-Github/scratchy/types"
)

func (p *Program) AddConst(expr *ast.BasicLit) (*types.Value, error) {
	switch expr.Kind {
	case token.INT:
		v, err := strconv.Atoi(expr.Value)
		if err != nil {
			return nil, p.NewErrorFromError(expr.ValuePos, err)
		}
		return &types.Value{
			Value: values.NewIntValue(v),
			Type:  types.NUMBER,
		}, nil

	case token.FLOAT:
		v, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			return nil, p.NewErrorFromError(expr.ValuePos, err)
		}
		return &types.Value{
			Value: values.NewFloatValue(v),
			Type:  types.NUMBER,
		}, nil

	case token.STRING:
		return &types.Value{
			Value: values.NewStringValue(expr.Value[1 : len(expr.Value)-1]), // remove quotes
			Type:  types.STRING,
		}, nil

	default:
		return nil, p.NewError(expr.ValuePos, "unsupported const type: %s", expr.Kind.String())
	}
}
