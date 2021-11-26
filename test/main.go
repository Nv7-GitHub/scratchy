package main

import . "github.com/Nv7-Github/scratchy/scratch"

var a float64

type Sprite struct {
	value map[string]float64
}

func add(a, b float64) float64 {
	return a + b
}

func (s *Sprite) main() {
	a = add(1, 2)
	s.value["a"] = a
	Say(NumberToString(s.value["a"]))
}
