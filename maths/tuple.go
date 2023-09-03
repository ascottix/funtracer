// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package maths

import (
	"math"
)

type Tuple struct {
	X, Y, Z, W float64
}

func Point(x, y, z float64) Tuple {
	return Tuple{x, y, z, 1}
}

func Vector(x, y, z float64) Tuple {
	return Tuple{x, y, z, 0}
}

func PointAtInfinity(sign int) Tuple {
	inf := math.Inf(sign)

	return Point(inf, inf, inf)
}

// CompByIdx allows to reference X, Y and Z by index, which is useful for example when building a BVH
func (t Tuple) CompByIdx(idx int) float64 {
	switch idx {
	case 0:
		return t.X
	case 1:
		return t.Y
	default:
		return t.Z
	}
}

func (t Tuple) IsPoint() bool {
	return t.W == 1
}

func (t Tuple) IsVector() bool {
	return t.W == 0
}

func (t Tuple) Equals(u Tuple) bool {
	return FloatEqual(t.X, u.X) && FloatEqual(t.Y, u.Y) && FloatEqual(t.Z, u.Z) && FloatEqual(t.W, u.W)
}

func (t Tuple) StrictEquals(u Tuple) bool {
	return (t.X == u.X) && (t.Y == u.Y) && (t.Z == u.Z) && (t.W == u.W)
}

func (t Tuple) ApproxEquals(u Tuple) bool {
	return t.Equals(u)
}

func (t Tuple) Add(u Tuple) Tuple {
	return Tuple{t.X + u.X, t.Y + u.Y, t.Z + u.Z, t.W + u.W}
}

func (t Tuple) Sub(u Tuple) Tuple {
	return Tuple{t.X - u.X, t.Y - u.Y, t.Z - u.Z, t.W - u.W}
}

func (t Tuple) Neg() Tuple {
	return Tuple{-t.X, -t.Y, -t.Z, -t.W}
}

func (t Tuple) Mul(a float64) Tuple {
	return Tuple{a * t.X, a * t.Y, a * t.Z, a * t.W}
}

func (t Tuple) Div(a float64) Tuple {
	return Tuple{t.X / a, t.Y / a, t.Z / a, t.W / a}
}

func (t Tuple) Length() float64 {
	// Should be used only on vectors
	return math.Sqrt(t.X*t.X + t.Y*t.Y + t.Z*t.Z)
}

func (t Tuple) Normalize() Tuple {
	// Should be used only on vectors
	return t.Div(t.Length())
}

func (t Tuple) DotProduct(u Tuple) float64 {
	// Should be used only on vectors
	return t.X*u.X + t.Y*u.Y + t.Z*u.Z
}

func (t Tuple) CrossProduct(u Tuple) Tuple {
	// Should be used only on vectors
	return Tuple{t.Y*u.Z - t.Z*u.Y, t.Z*u.X - t.X*u.Z, t.X*u.Y - t.Y*u.X, 0}
}

func (t Tuple) Reflect(n Tuple) Tuple {
	return t.Sub(n.Mul(2 * t.DotProduct(n)))
}
