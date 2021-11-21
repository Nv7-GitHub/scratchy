package scratchy

import (
	"fmt"
	"go/ast"
)

func (p *Program) ConvType(typ ast.Expr) (Type, error) {
	switch t := typ.(type) {
	case *ast.Ident:
		name := t.Name
		bTyp, ok := basicTypeConv[name]
		if !ok {
			return nil, p.NewError(typ.Pos(), "unknown type: %s", name)
		}
		return bTyp, nil

	case *ast.ArrayType:
		valTyp := t.Elt.(*ast.Ident).Name
		valType, ok := basicTypeConv[valTyp]
		if !ok {
			return nil, p.NewError(typ.Pos(), "unknown array value type: %s", valTyp)
		}
		return NewArrayType(valType), nil

	case *ast.MapType:
		keyTyp := t.Key.(*ast.Ident).Name
		valTyp := t.Value.(*ast.Ident).Name
		keyType, ok := basicTypeConv[keyTyp]
		if !ok {
			return nil, p.NewError(typ.Pos(), "unknown map key type: %s", keyTyp)
		}
		valType, ok := basicTypeConv[valTyp]
		if !ok {
			return nil, p.NewError(typ.Pos(), "unknown map value type: %s", valTyp)
		}
		return NewMapType(keyType, valType), nil

	default:
		return nil, p.NewError(typ.Pos(), "unknown type: %v", typ)
	}
}

var basicTypeConv = map[string]BasicType{
	"bool":    BOOL,
	"float64": NUMBER,
	"string":  STRING,
}

type Type interface {
	fmt.Stringer

	BasicType() BasicType
	Equal(Type) bool
}

type BasicType int

const (
	NUMBER BasicType = iota
	STRING
	BOOL
	ARRAY
	MAP
)

var basicTypeNames = map[BasicType]string{
	NUMBER: "number",
	STRING: "string",
	BOOL:   "bool",
}

func (b BasicType) BasicType() BasicType {
	return b
}

func (b BasicType) Equal(t Type) bool {
	return t.BasicType() == b
}

func (b BasicType) String() string {
	return basicTypeNames[b]
}

type ArrayType struct {
	ValueType BasicType
}

func NewArrayType(ValType BasicType) *ArrayType {
	return &ArrayType{ValType}
}

func (a *ArrayType) BasicType() BasicType {
	return ARRAY
}

func (a *ArrayType) Equal(t Type) bool {
	if t.BasicType() != ARRAY {
		return false
	}
	return t.(*ArrayType).ValueType == a.ValueType
}

func (a *ArrayType) String() string {
	return fmt.Sprintf("array<%s>", a.ValueType.String())
}

type MapType struct {
	KeyType BasicType
	ValType BasicType
}

func NewMapType(keyType, valType BasicType) *MapType {
	return &MapType{keyType, valType}
}

func (m *MapType) BasicType() BasicType {
	return MAP
}

func (m *MapType) Equal(t Type) bool {
	if t.BasicType() != MAP {
		return false
	}
	return t.(*MapType).KeyType == m.KeyType && t.(*MapType).ValType == m.ValType
}

func (m *MapType) String() string {
	return fmt.Sprintf("map<%s,%s>", m.KeyType.String(), m.ValType.String())
}
