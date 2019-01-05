// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
)

var permutation [256]int = [256]int{151, 160, 137, 91, 90, 15,
	131, 13, 201, 95, 96, 53, 194, 233, 7, 225, 140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23,
	190, 6, 148, 247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32, 57, 177, 33,
	88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175, 74, 165, 71, 134, 139, 48, 27, 166,
	77, 146, 158, 231, 83, 111, 229, 122, 60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244,
	102, 143, 54, 65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169, 200, 196,
	135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64, 52, 217, 226, 250, 124, 123,
	5, 202, 38, 147, 118, 126, 255, 82, 85, 212, 207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42,
	223, 183, 170, 213, 119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104, 218, 246, 97, 228,
	251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241, 81, 51, 145, 235, 249, 14, 239, 107,
	49, 192, 214, 31, 181, 199, 106, 157, 184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254,
	138, 236, 205, 93, 222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180}

// PerlinNoise implements the improved Ken Perlin noise as described here:
// https://mrl.nyu.edu/~perlin/noise/
type PerlinNoise struct {
	p [512]int
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func grad(hash int, x, y, z float64) float64 {
	// Convert lower 4 bits of hash code into 12 gradient directions
	var u, v float64

	h := hash & 15

	if h < 8 {
		u = x
	} else {
		u = y
	}

	if (h & 1) != 0 {
		u = -u
	}

	if h < 4 {
		v = y
	} else if h == 12 || h == 14 {
		v = x
	} else {
		v = z
	}

	if (h & 2) != 0 {
		v = -v
	}

	return u + v
}

func NewPerlinNoise() *PerlinNoise {
	pn := &PerlinNoise{}

	for i, v := range permutation {
		pn.p[i] = v
		pn.p[256+i] = v
	}

	return pn
}

func (pn *PerlinNoise) at(x, y, z float64) float64 {
	// Find unit cube that contais point
	ix := int(math.Floor(x)) & 0xFF
	iy := int(math.Floor(y)) & 0xFF
	iz := int(math.Floor(z)) & 0xFF

	// Find relative coordinates of point in cube
	x -= math.Floor(x)
	y -= math.Floor(y)
	z -= math.Floor(z)

	// Compute fade curves for each coordinate
	u, v, w := fade(x), fade(y), fade(z)

	// Hash coordinates of the 8 cube corners
	a := pn.p[ix] + iy
	b := pn.p[ix+1] + iy
	aa := pn.p[a] + iz
	ab := pn.p[a+1] + iz
	ba := pn.p[b] + iz
	bb := pn.p[b+1] + iz

	// Add blended results from 8 corners of cube
	return lerp(w, lerp(v, lerp(u, grad(pn.p[aa], x, y, z),
		grad(pn.p[ba], x-1, y, z)),
		lerp(u, grad(pn.p[ab], x, y-1, z),
			grad(pn.p[bb], x-1, y-1, z))),
		lerp(v, lerp(u, grad(pn.p[aa+1], x, y, z-1),
			grad(pn.p[ba+1], x-1, y, z-1)),
			lerp(u, grad(pn.p[ab+1], x, y-1, z-1),
				grad(pn.p[bb+1], x-1, y-1, z-1))))
}

var defaultPerlinNoise *PerlinNoise = NewPerlinNoise()

func Perlin(p Tuple) float64 {
	return defaultPerlinNoise.at(p.X, p.Y, p.Z)
}
