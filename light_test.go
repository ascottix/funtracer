// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
	"testing"
)

func TestLightning(t *testing.T) {
	s := NewSphere() // Object is not important here, only the material
	position := Point(0, 0, 0)

	test := func(eyev, normalv Tuple, light *PointLight, expected Color, scenario string) {
		c, _ := Lighten(light, s, position, eyev, normalv, false)
		c = c.Add(Gray(0.1)) // Ambient
		if !c.Equals(expected) {
			t.Errorf("%s failed: %+v should be %+v", scenario, c, expected)
		}
	}

	test(Vector(0, 0, -1), Vector(0, 0, -1), NewPointLight(Point(0, 0, -10), White), Gray(1.9), "eye between light and surface")

	test(Vector(0, math.Sqrt(2)/2, -math.Sqrt(2)/2), Vector(0, 0, -1), NewPointLight(Point(0, 0, -10), White), White, "eye between light and surface, eye offset 45°")

	test(Vector(0, 0, -1), Vector(0, 0, -1), NewPointLight(Point(0, 10, -10), White), Gray(0.7364), "eye opposite surface, light offset 45°")

	test(Vector(0, -math.Sqrt(2)/2, -math.Sqrt(2)/2), Vector(0, 0, -1), NewPointLight(Point(0, 10, -10), White), Gray(1.6364), "eye in the path of the reflection vector")

	test(Vector(0, 0, -1), Vector(0, 0, -1), NewPointLight(Point(0, 0, 10), White), Gray(0.1), "light behind the surface")
}

func TestShadow(t *testing.T) {
	s := NewSphere() // Object is not important here, only the material
	position := Point(0, 0, 0)
	eyev := Vector(0, 0, -1)
	normalv := Vector(0, 0, -1)
	light := NewPointLight(Point(0, 0, -10), RGB(1, 1, 1))
	shadowed := true

	c, _ := Lighten(light, s, position, eyev, normalv, shadowed)
	if !c.Equals(Black) {
		t.Errorf("shadow failed")
	}
}

func TestPattern(t *testing.T) {
	s := NewSphere()
	s.Material().SetPattern(NewStripePattern(White, Black))
	s.Material().SetAmbient(1).SetDiffuse(0).SetSpecular(0)
	eyev := Vector(0, 0, -1)
	normalv := Vector(0, 0, -1)
	light := NewPointLight(Point(0, 0, -10), White)

	c1, o1 := Lighten(light, s, Point(0.9, 0, 0), eyev, normalv, false)
	c2, o2 := Lighten(light, s, Point(1.1, 0, 0), eyev, normalv, false)

	// We need to add the colors to the light to simulate ambient
	if !c1.Add(o1).Equals(White) || !c2.Add(o2).Equals(Black) {
		t.Errorf("light on strip pattern failed: %+v, %+v", c1, c2)
	}
}

func TestRectLight(t *testing.T) {
	TestWithImage(t)

	floor := NewPlane()

	s1 := NewSphere()
	s1.SetTransform(Translation(0, 1, 0))
	s1.SetMaterial(NewMaterial().SetDiffuseColor(CSS("orange")).SetDiffuse(1))

	s2 := NewSphere()
	s2.SetTransform(Translation(-2.5, 1, 0))
	s2.SetMaterial(NewMaterial().SetDiffuseColor(CSS("dodgerblue")).SetDiffuse(1))

	light := NewRectLight(RGB(1, 1, 1).Mul(0.9))
	light.SetSize(2, 2)
	light.SetDirection(Point(3, 5, -4), Point(0, 0, 0))

	world := NewWorld()
	world.SetAmbient(Gray(0.05))

	world.AddObjects(floor, s1, s2)

	world.AddLights(light)

	camera := NewCamera(640, 320, Pi/3)
	camera.SetTransform(EyeViewpoint(Point(0, 1.5, -5), Point(0, 1, 0), Vector(0, 1, 0)))

	world.RenderToPNG(camera, "test_rect_light.png")
}

func TestDepthOfField(t *testing.T) {
	TestWithImage(t)

	s1 := NewSphere()
	s1.SetTransform(Translation(0, 1, 0))
	s1.SetMaterial(NewMaterial().SetDiffuseColor(CSS("orange")).SetDiffuse(1))

	s2 := NewSphere()
	s2.SetTransform(Translation(-2, 1, +4))
	s2.SetMaterial(NewMaterial().SetDiffuseColor(CSS("dodgerblue")).SetDiffuse(1))

	light := NewPointLight(Point(0,4,-4), RGB(1, 1, 1).Mul(0.9))

	world := NewWorld()
	world.SetAmbient(Gray(0.05))

	world.AddObjects(s1, s2)

	world.AddLights(light)

	world.Options.Supersampling = 16
	world.Options.LensRadius = 0.2
	world.Options.FocalDistance = 5 // 5 sets focus on orange ball, 9 on blue ball

	camera := NewCamera(300, 300, Pi/4)
	camera.SetTransform(EyeViewpoint(Point(0, 1.5, -5), Point(0, 1, 0), Vector(0, 1, 0)))

	world.RenderToPNG(camera, "test_depth_of_field.png")
}
