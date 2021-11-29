package main

import . "github.com/Nv7-Github/scratchy/scratch"

type Sprite struct {
	console []string
}

func add(a, b float64) float64 {
	return a + b
}

func (s *Sprite) main() {
	v := add(1, 2)
	s.console = append(s.console, "v: "+NumberToString(v))
}
