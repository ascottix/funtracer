// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
	"testing"
)

func TestGroupIntersect(t *testing.T) {
	g := NewGroup()

	r := NewRay(Point(0, 0, 0), Vector(0, 0, 1))

	xs := g.LocalIntersect(r)

	if xs.Len() != 0 {
		t.Errorf("group intersect should be empty")
	}

	s1 := NewSphere()
	s2 := NewSphere()
	s2.SetTransform(Translation(0, 0, -3))
	s3 := NewSphere()
	s3.SetTransform(Translation(5, 0, 0))

	g.Add(s1, s2, s3)

	r = NewRay(Point(0, 0, -5), Vector(0, 0, 1))

	xs = g.LocalIntersect(r)

	if xs.Len() != 4 || xs.At(0).O != s2 || xs.At(1).O != s2 || xs.At(2).O != s1 || xs.At(3).O != s1 {
		t.Errorf("group intersect failed: %+v", *xs)
	}

	g = NewGroup()
	g.SetTransform(Scaling(2, 2, 2))

	s := NewSphere()
	s.SetTransform(Translation(5, 0, 0))

	g.Add(s)

	r = NewRay(Point(10, 0, -10), Vector(0, 0, 1))

	xs = NewIntersections()
	g.AddIntersections(r, xs)

	if xs.Len() != 2 {
		t.Errorf("group intersect (transformed) failed")
	}
}

func TestGroupWorldToObject(t *testing.T) {
	g1 := NewGroup()
	g1.SetTransform(RotationY(Pi / 2))

	g2 := NewGroup()
	g2.SetTransform(Scaling(2, 2, 2))

	g1.Add(g2)

	s := NewSphere()
	s.SetTransform(Translation(5, 0, 0))

	g2.Add(s)

	p := s.WorldToObject(Point(-2, 0, -10))

	if !p.Equals(Point(0, 0, -1)) {
		t.Errorf("world to object failed: %+v", p)
	}
}

func TestGroupNormalToWorld(t *testing.T) {
	g1 := NewGroup()
	g1.SetTransform(RotationY(Pi / 2))

	g2 := NewGroup()
	g2.SetTransform(Scaling(1, 2, 3))

	g1.Add(g2)

	s := NewSphere()
	s.SetTransform(Translation(5, 0, 0))

	g2.Add(s)

	n := s.NormalToWorld(Vector(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3))

	if !n.Equals(Vector(0.285714, 0.428571, -0.857142)) {
		t.Errorf("normal to world failed: %+v", n)
	}
}

func createHexagonScene() *Scene {
	scene := NewScene()

	world := scene.World

	material := func() *Material {
		m := NewMaterial()
		p := WhiteLinesPattern()
		p.SetTransform(Scaling(0.6))
		m.SetPattern(p)
		return m
	}

	hexagon_corner := func() *Shape {
		corner := NewSphere()
		corner.SetTransform(Translation(0, 0, -1), Scaling(0.25, 0.25, 0.25))
		return corner
	}

	hexagon_edge := func() *Shape {
		edge := NewCylinder(0, 1, false)
		edge.SetTransform(Translation(0, 0, -1), RotationY(-Pi/6), RotationZ(-Pi/2), Scaling(0.25, 1, 0.25))
		return edge
	}

	hexagon_side := func() *Group {
		side := NewGroup()
		side.Add(hexagon_corner(), hexagon_edge())
		return side
	}

	hexagon := func() *Group {
		hex := NewGroup()
		for n := 0; n < 6; n++ {
			side := hexagon_side()
			side.SetTransform(RotationY(float64(n) * Pi / 3))
			hex.Add(side)

			if n == 0 {
				hex.Add(side.BoundingBox())
			}
		}

		hex.SetMaterial(material())

		return hex
	}

	world.AddLights(NewPointLight(Point(2, 5, -5), Gray(0.8)), NewDirectionalLight(Vector(0, -1, 0), Gray(0.2)))

	floor := NewPlane()
	floor.SetTransform(Translation(0, -10, 0))

	world.AddObjects(hexagon(), floor)

	camera := NewCamera(400, 400, 0.8)
	camera.SetTransform(EyeViewpoint(Point(-1, 3, -4), Point(0, 0, 0), Vector(0, 1, 0)))

	scene.Camera = camera

	return scene
}

func TestGroupHexagon(t *testing.T) {
	TestWithImage(t)

	scene := createHexagonScene()

	scene.World.ErpCanvasToImage = ErpLinear
	scene.World.RenderToPNG(scene.Camera, "test_hexagon.png")
}
