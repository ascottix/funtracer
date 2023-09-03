// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package textures

import (
	. "ascottix/funtracer/maths"
)

type Ray struct {
	Origin    Tuple
	Direction Tuple
}

func NewRay(p, v Tuple) Ray {
	r := Ray{p, v}

	return r
}

func (r Ray) Position(t float64) Tuple {
	return r.Origin.Add(r.Direction.Mul(t))
}

func (r Ray) Transform(m Matrix) Ray {
	return Ray{
		m.MulT(r.Origin),
		m.MulT(r.Direction),
	}

	// A bit faster but probably not really worth it
	// return Ray{
	// 	Tuple{
	// 		X: m.A[0]*r.Origin.X + m.A[1]*r.Origin.Y + m.A[2]*r.Origin.Z + m.A[3]*r.Origin.W,
	// 		Y: m.A[4]*r.Origin.X + m.A[5]*r.Origin.Y + m.A[6]*r.Origin.Z + m.A[7]*r.Origin.W,
	// 		Z: m.A[8]*r.Origin.X + m.A[9]*r.Origin.Y + m.A[10]*r.Origin.Z + m.A[11]*r.Origin.W,
	// 		W: r.Origin.W,
	// 	},
	// 	Tuple{
	// 		X: m.A[0]*r.Direction.X + m.A[1]*r.Direction.Y + m.A[2]*r.Direction.Z + m.A[3]*r.Direction.W,
	// 		Y: m.A[4]*r.Direction.X + m.A[5]*r.Direction.Y + m.A[6]*r.Direction.Z + m.A[7]*r.Direction.W,
	// 		Z: m.A[8]*r.Direction.X + m.A[9]*r.Direction.Y + m.A[10]*r.Direction.Z + m.A[11]*r.Direction.W,
	// 		W: r.Direction.W,
	// 	}}
}
