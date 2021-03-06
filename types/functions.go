package types

import (
	"github.com/Nv7-Github/scratch/blocks"
	"github.com/Nv7-Github/scratch/types"
)

type Value struct {
	Value types.Value
	Type  Type
}

type Stack interface {
	Add(blk blocks.Block)
}
