package functions

import (
	"github.com/Nv7-Github/scratch/sprites"
	"github.com/Nv7-Github/scratchy/types"
)

func init() {
	functions["Wait"] = Function{
		ParamTypes: []types.Type{types.NUMBER},
		ParamNames: []string{"time"},
		ReturnType: nil,
		Function: func(sprite *sprites.Sprite, stack types.Stack, params []*types.Value) (*types.Value, error) {
			blk := sprite.NewWait(params[0].Value)
			stack.Add(blk)
			return nil, nil
		},
	}
}
