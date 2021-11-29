package main

import . "github.com/Nv7-Github/scratchy/scratch"

type Sprite struct {
	vals    map[string]float64
	console []string
}

func add(a, b float64) float64 {
	return a + b
}

func (s *Sprite) main() {
	s.vals["a"] = 1
	s.vals["b"] = 2
	s.vals["c"] = 3

	for k, v := range s.vals {
		s.console = append(s.console, "key: "+k)
		s.console = append(s.console, "val: "+NumberToString(v))
	}
}
