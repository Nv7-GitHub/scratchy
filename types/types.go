package types

import (
	"fmt"
)

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
	ANY
)

func (b BasicType) BasicType() BasicType {
	return b
}

func (b BasicType) Equal(t Type) bool {
	if b == ANY || t == ANY {
		return true
	}
	return t.BasicType() == b
}

func (b BasicType) String() string {
	return [...]string{"numnber", "string", "bool", "array", "any"}[b]
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
