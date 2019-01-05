// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

func JadePattern() Pattern {
	p := NewPerlinNoisePattern(RGB(0, 0.8, 0), RGB(0.0, 0.2, 0.0), ErpQuadratic, ErpQuadratic)
	p.SetTransform(Scaling(3, 0.3, 0.3))
	return p
}

func AmberPattern() Pattern {
	p := NewPerlinNoisePattern(RGB(1, 0.6, 0.2), RGB(1, 0.8, 0.4), ErpQuadratic, ErpQuadratic)
	p.SetTransform(Scaling(0.3, 0.3, 0.3))
	return p
}

func WhiteLinesPattern() Pattern {
	p := NewPerlinNoisePattern(RGB(0.9, 0.9, 0.9), RGB(0.0, 0.2, 0.9), ErpToPerlinRange, ErpBezier, ErpClip)
	p.SetTransform(Scaling(0.07))
	return p
}

func MatMatte(c Color) *Material {
	return NewMaterial().
		SetPattern(NewSolidColorPattern(c))
}

func MatGlass() *Material {
	return NewMaterial().
		SetReflective(0.05).
		SetRefractive(0.95, 1.5)
}

func MatColoredGlass(c Color) *Material {
	return MatGlass().
		SetPattern(NewSolidColorPattern(c))
}
