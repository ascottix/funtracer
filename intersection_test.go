// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func NewGlassSphere() *Shape {
	s := NewSphere()

	s.Material().SetRefractive(1, 1.5)

	return s
}

func TestIntersectionRefraction(t *testing.T) {
	a := NewGlassSphere()
	b := NewGlassSphere()
	c := NewGlassSphere()

	a.SetTransform(Scaling(2, 2, 2))
	a.Material().SetRefractive(1, 1.5)

	b.SetTransform(Translation(0, 0, -0.25))
	b.Material().SetRefractive(1, 2.0)

	c.SetTransform(Translation(0, 0, +0.25))
	c.Material().SetRefractive(1, 2.5)

	ray := NewRay(Point(0, 0, -4), Vector(0, 0, 1))

	xs := NewIntersections()

	xs.Add(a, 2, 6)
	xs.Add(b, 2.75, 4.75)
	xs.Add(c, 3.25, 5.25)
	xs.Sort()

	n1x := []float64{1.0, 1.5, 2.0, 2.5, 2.5, 1.5}
	n2x := []float64{1.5, 2.0, 2.5, 2.5, 1.5, 1.0}

	for i, n1 := range n1x {
		n2 := n2x[i]

		ii := NewIntersectionInfo(xs.At(i), ray, xs)

		if !FloatEqual(ii.N1, n1) {
			t.Errorf("bad n1(%d): %f should be %f", i, ii.N1, n1)
		}

		if !FloatEqual(ii.N2, n2) {
			t.Errorf("bad n2(%d): %f should be %f", i, ii.N2, n2)
		}
	}
}

func TestIntersectionFooBar(t *testing.T) {
	t.SkipNow()

	world := NewWorld()

	light1 := NewPointLight(Point(-5, 5, -3), Gray(0.5))
	light2 := NewPointLight(Point(0, 4, -8), Gray(0.4))
	world.AddLights(light1, light2)

	back := NewPlane()
	back.SetTransform(RotationX(Pi/2), Translation(0, 5, 0))
	back.Material().SetPattern(NewCheckerPattern(Gray(0.7), Gray(0.9))).SetSpecular(0)
	back.Material().Pattern.SetTransform(Translation(0, 0.1, 0), Scaling(0.7))

	floor := NewPlane()
	floor.SetTransform(Translation(0, -1, 0))
	floor.Material().SetPattern(back.Material().Pattern)

	world.AddObjects(back, floor)

	cylinder := NewCylinder(1, 2, true) // (-1, +1, false)
	cylinder.SetTransform(Translation(0, -0.9, 0), Scaling(0.05, 1, 0.05))
	cylinder.SetMaterial(MatMatte(CSS("forestgreen")))

	cylinder2 := NewCylinder(1, 2, true) // (-1, +1, false)
	cylinder2.SetTransform(Translation(0.165, 0.75, 0), RotationZ(0.5), Scaling(0.049, 0.3, 0.049))
	cylinder2.SetMaterial(MatMatte(CSS("forestgreen")))

	// TODO!
	// To render the image sphere's LocalNormalAt method must be modified as follows:
	// func (s *Sphere) LocalNormalAt(point Tuple) Tuple {
	// 	// Just convert the point into a vector, which is the normal in object coordinates
	// 	point.W = 0

	// 	R := float64()
	// 	N := Perlin(point.Mul(10))
	// 	point.X += N*R
	// 	point.Y += N*R
	// 	point.Z += N*R

	// 	return point.Normalize()

	// 	return point
	// }

	world.AddObjects(cylinder, cylinder2)

	sphere := NewSphere()
	sphere.SetMaterial(MatMatte(CSS("orangered")))
	sphere.SetTransform(Scaling(1, 0.88, 1))
	world.AddObjects(sphere)

	camera := NewCamera(800, 400, 1.5)
	camera.SetTransform(EyeViewpoint(Point(0, 2, -7), Point(0, 1, 0), Vector(0, 1, 0)))

	world.RenderToPNG(camera, "test_foobar.png")
}
