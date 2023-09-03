// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package shapes

import (
	"math"

	. "ascottix/funtracer/maths"
	. "ascottix/funtracer/textures"
)

type Cylinder struct {
	MinY   float64
	MaxY   float64
	Capped bool
}

func NewCylinder(miny, maxy float64, capped bool) *Shape {
	return NewShape("cylinder", &Cylinder{miny, maxy, capped})
}

func NewInfiniteCylinder() *Shape {
	return NewCylinder(math.Inf(-1), math.Inf(+1), false)
}

func (p *Cylinder) Bounds() Box {
	return Box{Point(-1, p.MinY, -1), Point(+1, p.MaxY, +1)}
}

func (p *Cylinder) LocalIntersect(ray Ray) (xs []float64) {
	a := Square(ray.Direction.X) + Square(ray.Direction.Z)

	// Check intersection with cylinder surface
	if math.Abs(a) > Epsilon {
		b := 2 * (ray.Origin.X*ray.Direction.X + ray.Origin.Z*ray.Direction.Z)

		c := Square(ray.Origin.X) + Square(ray.Origin.Z) - 1

		disc := b*b - 4*a*c

		if disc >= 0 {
			sqrt := math.Sqrt(disc)
			t0 := (-b - sqrt) / (2 * a)
			t1 := (-b + sqrt) / (2 * a)

			// Ok: the ray hits the infinite cylinder, now check if hits are within bounds
			y0 := ray.Origin.Y + t0*ray.Direction.Y
			if p.MinY < y0 && y0 < p.MaxY {
				xs = append(xs, t0)
			}

			y1 := ray.Origin.Y + t1*ray.Direction.Y
			if p.MinY < y1 && y1 < p.MaxY {
				xs = append(xs, t1)
			}
		}
	}
	// ...else the ray is parallel to the y axis

	// If cylinder is capped, we need to check for those intersections too
	checkCap := func(t float64) bool {
		x := ray.Origin.X + t*ray.Direction.X
		z := ray.Origin.Z + t*ray.Direction.Z

		return (x*x + z*z) <= 1
	}

	if p.Capped {
		t0 := (p.MinY - ray.Origin.Y) / ray.Direction.Y
		if checkCap(t0) {
			xs = append(xs, t0)
		}

		t1 := (p.MaxY - ray.Origin.Y) / ray.Direction.Y
		if checkCap(t1) {
			xs = append(xs, t1)
		}
	}

	return
}

func (p *Cylinder) LocalNormalAt(point Tuple) Tuple {
	dist := point.X*point.X + point.Z*point.Z

	if dist < 1 {
		if point.Y <= (p.MinY + Epsilon) {
			return Vector(0, -1, 0)
		}

		if point.Y >= (p.MaxY - Epsilon) {
			return Vector(0, +1, 0)
		}
	}

	return Vector(point.X, 0, point.Z)
}

func (p *Cylinder) NormalAtHit(point Tuple, ii *IntersectionInfo) Tuple {
	// See http://cse.csusb.edu/tongyu/courses/cs520/notes/texture.php for (u,v)

	return p.LocalNormalAt(point)
}
