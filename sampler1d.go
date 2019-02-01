// Copyright (c) 2019 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

type Sampler1d interface {
	Reset()        // Reset the sampler to the start position
	Next() float64 // Advances to the next position
	Get() float64  // Get a sample value in the current position
}

// Stratified1d generates one sample in the middle of each cell
type Stratified1d struct {
	CellSize     float64
	CellHalfSize float64
	curSample    float64
}

// Jittered1d adds a random perturbation to a 1d-sampler
type Jittered1d struct {
	sampler Sampler1d
	rand    FloatGenerator
	size    float64 // Maximum range for the random jitter, actual value goes from -size/2 to +size/2
}

func NewStratified1d(strata int) *Stratified1d {
	ss := Stratified1d{
		CellSize: 1 / float64(strata),
	}

	ss.CellHalfSize = ss.CellSize / 2
	ss.Reset()

	return &ss
}

func (ss *Stratified1d) Reset() {
	ss.curSample = -ss.CellHalfSize
}

func (ss *Stratified1d) Next() float64 {
	ss.curSample += ss.CellSize

	return ss.curSample
}

func (ss *Stratified1d) Get() float64 {
	return ss.curSample
}

func NewJittered1d(sampler Sampler1d, size float64, rand FloatGenerator) *Jittered1d {
	ss := Jittered1d{
		sampler: sampler,
		rand:    rand,
		size:    size,
	}

	return &ss
}

func (ss *Jittered1d) Reset() {
	ss.sampler.Reset()
}

func (ss *Jittered1d) Next() float64 {
	ss.sampler.Next()

	return ss.Get()
}

func (ss *Jittered1d) Get() float64 {
	r := (ss.rand() * ss.size) - (ss.size / 2)

	return ss.sampler.Get() + r
}
