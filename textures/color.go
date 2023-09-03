// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package textures

import (
	. "ascottix/funtracer/maths"
)

type Color struct {
	R, G, B float64
}

var Black = RGB(0, 0, 0)
var White = RGB(1, 1, 1)

func RGB(r, g, b float64) Color {
	return Color{r, g, b}
}

func Gray(v float64) Color {
	return RGB(v, v, v)
}

// Implement the image.Color interface
func (c Color) RGBA() (r, g, b, a uint32) {
	conv := func(c float64) uint32 {
		if c > 1 {
			c = 1
		}

		return uint32(c * 65535.99)
	}

	a = 0xFFFF
	r = conv(c.R)
	g = conv(c.G)
	b = conv(c.B)

	return
}

func (c Color) IsBlack() int32 {
	if c.R == 0 && c.G == 0 && c.B == 0 {
		return 1
	}
	return 0
}

func (c Color) Equals(d Color) bool {
	return FloatEqual(c.R, d.R) && FloatEqual(c.G, d.G) && FloatEqual(c.B, d.B)
}

func (c Color) Add(d Color) Color {
	return Color{c.R + d.R, c.G + d.G, c.B + d.B}
}

func (c Color) Sub(d Color) Color {
	return Color{c.R - d.R, c.G - d.G, c.B - d.B}
}

func (c Color) Mul(a float64) Color {
	return Color{a * c.R, a * c.G, a * c.B}
}

func (c Color) Erp(erp Interpolator) Color {
	return Color{erp(c.R), erp(c.G), erp(c.B)}
}

// Blend colors together by multiplying their components (Hadamard product)
func (c Color) Blend(d Color) Color {
	return Color{c.R * d.R, c.G * d.G, c.B * d.B}
}

// Blend colors together by multiplying their components (Hadamard product)
func (c Color) Gradient(d Color, t float64) Color {
	return Color{c.R + t*(d.R-c.R), c.G + t*(d.G-c.G), c.B + t*(d.B-c.B)}
}
