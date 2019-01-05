// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
)

type Sphere struct {
}

// NewSphere returns a Shape based on a sphere with radius=1 and center (0,0,0)
func NewSphere() *Shape {
	return NewShape("sphere", &Sphere{})
}

func (s *Sphere) Bounds() Box {
	return Box{Point(-1, -1, -1), Point(+1, +1, +1)}
}

func (s *Sphere) LocalIntersect(ray Ray) []float64 {
	sphereToRay := ray.Origin.Sub(Point(0, 0, 0))
	a := ray.Direction.DotProduct(ray.Direction)
	b := 2 * ray.Direction.DotProduct(sphereToRay)
	c := sphereToRay.DotProduct(sphereToRay) - 1
	discriminant := b*b - 4*a*c

	if discriminant >= 0 {
		t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
		t2 := (-b + math.Sqrt(discriminant)) / (2 * a)

		return []float64{t1, t2}
	}

	return nil
}

func (s *Sphere) LocalNormalAt(point Tuple) Tuple {
	// Just convert the point into a vector, which is the normal in object coordinates
	point.W = 0

	return point
}
