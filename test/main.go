package main

import . "github.com/Nv7-Github/scratchy/scratch"

var a float64

type Sprite struct {
	value map[string]float64
	v     []float64
}

/*func (s *Sprite) add(a, b float64) float64 {
	return a + b
}*/

func (s *Sprite) addonetwo() float64 {
	return 1 + 2
}

func (s *Sprite) main() {
	a = 1
	s.v[1] = 1 // Note: Arrays start with 1 in scratch
	s.value["hi"] = 1.2

	/*a = s.add(1, 2)
	s.value["a"] = a*/
}
