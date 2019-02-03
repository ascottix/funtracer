// Copyright (c) 2019 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
)

type Sampler2d interface {
	Reset()
	Next() (float64, float64)
}

type Combined1d1d struct {
	cx       int // How many x samples before we need to increment y
	sx       int
	sy       int
	samplerX Sampler1d
	samplerY Sampler1d
}

// Returns samples in a 2D unit disk centered at the origin (0,0) 
type ConcentricSampleDisk struct {
	rand 	FloatGenerator
}

func NewCombined1d1d(sx, sy int, samplerx, samplery Sampler1d) *Combined1d1d {
	ss := Combined1d1d{
		sx:       sx,
		sy:       sy,
		samplerX: samplerx,
		samplerY: samplery,
	}

	ss.Reset()

	return &ss
}

func (ss *Combined1d1d) Reset() {
	ss.samplerY.Reset()
	ss.cx = 0
}

func (ss *Combined1d1d) Next() (float64, float64) {
	if ss.cx == 0 {
		ss.cx = ss.sx
		ss.samplerX.Reset()
		ss.samplerY.Next()
	}
	ss.cx--

	return ss.samplerX.Next(), ss.samplerY.Get()
}

func NewStratified2d(sx, sy int) *Combined1d1d {
	return NewCombined1d1d(sx, sy, NewStratified1d(sx), NewStratified1d(sy))
}

func NewJitteredStratified2d(sx, sy int, rand FloatGenerator) *Combined1d1d {
	return NewCombined1d1d(sx, sy,
		NewJittered1d(NewStratified1d(sx), 1/float64(sx), rand),
		NewJittered1d(NewStratified1d(sy), 1/float64(sy), rand))
}

func NewConcentricSampleDisk(rand FloatGenerator) *ConcentricSampleDisk {
	return &ConcentricSampleDisk{rand: rand}
}

func (s *ConcentricSampleDisk) Reset() {
}

func (s *ConcentricSampleDisk) Next() (float64, float64) {
	// Get two random samples in the [-1,+1] interval
	x := s.rand()*2 - 1
	y := s.rand()*2 - 1

	if x != 0 || y != 0 {
		// Apply concentric mapping to point
		var r, theta float64

		if math.Abs(x) > math.Abs(y) {
			r = x
			theta = (y / x) * Pi / 4
		} else {
			r = y
			theta = Pi / 2 - (x / y) * Pi / 4
		}

		x = r * math.Cos(theta)
		y = r * math.Sin(theta)
	}

    return x, y
}
