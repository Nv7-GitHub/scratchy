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
	return fmt.Sprintf("map<%s, %s>", m.KeyType.String(), m.ValType.String())
}
