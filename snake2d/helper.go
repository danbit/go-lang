package main

import (
	"math"
	"math/rand"
	"strconv"
)

func formatInt32(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}

func random(min, max int32) float64 {
	fmin := float64(min)
	fmax := float64(max)
	return math.Round((rand.Float64()*(fmax-fmin)+fmin)/float64(CellSize)) * float64(CellSize)
}
