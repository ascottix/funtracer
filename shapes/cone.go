// Copyright (c) 2019 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package shapes

import (
	"math"

	. "ascottix/funtracer/maths"
	. "ascottix/funtracer/textures"
)

type Cone struct {
	MinY   float64
	MaxY   float64
	Capped bool
}

func NewCone(miny, maxy float64, capped bool) *Shape {
	return NewShape("cone", &Cone{miny, maxy, capped})
}

func NewInfiniteCone() *Shape {
	return NewCone(math.Inf(-1), math.Inf(+1), false)
}

func (p *Cone) Bounds() Box {
	maxr := math.Max(math.Abs(p.MinY), math.Abs(p.MaxY))

	return Box{Point(-maxr, p.MinY, -maxr), Point(+maxr, p.MaxY, +maxr)}
}

func (p *Cone) LocalIntersect(ray Ray) (xs []float64) {
	a := Square(ray.Direction.X) - Square(ray.Direction.Y) + Square(ray.Direction.Z)
	b := 2 * (ray.Origin.X*ray.Direction.X - ray.Origin.Y*ray.Direction.Y + ray.Origin.Z*ray.Direction.Z)
	c := Square(ray.Origin.X) - Square(ray.Origin.Y) + Square(ray.Origin.Z)

	// Check intersection with cone surface
	if math.Abs(a) > Epsilon {
		disc := b*b - 4*a*c

		if disc >= 0 {
			sqrt := math.Sqrt(disc)
			t0 := (-b - sqrt) / (2 * a)
			t1 := (-b + sqrt) / (2 * a)

			// Ok: the ray hits the infinite cone, now check if hits are within bounds
			y0 := ray.Origin.Y + t0*ray.Direction.Y
			if p.MinY < y0 && y0 < p.MaxY {
				xs = append(xs, t0)
			}

			y1 := ray.Origin.Y + t1*ray.Direction.Y
			if p.MinY < y1 && y1 < p.MaxY {
				xs = append(xs, t1)
			}
		}
	} else {
		// The ray is parallel to one of the cone halves, but may intersect the other
		if math.Abs(b) > Epsilon {
			t0 := -c / (2 * b)

			y0 := ray.Origin.Y + t0*ray.Direction.Y
			if p.MinY < y0 && y0 < p.MaxY {
				xs = append(xs, t0)
			}
		}
		// ...else the ray misses
	}

	// If cone is capped, we need to check for those intersections too
	checkCap := func(t, r float64) bool {
		x := ray.Origin.X + t*ray.Direction.X
		z := ray.Origin.Z + t*ray.Direction.Z

		return (x*x + z*z) <= r*r
	}

	if p.Capped {
		t0 := (p.MinY - ray.Origin.Y) / ray.Direction.Y
		if checkCap(t0, math.Abs(p.MinY)) {
			xs = append(xs, t0)
		}

		t1 := (p.MaxY - ray.Origin.Y) / ray.Direction.Y
		if checkCap(t1, math.Abs(p.MaxY)) {
			xs = append(xs, t1)
		}
	}

	return
}

func (p *Cone) LocalNormalAt(point Tuple) Tuple {
	if point.Y <= (p.MinY + Epsilon) {
		return Vector(0, -1, 0)
	}

	if point.Y >= (p.MaxY - Epsilon) {
		return Vector(0, +1, 0)
	}

	dist := point.X*point.X + point.Z*point.Z
	y := math.Sqrt(dist)
	if point.Y > 0 {
		y = -y
	}

	return Vector(point.X, y, point.Z)
}

func (p *Cone) NormalAtHit(point Tuple, ii *IntersectionInfo) Tuple {
	return p.LocalNormalAt(point)
}
