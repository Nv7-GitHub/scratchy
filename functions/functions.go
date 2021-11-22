package functions

import (
	styps "github.com/Nv7-Github/scratch/types"
	"github.com/Nv7-Github/scratchy/types"
)

type Function struct {
	ParamTypes []types.Type
	ReturnType types.Type
	Build      func(params []styps.Value) (styps.Value, error)
}
