package main

import (
	"github.com/marni/goigc"
)

//lengthCalc calculates track length based on example found here:
//https://github.com/marni/goigc/blob/master/doc_test.go
func lengthCalc(track igc.Track) float64 {
	totalDistance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}
	return totalDistance
}
