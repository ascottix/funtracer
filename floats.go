// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
	"math/rand"
)

const Pi = math.Pi

type FloatGenerator func() float64

// Because of floating point precision, we don't test for strict equality but
// rather assume that two float numbers a, b are equal if abs(a-b) < Epsilon
const Epsilon = 0.00001

func FloatEqual(a, b float64) bool {
	return math.Abs(a-b) < Epsilon
}

func SliceFloatEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if math.Abs(v-b[i]) >= Epsilon {
			return false
		}
	}

	return true
}

func Min3(a, b, c float64) float64 {
	return math.Min(a, math.Min(b, c))
}

func Max3(a, b, c float64) float64 {
	return math.Max(a, math.Max(b, c))
}

func DegToRad(f float64) float64 {
	return f * 0.0174533
}

// Utilities
func Square(f float64) float64 {
	return f * f
}

func NewRandomGenerator(seed int64) FloatGenerator {
	rand := rand.New(rand.NewSource(seed))

	return rand.Float64
}
