package main

// Test scratchy program

import "github.com/Nv7-Github/scratchy/scratch"

var global int

type Sprite struct {
	vals      []float64
	numSorted int
	i         int
	tmp       float64 // buffer
}

//scratchy:main
func (s *Sprite) main() {
	global = 10

	scratch.Clear(s.vals)
	scratch.Append(s.vals, 3)
	scratch.Append(s.vals, 2)
	scratch.Append(s.vals, 2.5)
	scratch.Append(s.vals, 1)

	scratch.Say("It's unsorted!")
	scratch.Wait(1)
	scratch.Say("Sorting...")

	s.numSorted = 1
	for (s.numSorted < len(s.vals)) || (s.numSorted == len(s.vals)) {
		smallest := s.vals[s.numSorted]
		smallestInd := s.numSorted
		for s.i = s.numSorted; s.i < len(s.vals)+1; s.i++ {
			if s.vals[s.i] < smallest {
				smallest = s.vals[s.i]
				smallestInd = s.i
			}
		}
		s.tmp = s.vals[s.numSorted]
		s.vals[s.numSorted] = s.vals[smallestInd]
		s.vals[smallestInd] = s.tmp
		s.numSorted++
	}

	scratch.SayFor("It's sorted!", 1)
}
