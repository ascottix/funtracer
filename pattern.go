// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
)

type Texture interface {
	ApplyAtHit(ii *IntersectionInfo)
}

type Pattern interface {
	Transformable
	Texture
	PatternAt(point Tuple) Color
	ColorAt(object Patternable, point Tuple) Color
}

type LocalPatternAt func(point Tuple) Color

type BasicPattern struct {
	Transformer
	LocalPatternAt LocalPatternAt
}

type ProxyPattern struct {
	P Pattern
	O Patternable
}

func NewProxyPattern(object Patternable, pattern Pattern) (p Pattern) {
	if pattern != nil {
		p = &ProxyPattern{O: object, P: pattern}
	}

	return p
}

func (p *ProxyPattern) SetTransform(transforms ...Matrix) {
	p.P.SetTransform(transforms...)
}

func (p *ProxyPattern) Transform() Matrix {
	return p.P.Transform()
}

func (p *ProxyPattern) InverseTransform() Matrix {
	return p.P.InverseTransform()
}

func (p *ProxyPattern) PatternAt(point Tuple) Color {
	return p.P.PatternAt(point)
}

func (p *ProxyPattern) ColorAt(object Patternable, point Tuple) Color {
	return p.P.ColorAt(p.O, point)
}

func (p *ProxyPattern) ApplyAtHit(ii *IntersectionInfo) {
	ii.Mat.DiffuseColor = p.P.ColorAt(ii.O, ii.Point)
}

func NewBasicPattern(localPatternAt LocalPatternAt) Pattern {
	p := &BasicPattern{LocalPatternAt: localPatternAt}
	p.SetTransform()

	return p
}

func (p *BasicPattern) PatternAt(point Tuple) Color {
	return p.LocalPatternAt(point)
}

func (p *BasicPattern) ColorAt(object Patternable, point Tuple) Color {
	objectPoint := object.WorldToObject(point)   // Convert from world to object space
	patternPoint := p.Tinverse.MulT(objectPoint) // Convert from object to pattern space
	return p.LocalPatternAt(patternPoint)
}

func (p *BasicPattern) ApplyAtHit(ii *IntersectionInfo) {
	ii.Mat.DiffuseColor = p.ColorAt(ii.O, ii.Point)
}

func NewSolidColorPattern(c Color) Pattern {
	return NewBasicPattern(func(point Tuple) Color {
		return c
	})
}

func NewStripePattern(a, b Color) Pattern {
	c := [2]Color{a, b}

	return NewBasicPattern(func(point Tuple) Color {
		return c[uint32(math.Floor(point.X))%2]
	})
}

func NewCheckerPattern(a, b Color) Pattern {
	c := [2]Color{a, b}

	return NewBasicPattern(func(point Tuple) Color {
		return c[uint32(math.Floor(point.X)+math.Floor(point.Y)+math.Floor(point.Z))%2]
	})
}

func NewGradientPattern(a, b Color) Pattern {
	return NewBasicPattern(func(point Tuple) Color {
		return a.Add((b.Sub(a)).Mul(point.X - math.Floor(point.X)))
	})
}

func NewBlendedPattern(patterns ...Pattern) Pattern {
	return NewBasicPattern(func(point Tuple) Color {
		c := Black
		for _, p := range patterns {
			c = c.Add(p.PatternAt(point))
		}
		return c
	})
}

func NewPointPattern() Pattern { // Used mainly for test: returns the point coordinates as a color
	return NewBasicPattern(func(point Tuple) Color {
		return RGB(point.X, point.Y, point.Z)
	})
}

func NewPerlinNoisePattern(a, b Color, interpolators ...Interpolator) Pattern {
	pn := NewPerlinNoise()

	return NewBasicPattern(func(point Tuple) Color {
		noise := pn.at(point.X, point.Y, point.Z)
		noise = (noise + 1) / 2 // Perlin noise is in the [-1,+1] range, take it to [0,1]

		for _, interp := range interpolators {
			noise = interp(noise)
		}

		return a.Gradient(b, noise)
	})
}
