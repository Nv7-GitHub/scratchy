package scratchy

import "strconv"

func (p *Program) GetVarName(name string) string {
	_, exists := p.CurrSprite.Variables[name]
	if exists {
		i := 0
		var nameV string
		for exists {
			nameV = name + strconv.Itoa(i)
			_, exists = p.CurrSprite.Variables[nameV]
			i++
		}
		return nameV
	}
	return name
}
