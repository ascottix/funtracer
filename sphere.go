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

func (s *Sphere) NormalAtHit(point Tuple, ii *IntersectionInfo) Tuple {
	// Compute (u,v) coordinates of point, starting from the consideration that a point (x, y, z)
	// on the unit sphere can be expressed as (sin θ cos φ, cos θ, sin θ sin φ)
	phi := math.Atan2(point.Z, point.X) + Pi
	theta := math.Acos(point.Y)

	ii.U = phi / (2 * Pi)
	ii.V = theta / Pi

	normal := s.LocalNormalAt(point)

	if nmap := ii.GetNormalMap(); nmap != nil {
		n := nmap.NormalAtHit(ii)

		A := Vector(0, 1, 0)
		T := normal.CrossProduct(A).Normalize() // Tangent vector
		B := T.CrossProduct(normal)             // Bitangent vector

		// TODO: we're multiplying by -n.Y here but it may be texture-dependent!
		// Yes... it was texture dependent... we need to find a way to handle that
		ii.SurfNormalv = (T.Mul(n.X).Add(B.Mul(n.Y)).Add(normal.Mul(n.Z))).Normalize()
		ii.HasSurfNormalv = true
	}

	return normal
}
