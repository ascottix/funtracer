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

	test := func(eyev, normalv Tuple, light Light, expected Color, scenario string) {
		c := light.Lighten(s.Material().Pattern.ColorAt(s, position), s, position, eyev, normalv, false)
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

	oc := s.Material().Pattern.ColorAt(s, position)
	if !light.Lighten(oc, s, position, eyev, normalv, shadowed).Equals(Black) {
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

	o1 := s.Material().Pattern.ColorAt(s, Point(0.9, 0, 0))
	o2 := s.Material().Pattern.ColorAt(s, Point(1.1, 0, 0))
	c1 := light.Lighten(o1, s, Point(0.9, 0, 0), eyev, normalv, false)
	c2 := light.Lighten(o2, s, Point(1.1, 0, 0), eyev, normalv, false)

	// We need to add the colors to the light to simulate ambient
	if !c1.Add(o1).Equals(White) || !c2.Add(o2).Equals(Black) {
		t.Errorf("light on strip pattern failed: %+v, %+v", c1, c2)
	}
}
