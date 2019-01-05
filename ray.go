// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

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
}
