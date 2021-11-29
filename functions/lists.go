package functions

import (
	"github.com/Nv7-Github/scratch/sprites"
	"github.com/Nv7-Github/scratch/values"
	"github.com/Nv7-Github/scratchy/types"
)

func init() {
	functions["append"] = Function{
		ParamTypes: []types.Type{types.ARRAY, types.ANY},
		ParamNames: []string{"value"},
		ReturnType: nil,
		Function: func(sprite *sprites.Sprite, stack types.Stack, params []*types.Value) (*types.Value, error) {
			blk := sprite.NewAddToList(params[0].Value.(*values.ListValue).List, params[1].Value)
			stack.Add(blk)
			return nil, nil
		},
	}
}
