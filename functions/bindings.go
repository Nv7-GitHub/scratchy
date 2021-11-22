package functions

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/scratchy/types"
)

func ConvertType(typ types.Type) string {
	switch t := typ.(type) {
	case types.BasicType:
		switch t.BasicType() {
		case types.NUMBER:
			return "float64"

		case types.STRING:
			return "string"

		case types.BOOL:
			return "bool"
		}

	case *types.ArrayType:
		return "[]" + ConvertType(t.ValueType)

	case *types.MapType:
		return "map[" + ConvertType(t.KeyType) + "]" + ConvertType(t.ValType)
	}
	return "unknown"
}

func GetZero(typ types.Type) string {
	switch t := typ.(type) {
	case types.BasicType:
		switch t {
		case types.NUMBER:
			return "0"

		case types.STRING:
			return "\"\""

		case types.BOOL:
			return "false"
		}

	case *types.ArrayType:
		return fmt.Sprintf("make([]%s, 0)", ConvertType(t.ValueType))

	case *types.MapType:
		return fmt.Sprintf("make(map[%s]%s)", ConvertType(t.KeyType), ConvertType(t.ValType))
	}

	return "unknown"
}

func GenBindings() string {
	out := &strings.Builder{}
	out.WriteString("package scratch\n\n")

	for name, fn := range functions {
		if name[0] != strings.ToUpper(string(name[0]))[0] {
			continue
		}

		out.WriteString("func ")
		out.WriteString(name)
		out.WriteString("(")
		for i, arg := range fn.ParamNames {
			out.WriteString(arg)
			out.WriteString(" ")
			out.WriteString(ConvertType(fn.ParamTypes[i]))

			if i != len(fn.ParamNames)-1 {
				out.WriteString(", ")
			}
		}
		out.WriteString(")")
		if fn.ReturnType != nil {
			out.WriteString(" ")
			out.WriteString(ConvertType(fn.ReturnType))
			out.WriteString(" {\n\t")

			// return
			out.WriteString("return ")
			out.WriteString(GetZero(fn.ReturnType))
			out.WriteString("\n}\n\n")
		} else {
			out.WriteString(" {}\n\n")
		}
	}

	return out.String()
}
