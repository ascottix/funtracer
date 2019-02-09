// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"image/png"
	"math"
	"os"
	"testing"
)

func TestNewRay(t *testing.T) {
	r := NewRay(Point(1, 2, 3), Vector(4, 5, 6))

	if !r.Origin.Equals(Point(1, 2, 3)) || !r.Direction.Equals(Vector(4, 5, 6)) {
		t.Errorf("new ray failed")
	}
}

func TestRayPosition(t *testing.T) {
	r := NewRay(Point(2, 3, 4), Vector(1, 0, 0))

	if !r.Position(0).Equals(Point(2, 3, 4)) || !r.Position(1).Equals(Point(3, 3, 4)) || !r.Position(-1).Equals(Point(1, 3, 4)) || !r.Position(2.5).Equals(Point(4.5, 3, 4)) {
		t.Errorf("ray position failed")
	}
}

func TestRayIntersectSphere(t *testing.T) {
	s1 := NewSphere()
	xs := s1.Intersect(NewRay(Point(0, 0, -5), Vector(0, 0, 1)))

	if xs == nil || xs.Len() != 2 || xs.At(0).T != 4.0 || xs.At(1).T != 6.0 {
		t.Errorf("intersect front fail: %+v", xs)
	}

	s2 := NewSphere()
	xs = s2.Intersect(NewRay(Point(0, 1, -5), Vector(0, 0, 1)))

	if xs == nil || xs.Len() != 2 || xs.At(0).T != 5.0 || xs.At(1).T != 5.0 {
		t.Errorf("intersect tangent fail")
	}

	xs = NewSphere().Intersect(NewRay(Point(0, 2, -5), Vector(0, 0, 1)))

	if xs.Len() != 0 {
		t.Errorf("no intersect fail")
	}

	xs = NewSphere().Intersect(NewRay(Point(0, 0, 0), Vector(0, 0, 1))) // Origin inside sphere

	if xs == nil || xs.Len() != 2 || xs.At(0).T != -1.0 || xs.At(1).T != 1.0 {
		t.Errorf("intersect inside fail")
	}

	xs = NewSphere().Intersect(NewRay(Point(0, 0, 5), Vector(0, 0, 1))) // Sphere behind ray

	if xs == nil || xs.Len() != 2 || xs.At(0).T != -6.0 || xs.At(1).T != -4.0 {
		t.Errorf("intersect behind fail")
	}

	// Test transformations
	s1 = NewSphere()
	s1.SetTransform(Scaling(2, 2, 2))
	xs = s1.Intersect(NewRay(Point(0, 0, -5), Vector(0, 0, 1)))

	if xs == nil || xs.Len() != 2 || xs.At(0).T != 3.0 || xs.At(1).T != 7.0 {
		t.Errorf("intersect on scaling fail")
	}

	s2 = NewSphere()
	s2.SetTransform(Translation(5, 0, 0))
	xs = s2.Intersect(NewRay(Point(0, 0, -5), Vector(0, 0, 1)))

	if xs.Len() != 0 {
		t.Errorf("no intersect on translation fail")
	}
}

func TestRayHit(t *testing.T) {
	i1 := NewIntersections().Add(NewSphere(), 1, 2)

	if i1.Hit().T != 1 || i1.Hit().O == nil {
		t.Errorf("hit 1 failed")
	}

	i2 := NewIntersections().Add(NewSphere(), 1, -1)

	if i2.Hit().T != 1 || !i2.Hit().Valid() {
		t.Errorf("hit 2 failed")
	}

	i3 := NewIntersections().Add(NewSphere(), -1, -2)

	if i3.Hit().Valid() {
		t.Errorf("hit 3 failed %+v", i3.Hit())
	}

	i4 := NewIntersections().Add(NewSphere(), 5, 7, -3, 2)

	if i4.Hit().T != 2 {
		t.Errorf("hit 4 failed %+v", i4.Hit())
	}
}

func TestRayTransform(t *testing.T) {
	r1 := NewRay(Point(1, 2, 3), Vector(0, 1, 0)).Transform(Translation(3, 4, 5))

	if !r1.Origin.Equals(Point(4, 6, 8)) || !r1.Direction.Equals(Vector(0, 1, 0)) {
		t.Errorf("ray transform 1 failed")
	}

	r2 := NewRay(Point(1, 2, 3), Vector(0, 1, 0)).Transform(Scaling(2, 3, 4))

	if !r2.Origin.Equals(Point(2, 6, 12)) || !r2.Direction.Equals(Vector(0, 3, 0)) {
		t.Errorf("ray transform 2 failed %+v", r2)
	}
}

func TestRayReflect(t *testing.T) {
	s := NewPlane()
	r := NewRay(Point(0, 1, -1), Vector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	i := NewIntersection(math.Sqrt(2), s)
	ii := NewIntersectionInfo(i, r, nil)

	if !ii.Reflectv.Equals(Vector(0, math.Sqrt(2)/2, math.Sqrt(2)/2)) {
		t.Errorf("basic reflection test failed")
	}
}

func TestFlatSphereImage(t *testing.T) {
	t.SkipNow()

	p1 := Point(-5, -5, +5) // Wall top left
	p2 := Point(+5, +5, +5) // Wall bottom right
	pe := Point(0, 0, -3)   // Eye
	w := 240
	h := 240
	dx := (p2.X - p1.X) / float64(w)
	dy := (p2.Y - p1.Y) / float64(h)
	m := RGB(1, 0.5, 0.5) // Material

	c := NewCanvas(w, h)

	// Render!
	s := NewSphere()
	s.SetTransform(Shearing(1, 0, 0, 0, 0, 0).Mul(Scaling(0.5, 1, 1)))
	p := p1
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dir := p.Sub(pe) // Vector from eye to point on wall
			ray := NewRay(pe, dir)
			xs := s.Intersect(ray)
			if xs.Hit().Valid() {
				c.FastSetPixelAt(x, h-1-y, m)
			}
			p.X += dx
		}
		p.X = p1.X
		p.Y += dy
	}

	c.WriteAsPPM(os.Stdout)
}

func TestLightSphereImage(t *testing.T) {
	TestWithImage(t)

	p1 := Point(-5, -5, +5) // Wall top left
	p2 := Point(+5, +5, +5) // Wall bottom right
	pe := Point(0, 0, -2.5) // Eye
	w := 480
	h := 480
	dx := (p2.X - p1.X) / float64(w)
	dy := (p2.Y - p1.Y) / float64(h)

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Prepare scene
	s := NewSphere()
	s.Material().SetPattern(NewSolidColorPattern(RGB(1, 0.2, 1)))
	s.Material().Ambient = 0.1
	light := NewPointLight(Point(-10, 10, -10), RGB(1, 1, 1))

	// Render!
	p := p1
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dir := p.Sub(pe).Normalize() // Normalized vector from eye to point on wall
			ray := NewRay(pe, dir)       // Cast ray from eye to target point
			xs := s.Intersect(ray)       // Get intersections
			hit := xs.Hit()              // Get hit if any
			if hit.Valid() {
				point := ray.Position(hit.T) // Point hit by the ray
				ii := IntersectionInfo{Intersection: hit, Point: point}
				normal := hit.O.NormalAtHit(&ii, xs) // Normal to the surface in the hit point
				eye := dir.Neg()
				color, _ := Lighten(light, hit.O, point, eye, normal, false)
				img.Set(x, h-1-y, color)
			}
			p.X += dx
		}
		p.X = p1.X
		p.Y += dy
	}

	f, _ := os.Create("test_light_sphere.png")
	defer f.Close()
	png.Encode(f, img)
}
