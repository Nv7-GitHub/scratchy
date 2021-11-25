package functions

import (
	"github.com/Nv7-Github/scratch/blocks"
	"github.com/Nv7-Github/scratch/sprites"
	"github.com/Nv7-Github/scratchy/types"
)

func init() {
	functions["NumberToString"] = Function{
		ParamTypes: []types.Type{types.NUMBER},
		ParamNames: []string{"value"},
		ReturnType: types.STRING,
		Function: func(sprite *sprites.Sprite, stack blocks.Stack, params []*types.Value) (*types.Value, error) {
			return &types.Value{
				Value: params[0].Value,
				Type:  types.STRING,
			}, nil
		},
	}

	functions["StringToNumber"] = Function{
		ParamTypes: []types.Type{types.STRING},
		ParamNames: []string{"value"},
		ReturnType: types.NUMBER,
		Function: func(sprite *sprites.Sprite, stack blocks.Stack, params []*types.Value) (*types.Value, error) {
			return &types.Value{
				Value: params[0].Value,
				Type:  types.NUMBER,
			}, nil
		},
	}

}
