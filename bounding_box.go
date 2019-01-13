// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
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

var _Debug = false

func (b Box) Intersects(ray Ray) bool {
	ray.Direction.X = 1 / ray.Direction.X
	ray.Direction.Y = 1 / ray.Direction.Y
	ray.Direction.Z = 1 / ray.Direction.Z

	return b.IntersectsInvDir(ray)
}

func (b Box) IntersectsInvDir(ray Ray) bool {
	checkAxis := func(origin, invdir, min, max float64) (tmin, tmax float64) {
		if invdir >= 0 {
			tmin = (min - origin) * invdir
			tmax = (max - origin) * invdir
		} else {
			tmax = (min - origin) * invdir
			tmin = (max - origin) * invdir
		}

		return tmin, tmax
	}

	xtmin, xtmax := checkAxis(ray.Origin.X, ray.Direction.X, b.Min.X, b.Max.X)
	ytmin, ytmax := checkAxis(ray.Origin.Y, ray.Direction.Y, b.Min.Y, b.Max.Y)
	ztmin, ztmax := checkAxis(ray.Origin.Z, ray.Direction.Z, b.Min.Z, b.Max.Z)

	tmin := math.Inf(-1)
	tmax := math.Inf(+1)

	// Note to self: do _not_ use something like math.Max(xtmin, math.Max(ytmin, ztmin)) because
	// it will fail completely when one of the operands is NaN, whereas the method below will
	// simply skip over the invalid value (we get a 0/0 NaN from checkAxis when direction=0 and min/max=origin)

	if xtmin > tmin {
		tmin = xtmin
	}
	if ytmin > tmin {
		tmin = ytmin
	}
	if ztmin > tmin {
		tmin = ztmin
	}

	if xtmax < tmax {
		tmax = xtmax
	}
	if ytmax < tmax {
		tmax = ytmax
	}
	if ztmax < tmax {
		tmax = ztmax
	}

	if _Debug {
		Debugln("xt=", xtmin, ",", xtmax, ", yt=", ytmin, ",", ytmax, "zt=", ztmin, ztmax)
		Debugln("t=", tmin, tmax)
	}

	return tmin <= tmax
}

func (b Box) ToCube() *Shape {
	s := NewCube()

	sx := (b.Max.X - b.Min.X) / 2
	sy := (b.Max.Y - b.Min.Y) / 2
	sz := (b.Max.Z - b.Min.Z) / 2

	s.SetTransform(Translation(b.Max.X-sx, b.Max.Y-sy, b.Max.Z-sz), Scaling(sx, sy, sz))

	return s
}

func (b Box) Diagonal() Tuple {
	return Vector(b.Max.X-b.Min.X, b.Max.Y-b.Min.Y, b.Max.Z-b.Min.Z)
}

func (b Box) SurfaceArea() float64 {
	d := b.Diagonal()

	return 2 * (d.X*d.Y + d.X*d.Z + d.Y*d.Z)
}

// Offset computes the offset of a point inside the bounding box,
// it is a vector with components in the [0,1] interval
func (b Box) Offset(p Tuple) Tuple {
	o := p.Sub(b.Min)

	if b.Max.X > b.Min.X {
		o.X /= (b.Max.X - b.Min.X)
	}

	if b.Max.Y > b.Min.Y {
		o.Y /= (b.Max.Y - b.Min.Y)
	}

	if b.Max.Z > b.Min.Z {
		o.Z /= (b.Max.Z - b.Min.Z)
	}

	return o
}

func (b Box) String() string {
	return fmt.Sprintf("Box (%.2f,%.2f,%.2f) - (%.2f, %.2f, %.2f)", b.Min.X, b.Min.Y, b.Min.Z, b.Max.X, b.Max.Y, b.Max.Z)
}
