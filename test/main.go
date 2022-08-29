package main

// Test scratchy program

import . "github.com/Nv7-Github/scratchy/scratch"

var global int

type Sprite struct {
	vals      []float64
	numSorted int
	i         int
	tmp       float64 // buffer
}

func (s *Sprite) main() {
	global = 10

	Clear(s.vals)
	Append(s.vals, 3)
	Append(s.vals, 2)
	Append(s.vals, 2.5)
	Append(s.vals, 1)

	Say("It's unsorted!")
	Wait(1)
	Say("Sorting...")

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

	SayFor("It's sorted!", 1)
}
