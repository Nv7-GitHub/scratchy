package main

import . "github.com/Nv7-Github/scratchy/scratch"

type Sprite struct {
	vals []float64
}

func (s *Sprite) main() {
	ClearNumberArray(s.vals)
	s.vals = append(s.vals, 3)
	s.vals = append(s.vals, 2)
	s.vals = append(s.vals, 2.5)
	s.vals = append(s.vals, 1)

	Say("It's unsorted!")
	Wait(1)
	Say("Sorting...")

	numSorted := 1
	for (numSorted < len(s.vals)) || (numSorted == len(s.vals)) {
		smallest := s.vals[numSorted]
		smallestInd := numSorted
		for i := numSorted; i < len(s.vals)+1; i++ {
			if s.vals[i] < smallest {
				smallest = s.vals[i]
				smallestInd = i
			}
		}
		v := s.vals[numSorted]
		s.vals[numSorted] = s.vals[smallestInd]
		s.vals[smallestInd] = v
		numSorted++
	}

	SayFor("It's sorted!", 1)
}
