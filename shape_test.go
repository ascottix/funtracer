// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
	"testing"
)

// Sphere
func TestSphereNormal(t *testing.T) {
	s := NewSphere()

	if !s.NormalAt(Point(1, 0, 0)).Equals(Vector(1, 0, 0)) {
		t.Errorf("normal at x failed")
	}

	if !s.NormalAt(Point(0, 1, 0)).Equals(Vector(0, 1, 0)) {
		t.Errorf("normal at y failed")
	}

	if !s.NormalAt(Point(0, 0, 1)).Equals(Vector(0, 0, 1)) {
		t.Errorf("normal at z failed")
	}

	d := math.Sqrt(3) / 3
	n := s.NormalAt(Point(d, d, d))
	if !n.Equals(Vector(d, d, d)) {
		t.Errorf("normal at non-axial point failed")
	}

	if !n.Equals(n.Normalize()) {
		t.Errorf("normal vector is not normalized")
	}

	s.SetTransform(Translation(0, 1, 0))
	if !s.NormalAt(Point(0, 1.70711, -0.70711)).Equals(Vector(0, 0.70711, -0.70711)) {
		t.Errorf("normal of translated sphere failed")
	}

	s.SetTransform(Scaling(1, 0.5, 1).Mul(RotationZ(Pi / 5)))
	if !s.NormalAt(Point(0, math.Sqrt(2)/2, -math.Sqrt(2)/2)).Equals(Vector(0, 0.97014, -0.24254)) {
		t.Errorf("normal of scaled and rotated sphere failed")
	}
}

// Plane
func TestPlaneNormalAt(t *testing.T) {
	p := NewPlane()

	n1 := p.NormalAt(Point(0, 0, 0))
	n2 := p.NormalAt(Point(10, 0, -10))
	n3 := p.NormalAt(Point(-5, 0, 150))

	expected := Vector(0, 1, 0)

	if !n1.Equals(expected) || !n1.Equals(n2) || !n1.Equals(n3) {
		t.Errorf("plane normal failed")
	}
}

func TestPlaneIntersect(t *testing.T) {
	p := NewPlane()

	r := NewRay(Point(0, 10, 0), Vector(0, 0, 1))
	if p.Intersect(r).Len() != 0 {
		t.Errorf("plane intersection (parallel ray) should be empty")
	}

	r = NewRay(Point(0, 10, 0), Vector(0, 0, 1))
	if p.Intersect(r).Len() != 0 {
		t.Errorf("plane intersection (coplanar ray) should be empty")
	}

	r = NewRay(Point(0, 1, 0), Vector(0, -1, 0))
	xs := p.Intersect(r)
	if xs == nil || xs.Len() != 1 {
		t.Errorf("plane intersection (ray from above) should hit once")
	}

	r = NewRay(Point(0, -1, 0), Vector(0, 1, 0))
	xs = p.Intersect(r)
	if xs == nil || xs.Len() != 1 {
		t.Errorf("plane intersection (ray from below) should hit once")
	}
}

// Cube
func TestCubeIntersect(t *testing.T) {
	c := NewCube()

	hit := func(px, py, pz, dx, dy, dz, t1, t2 float64) {
		r := NewRay(Point(px, py, pz), Vector(dx, dy, dz))
		xs := c.shapable.LocalIntersect(r)

		if len(xs) != 2 || !FloatEqual(xs[0], t1) || !FloatEqual(xs[1], t2) {
			t.Errorf("cube intersect failed: %+v = %v", r, xs)
		}
	}

	hit(5, 0.5, 0, -1, 0, 0, 4, 6)
	hit(-5, 0.5, 0, 1, 0, 0, 4, 6)
	hit(0.5, 5, 0, 0, -1, 0, 4, 6)
	hit(0.5, -5, 0, 0, 1, 0, 4, 6)
	hit(0.5, 0, 5, 0, 0, -1, 4, 6)
	hit(0.5, 0, -5, 0, 0, 1, 4, 6)
	hit(0, 0.5, 0, 0, 0, 1, -1, 1) // Inside

	miss := func(px, py, pz, dx, dy, dz float64) {
		r := NewRay(Point(px, py, pz), Vector(dx, dy, dz))
		xs := c.shapable.LocalIntersect(r)

		if len(xs) != 0 {
			t.Errorf("cube miss failed: %+v", r)
		}
	}

	miss(-2, 0, 0, 0.2673, 0.5345, 0.8018)
	miss(0, -2, 0, 0.8018, 0.2673, 0.5345)
	miss(0, 0, -2, 0.5345, 0.8018, 0.2673)
	miss(2, 0, -2, 0, 0, -1)
	miss(0, 2, 2, 0, -1, 0)
	miss(2, 2, 0, -1, 0, 0)
}

func TestCubeNormal(t *testing.T) {
	c := NewCube()

	test := func(px, py, pz, nx, ny, nz float64) {
		n := c.shapable.LocalNormalAt(Point(px, py, pz))

		if !n.Equals(Vector(nx, ny, nz)) {
			t.Errorf("cube normal failed: %+v should be %+v", n, Vector(nx, ny, nz))
		}
	}

	test(1, 0.5, -0.8, 1, 0, 0)
	test(-1, -0.2, 0.9, -1, 0, 0)
	test(-0.4, 1, -0.1, 0, 1, 0)
	test(0.3, -1, -0.7, 0, -1, 0)
	test(-0.6, 0.3, 1, 0, 0, 1)
	test(0.4, 0.4, -1, 0, 0, -1)
	test(1, 1, 1, 1, 0, 0)
	test(-1, -1, -1, -1, 0, 0)
}

