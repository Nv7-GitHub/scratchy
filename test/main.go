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
	SayFor(NumberToString(1+2), 1)

	SayFor("Is 1 equal to 1?", 1)
	SayFor(BoolToString(1 == 1), 1)
	/*a = s.add(1, 2)
	s.value["a"] = a*/
}
