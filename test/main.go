package main

var a float64

type Sprite struct {
	value map[string]float64
}

func (s *Sprite) add(a, b float64) float64 {
	return a + b
}

func (s *Sprite) main() {
	a = s.add(1, 2)
	s.value["a"] = a
}
