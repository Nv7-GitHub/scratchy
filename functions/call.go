package functions

import (
	"fmt"

	"github.com/Nv7-Github/scratch/sprites"
	"github.com/Nv7-Github/scratchy/types"
)

func Call(sprite *sprites.Sprite, stack types.Stack, fnName string, params []*types.Value) (*types.Value, error) {
	fn, exists := functions[fnName]
	if !exists {
		return nil, fmt.Errorf("undeclared function: %s", fnName)
	}

	// Check par types
	for i, par := range params {
		if !par.Type.Equal(fn.ParamTypes[i]) {
			return nil, fmt.Errorf("invalid parameter type for %s: expected %s, got %s", fnName, fn.ParamTypes[i].String(), par.Type.String())
		}
	}

	return fn.Function(sprite, stack, params)
}
