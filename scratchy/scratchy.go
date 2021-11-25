package scratchy

import (
	"fmt"
	"go/ast"
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
	Code    *ast.BlockStmt
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
	ScratchFunction *blocks.Function
	ReturnVal       *types.Variable
	SpriteName      string
}

type Scope struct {
	Sprite *Sprite
	Fn     *Function
	Stack  blocks.Stack
	Vars   map[string]*Variable
}

type Program struct {
	GlobalFunctions map[string]*GlobalFunction
	GlobalVariables map[string]*Variable
	Sprites         map[string]*Sprite

	Fset  *token.FileSet
	Scope *Scope
}

func (p *Program) NewError(pos token.Pos, format string, args ...interface{}) error {
	return fmt.Errorf("%s: "+format, append([]interface{}{p.Fset.Position(pos).String()}, args...)...)
}

func (p *Program) NewErrorFromError(pos token.Pos, err error) error {
	return p.NewError(pos, err.Error())
}

func newProgram(fset *token.FileSet) *Program {
	return &Program{
		GlobalFunctions: make(map[string]*GlobalFunction),
		GlobalVariables: make(map[string]*Variable),
		Sprites:         make(map[string]*Sprite),
		Scope:           &Scope{},

		Fset: fset,
	}
}

func (p *Program) CalcArgs(args []ast.Expr) ([]*styps.Value, error) {
	out := make([]*styps.Value, len(args))
	var err error
	for i, arg := range args {
		out[i], err = p.AddExpr(arg)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
