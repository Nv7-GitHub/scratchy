package scratchy

import (
	"fmt"
	"go/token"

	"github.com/Nv7-Github/scratch/blocks"
	"github.com/Nv7-Github/scratch/sprites"
	"github.com/Nv7-Github/scratch/types"

	styps "github.com/Nv7-Github/scratchy/types"
)

type Param struct {
	Name string
	Type styps.Type
}

type GlobalFunction struct {
	Name    string
	Params  []Param
	RetType styps.Type
}

type Variable struct {
	Name  string
	Type  styps.Type
	Value interface{}
}

type MapValue struct {
	Key *types.List
	Val *types.List
}

type ArrayValue struct {
	Val *types.List
}

type Sprite struct {
	Name      string
	Functions map[string]*Function
	Variables map[string]*Variable
	Sprite    *sprites.Sprite
}

type Function struct {
	GlobalFunction
	ScratchFuntion *blocks.Function
}

type Program struct {
	GlobalFunctions map[string]*GlobalFunction
	GlobalVariables map[string]*Variable
	Sprites         map[string]*Sprite

	Fset *token.FileSet

	CurrSprite *Sprite
	CurrStack  blocks.Stack
}

func (p *Program) NewError(pos token.Pos, format string, args ...interface{}) error {
	return fmt.Errorf("%s: "+format, append([]interface{}{p.Fset.Position(pos).String()}, args...)...)
}

func newProgram() *Program {
	return &Program{
		GlobalFunctions: make(map[string]*GlobalFunction),
		GlobalVariables: make(map[string]*Variable),
		Sprites:         make(map[string]*Sprite),
	}
}
