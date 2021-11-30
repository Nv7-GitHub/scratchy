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

	functions["len"] = Function{
		ParamTypes: []types.Type{types.ARRAY},
		ParamNames: []string{"array"},
		ReturnType: types.NUMBER,
		Function: func(sprite *sprites.Sprite, stack types.Stack, params []*types.Value) (*types.Value, error) {
			length := sprite.NewLengthOfList(params[0].Value.(*values.ListValue).List)
			stack.Add(length)
			return &types.Value{
				Type:  types.NUMBER,
				Value: values.NewBlockValue(length),
			}, nil
		},
	}

	functions["ClearNumberArray"] = Function{
		ParamTypes: []types.Type{types.NewArrayType(types.NUMBER)},
		ParamNames: []string{"array"},
		ReturnType: nil,
		Function: func(sprite *sprites.Sprite, stack types.Stack, params []*types.Value) (*types.Value, error) {
			clearArray(stack, params[0].Value.(*values.ListValue), sprite)
			return nil, nil
		},
	}

	functions["ClearStringArray"] = Function{
		ParamTypes: []types.Type{types.NewArrayType(types.STRING)},
		ParamNames: []string{"array"},
		ReturnType: nil,
		Function: func(sprite *sprites.Sprite, stack types.Stack, params []*types.Value) (*types.Value, error) {
			clearArray(stack, params[0].Value.(*values.ListValue), sprite)
			return nil, nil
		},
	}

	functions["ClearBoolArray"] = Function{
		ParamTypes: []types.Type{types.NewArrayType(types.BOOL)},
		ParamNames: []string{"array"},
		ReturnType: nil,
		Function: func(sprite *sprites.Sprite, stack types.Stack, params []*types.Value) (*types.Value, error) {
			clearArray(stack, params[0].Value.(*values.ListValue), sprite)
			return nil, nil
		},
	}
}

func clearArray(stack types.Stack, array *values.ListValue, sprite *sprites.Sprite) {
	blk := sprite.NewDeleteAllFromList(array.List)
	stack.Add(blk)
}