// Triangle
func TestTriangleNew(t *testing.T) {
	p1, p2, p3 := Point(0, 1, 0), Point(-1, 0, 0), Point(1, 0, 0)

	s := NewTriangle(p1, p2, p3)

	if !s.P1.Equals(p1) || !s.P2.Equals(p2) || !s.P3.Equals(p3) {
		t.Errorf("bad triangle points")
	}

	if !s.E1.Equals(Vector(-1, -1, 0)) || !s.E2.Equals(Vector(1, -1, 0)) || !s.N.Equals(Vector(0, 0, -1)) {
		t.Errorf("bad triangle edges and normal")
	}
}

func TestTriangleIntersect(t *testing.T) {
	p1, p2, p3 := Point(0, 1, 0), Point(-1, 0, 0), Point(1, 0, 0)

	s := NewTriangle(p1, p2, p3)

	xs := s.LocalIntersect(NewRay(Point(0, -1, -2), Vector(0, 1, 0)))

	if len(xs) != 0 {
		t.Errorf("triangle intersections (parallel ray) should be empty")
	}

	xs = s.LocalIntersect(NewRay(Point(1, 1, -2), Vector(0, 0, 1)))

	if len(xs) != 0 {
		t.Errorf("triangle intersections (p1-p3 edge) should be empty")
	}

	xs = s.LocalIntersect(NewRay(Point(-1, 1, -2), Vector(0, 0, 1)))

	if len(xs) != 0 {
		t.Errorf("triangle intersections (p1-p2 edge) should be empty")
	}

	xs = s.LocalIntersect(NewRay(Point(0, -1, -2), Vector(0, 0, 1)))

	if len(xs) != 0 {
		t.Errorf("triangle intersections (p2-p3 edge) should be empty")
	}

	xs = s.LocalIntersect(NewRay(Point(0, 0.5, -2), Vector(0, 0, 1)))

	if len(xs) != 1 || !FloatEqual(xs[0], 2) {
		t.Errorf("triangle intersections should be a hit")
	}
}

// Cylinder
func TestCylinderIntersect(t *testing.T) {
	c := NewInfiniteCylinder()

	hit := func(px, py, pz, dx, dy, dz float64, ts ...float64) {
		r := NewRay(Point(px, py, pz), Vector(dx, dy, dz).Normalize())
		xs := c.shapable.LocalIntersect(r)

		if len(xs) != len(ts) {
			t.Errorf("cylinder intersect bad length: %+v = %d expected %d", r, len(xs), len(ts))
		} else {
			for i, v := range ts {
				if !FloatEqual(xs[i], v) {
					t.Errorf("cylinder intersect failed: %+v = %f, %f", r, xs[i], v)
				}
			}
		}
	}

	hit2 := func(px, py, pz, dx, dy, dz float64) {
		r := NewRay(Point(px, py, pz), Vector(dx, dy, dz).Normalize())
		xs := c.shapable.LocalIntersect(r)

		if len(xs) != 2 {
			t.Errorf("cylinder intersect bad length: %+v = %d should be 2", r, len(xs))
		}
	}

	hit(1, 0, 0, 0, 1, 0)
	hit(0, 0, 0, 0, 1, 0)
	hit(0, 0, -5, 1, 1, 1)

	hit(1, 0, -5, 0, 0, 1, 5, 5)
	hit(0, 0, -5, 0, 0, 1, 4, 6)
	hit(0.5, 0, -5, 0.1, 1, 1, 6.80798, 7.08872)

	c = NewCylinder(1, 2, false)

	hit(0, 1.5, 0, 0.1, 1, 0)
	hit(0, 3, -5, 0, 0, 1)
	hit(0, 0, -5, 0, 0, 1)
	hit(0, 2, -5, 0, 0, 1)
	hit(0, 1, -5, 0, 0, 1)
	hit(0, 1.5, -2, 0, 0, 1, 1, 3)

	c = NewCylinder(1, 2, true)

	hit2(0, 3, 0, 0, -1, 0)
	hit2(0, 3, -2, 0, -1, 2)
	hit2(0, 4, -2, 0, -1, 1)
	hit2(0, 0, -2, 0, 1, 2)
	hit2(0, -1, -2, 0, 1, 1)
}

func TestCylinderNormal(t *testing.T) {
	c := NewInfiniteCylinder()

	test := func(px, py, pz, nx, ny, nz float64) {
		n := c.shapable.LocalNormalAt(Point(px, py, pz))

		if !n.Equals(Vector(nx, ny, nz)) {
			t.Errorf("cylinder normal failed: %+v should be %+v", n, Vector(nx, ny, nz))
		}
	}

	test(1, 0, 0, 1, 0, 0)
	test(0, 5, -1, 0, 0, -1)
	test(0, -2, 1, 0, 0, 1)
	test(-1, 1, 0, -1, 0, 0)

	c = NewCylinder(1, 2, true)

	test(0, 1, 0, 0, -1, 0)
	test(0.5, 1, 0, 0, -1, 0)
	test(0, 1, 0.5, 0, -1, 0)
	test(0, 2, 0, 0, 1, 0)
	test(0.5, 2, 0, 0, 1, 0)
	test(0, 2, 0.5, 0, 1, 0)
}