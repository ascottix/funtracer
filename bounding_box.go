// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
)

type Box struct {
	Min Tuple
	Max Tuple
}

func NewBox(pmin, pmax Tuple) Box {
	return Box{pmin, pmax}
}

func (b Box) Union(a Box) Box {
	return Box{
		Point(math.Min(a.Min.X, b.Min.X), math.Min(a.Min.Y, b.Min.Y), math.Min(a.Min.Z, b.Min.Z)),
		Point(math.Max(a.Max.X, b.Max.X), math.Max(a.Max.Y, b.Max.Y), math.Max(a.Max.Z, b.Max.Z)),
	}
}

func (b Box) Transform(m Matrix) Box {
	// Get the box vertices
	vs := []Tuple{
		b.Min,
		Point(b.Min.X, b.Min.Y, b.Max.Z),
		Point(b.Min.X, b.Max.Y, b.Min.Z),
		Point(b.Min.X, b.Max.Y, b.Max.Z),
		Point(b.Max.X, b.Min.Y, b.Min.Z),
		Point(b.Max.X, b.Min.Y, b.Max.Z),
		Point(b.Max.X, b.Max.Y, b.Min.Z),
		b.Max,
	}

	a := Box{
		Point(math.Inf(+1), math.Inf(+1), math.Inf(+1)),
		Point(math.Inf(-1), math.Inf(-1), math.Inf(-1)),
	}

	for _, p := range vs {
		p := m.MulT(p) // Transform point

		a = a.Union(Box{p, p})
	}

	return a
}

func (b Box) Intersects(ray Ray) bool {
	checkAxis := func(origin, direction, min, max float64) (tmin, tmax float64) {
		tminNumerator := min - origin
		tmaxNumerator := max - origin

		if math.Abs(direction) >= Epsilon {
			tmin = tminNumerator / direction
			tmax = tmaxNumerator / direction
		} else {
			tmin = tminNumerator * math.Inf(+1)
			tmax = tmaxNumerator * math.Inf(+1)
		}

		if tmin > tmax {
			tmin, tmax = tmax, tmin
		}

		return tmin, tmax
	}

	xtmin, xtmax := checkAxis(ray.Origin.X, ray.Direction.X, b.Min.X, b.Max.X)
	ytmin, ytmax := checkAxis(ray.Origin.Y, ray.Direction.Y, b.Min.Y, b.Max.Y)
	ztmin, ztmax := checkAxis(ray.Origin.Z, ray.Direction.Z, b.Min.Z, b.Max.Z)

	tmin := math.Max(xtmin, math.Max(ytmin, ztmin))
	tmax := math.Min(xtmax, math.Min(ytmax, ztmax))

	return tmin < tmax
}

func (b Box) ToCube() *Shape {
	s := NewCube()

	sx := (b.Max.X - b.Min.X) / 2
	sy := (b.Max.Y - b.Min.Y) / 2
	sz := (b.Max.Z - b.Min.Z) / 2

	s.SetTransform(Translation(b.Max.X-sx, b.Max.Y-sy, b.Max.Z-sz), Scaling(sx, sy, sz))

	return s
}
