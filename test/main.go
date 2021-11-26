package main

import . "github.com/Nv7-Github/scratchy/scratch"

type Sprite struct {
	val string
}

func add(a, b float64) float64 {
	return a + b
}

func (s *Sprite) main() {
	s.val = BoolToString("a" == "A")
	SayFor(s.val, 1)

	s.val = "Hello, " + "World!"
	SayFor(s.val, 1)
}
