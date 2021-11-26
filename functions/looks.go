package functions

import (
	"github.com/Nv7-Github/scratch/sprites"
	"github.com/Nv7-Github/scratchy/types"
)

func init() {
	functions["Say"] = Function{
		ParamTypes: []types.Type{types.STRING},
		ParamNames: []string{"text"},
		ReturnType: nil,
		Function: func(sprite *sprites.Sprite, stack types.Stack, params []*types.Value) (*types.Value, error) {
			blk := sprite.NewSayBlock(params[0].Value)
			stack.Add(blk)
			return nil, nil
		},
	}

	functions["SayFor"] = Function{
		ParamTypes: []types.Type{types.STRING, types.NUMBER},
		ParamNames: []string{"text", "time"},
		ReturnType: nil,
		Function: func(sprite *sprites.Sprite, stack types.Stack, params []*types.Value) (*types.Value, error) {
			blk := sprite.NewSayForTimeBlock(params[0].Value, params[1].Value)
			stack.Add(blk)
			return nil, nil
		},
	}
}
