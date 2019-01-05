// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
)

func Identity() Matrix {
	return NewMatrix(4, 4,
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1)
}

func Translation(x, y, z float64) Matrix {
	return NewMatrix(4, 4,
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1)
}

func Scaling(s ...float64) Matrix {
	x := 1.0

	if len(s) > 0 {
		x = s[0]
	}

	y := x
	z := x

	if len(s) > 1 {
		y = s[1]
	}

	if len(s) > 2 {
		z = s[2]
	}

	return NewMatrix(4, 4,
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1)
}

func RotationX(a float64) Matrix {
	return NewMatrix(4, 4,
		1, 0, 0, 0,
		0, math.Cos(a), -math.Sin(a), 0,
		0, math.Sin(a), math.Cos(a), 0,
		0, 0, 0, 1)
}

func RotationY(a float64) Matrix {
	return NewMatrix(4, 4,
		math.Cos(a), 0, math.Sin(a), 0,
		0, 1, 0, 0,
		-math.Sin(a), 0, math.Cos(a), 0,
		0, 0, 0, 1)
}

func RotationZ(a float64) Matrix {
	return NewMatrix(4, 4,
		math.Cos(a), -math.Sin(a), 0, 0,
		math.Sin(a), math.Cos(a), 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1)
}

func Shearing(xy, xz, yx, yz, zx, zy float64) Matrix {
	return NewMatrix(4, 4,
		1, xy, xz, 0,
		yx, 1, yz, 0,
		zx, zy, 1, 0,
		0, 0, 0, 1)
}

// EyeViewpoint returns a world transformation that corresponds to the world being
// viewed from an eye placed at from and looking at to, with up being the up direction
func EyeViewpoint(from, to Tuple, up Tuple) Matrix {
	forward := to.Sub(from).Normalize()
	left := forward.CrossProduct(up.Normalize())
	up = left.CrossProduct(forward) // Get a "math-perfect" value for up (original parameter is ok even when approximated, this fixes that)
	orientation := NewMatrix(4, 4,
		left.X, left.Y, left.Z, 0,
		up.X, up.Y, up.Z, 0,
		-forward.X, -forward.Y, -forward.Z, 0,
		0, 0, 0, 1)

	return orientation.Mul(Translation(-from.X, -from.Y, -from.Z))
}
