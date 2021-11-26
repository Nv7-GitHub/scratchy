package main

import . "github.com/Nv7-Github/scratchy/scratch"

type Sprite struct {
	val string
}

func add(a, b float64) float64 {
	return a + b
}

func (s *Sprite) main() {
	for i := 0.0; i < 100; i++ {
		SayFor(NumberToString(i), 0.05)
	}
}
