package main

import . "github.com/Nv7-Github/scratchy/scratch"

var a float64

type Sprite struct {
	value map[string]float64
}

/*func (s *Sprite) add(a, b float64) float64 {
	return a + b
}*/

func (s *Sprite) main() {
	SayFor("Hello, World!", 1)
	SayFor("1 + 2 is", 1)
	Say(NumberToString(1 + 2))
	/*a = s.add(1, 2)
	s.value["a"] = a*/
}
