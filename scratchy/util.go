package scratchy

import "strconv"

func (p *Program) GetVarName(name string, global bool) string {
	var dat map[string]*Variable
	if global {
		dat = p.GlobalVariables
	} else {
		dat = p.CurrSprite.Variables
	}

	_, exists := dat[name]
	if exists {
		i := 0
		var nameV string
		for exists {
			nameV = name + strconv.Itoa(i)
			_, exists = dat[nameV]
			i++
		}
		return nameV
	}
	return name
}
