package functions

import (
	"github.com/Nv7-Github/scratch/sprites"
	"github.com/Nv7-Github/scratchy/types"
)

type Function struct {
	ParamTypes []types.Type
	ParamNames []string
	ReturnType types.Type
	Function   func(sprite *sprites.Sprite, stack types.Stack, params []*types.Value) (*types.Value, error)
}

var functions = make(map[string]Function)
