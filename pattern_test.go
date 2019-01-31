// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func TestPatternStripe(t *testing.T) {
	p := NewStripePattern(White, Black)

	if !p.PatternAt(Point(0, 0, 0)).Equals(White) || !p.PatternAt(Point(0, 1, 0)).Equals(White) || !p.PatternAt(Point(0, 2, 0)).Equals(White) {
		t.Errorf("stripe pattern should be constant in y")
	}

	if !p.PatternAt(Point(0, 0, 0)).Equals(White) || !p.PatternAt(Point(0, 0, 1)).Equals(White) || !p.PatternAt(Point(0, 0, 2)).Equals(White) {
		t.Errorf("stripe pattern should be constant in z")
	}

	if !p.PatternAt(Point(0, 0, 0)).Equals(White) || !p.PatternAt(Point(0.9, 0, 0)).Equals(White) || !p.PatternAt(Point(-1.1, 0, 0)).Equals(White) {
		t.Errorf("stripe pattern in x should be white")
	}

	if !p.PatternAt(Point(1, 0, 0)).Equals(Black) || !p.PatternAt(Point(-0.1, 0, 0)).Equals(Black) || !p.PatternAt(Point(-1, 0, 0)).Equals(Black) {
		t.Errorf("stripe pattern in x should be black")
	}
}

func TestPatternStripeAtObject(t *testing.T) {
	s := NewSphere()
	s.SetTransform(Scaling(2, 2, 2))

	p := NewStripePattern(White, Black)
	c := p.ColorAt(s, Point(1.5, 0, 0))

	if !c.Equals(White) {
		t.Errorf("bad stripe with object transform")
	}

	s = NewSphere()
	p = NewStripePattern(White, Black)
	p.SetTransform(Scaling(2, 2, 2))

	if !c.Equals(White) {
		t.Errorf("bad stripe with pattern transform")
	}

	s = NewSphere()
	s.SetTransform(Scaling(2, 2, 2))
	p = NewStripePattern(White, Black)
	p.SetTransform(Translation(0.5, 0, 0))

	if !c.Equals(White) {
		t.Errorf("bad stripe with both pattern and object transform")
	}
}

func TestPatternGradient(t *testing.T) {
	p := NewGradientPattern(White, Black)

	if !p.PatternAt(Point(0, 0, 0)).Equals(White) || !p.PatternAt(Point(0.25, 0, 0)).Equals(RGB(0.75, 0.75, 0.75)) || !p.PatternAt(Point(0.5, 0, 0)).Equals(RGB(0.5, 0.5, 0.5)) || !p.PatternAt(Point(0.75, 0, 0)).Equals(RGB(0.25, 0.25, 0.25)) {
		t.Errorf("gradient pattern failed")
	}
}

func TestPatternScene(t *testing.T) {
	TestWithImage(t)

	world := NewWorld()
	world.ErpCanvasToImage = ErpLinear

	newSphere := func(pattern Pattern, tx, tz float64) *Material {
		s := NewSphere()
		s.SetTransform(Translation(tx, +1, tz))

		if pattern != nil {
			s.Material().SetPattern(pattern)
		}

		world.AddObjects(s)

		return s.Material()
	}

	floor := NewPlane()
	floor.Material().SetPattern(NewCheckerPattern(White, RGB(1, 0.3, 0.3)))
	floor.Material().SetReflective(0.2)
	floor.Material().Pattern.SetTransform(Translation(0, 0.1, 0)) // Bump pattern a little bit to avoid "pattern acne"

	wall := NewPlane()
	wall.Material().SetPattern(NewSolidColorPattern(RGB(0.4, 0.5, 0.5)))
	wall.Material().SetReflective(0.2)
	wall.SetTransform(RotationX(Pi/2), Translation(0, 15, 0))

	light1 := NewPointLight(Point(-10, 10, -10), RGB(0.8, 0.8, 0.8))

	light2 := NewDirectionalLight(Vector(0, -1, 0), RGB(0.2, 0.2, 0.2))

	world.AddObjects(floor, wall)

	world.AddLights(light1, light2)

	newSphere(JadePattern(), 0, 1).SetRefractive(0.3, 1.6).SetReflective(0.2)

	newSphere(WhiteLinesPattern(), -2, -1).SetReflective(0.05).SetShininess(300)

	newSphere(NewSolidColorPattern(RGB(0.3, 0.3, 0.3)), +2, -2).SetReflective(1).SetRefractive(1, 1.52).SetDiffuse(0.1).SetShininess(300)

	newSphere(NewSolidColorPattern(RGB(0.9, 0.9, 1)), -5, +1).SetSpecular(0.3).SetReflective(0.1).SetShininess(5).SetDiffuse(0.6)

	newSphere(AmberPattern(), 5, 1).SetReflective(0.05).SetShininess(300)

	camera := NewCamera(800, 400, Pi/3)
	camera.SetTransform(EyeViewpoint(Point(0, 0+2.5, -10), Point(0, 0, 0), Vector(0, 1, 0)))

	world.RenderToPNG(camera, "test_scene_with_patterns.png")
}

func createReflectionsAndRefractionsScene() *Scene {
	world := NewWorld()

	light := NewPointLight(Point(-4.9, 4.9, -1), Gray(1))

	world.AddLights(light)

	wallMaterial := NewMaterial()
	wallMaterial.SetPattern(NewStripePattern(Gray(0.45), Gray(0.55)))
	wallMaterial.Pattern.SetTransform(Scaling(0.25, 0.25, 0.25), RotationY(1.5708))
	wallMaterial.SetAmbient(0).SetDiffuse(0.4).SetSpecular(0).SetReflective(0.3)

	floor := NewPlane()
	floor.SetTransform(RotationY(0.31415))
	floor.Material().SetPattern(NewCheckerPattern(Gray(0.65), Gray(0.35)))
	floor.Material().Pattern.SetTransform(Translation(0, 0.1, 0))
	floor.Material().SetSpecular(0).SetReflective(0.4)

	ceiling := NewPlane()
	ceiling.SetTransform(Translation(0, 5, 0))
	ceiling.Material().SetPattern(NewSolidColorPattern(RGB(0.8, 0.8, 0.8)))
	ceiling.Material().SetSpecular(0).SetAmbient(0.3)

	westWall := NewPlane()
	westWall.SetTransform(Translation(-5, 0, 0), RotationZ(1.5708), RotationY(1.5708))
	westWall.SetMaterial(wallMaterial)

	eastWall := NewPlane()
	eastWall.SetTransform(Translation(5, 0, 0), RotationZ(1.5708), RotationY(1.5708))
	eastWall.SetMaterial(wallMaterial)

	northWall := NewPlane()
	northWall.SetTransform(Translation(0, 0, 5), RotationX(1.5708))
	northWall.SetMaterial(wallMaterial)

	southWall := NewPlane()
	southWall.SetTransform(Translation(0, 0, -5), RotationX(1.5708))
	southWall.SetMaterial(wallMaterial)

	group := NewGroup()

	addBackSphere := func(scale, tx, ty, tz, r, g, b float64) {
		s := NewSphere()
		s.Material().SetPattern(NewSolidColorPattern(RGB(r, g, b)))
		s.Material().SetShininess(50)
		s.SetTransform(Translation(tx, ty, tz), Scaling(scale, scale, scale))

		group.Add(s)
	}

	addBackSphere(0.4, 4.6, 0.4, 1, 0.8, 0.5, 0.3)
	addBackSphere(0.3, 4.7, 0.3, 0.4, 0.9, 0.4, 0.5)
	addBackSphere(0.5, -1, 0.5, 4.5, 0.4, 0.9, 0.6)
	addBackSphere(0.3, -1.7, 0.3, 4.7, 0.4, 0.6, 0.9)

	addForeSphere := func(scale, tx, ty, tz, r, g, b float64, glass bool) {
		s := NewSphere()
		s.Material().SetPattern(NewSolidColorPattern(RGB(r, g, b)))
		s.SetTransform(Translation(tx, ty, tz), Scaling(scale, scale, scale))

		if glass {
			s.Material().SetAmbient(0).SetDiffuse(0.4).SetSpecular(0.9).SetShininess(300).SetReflective(0.9).SetRefractive(0.9, 1.5)
		} else {
			s.Material().SetSpecular(0.4).SetShininess(5)
		}

		group.Add(s)
	}

	addForeSphere(1, -0.6, 1, 0.6, 1, 0.3, 0.2, false)
	addForeSphere(0.7, 0.6, 0.7, -0.6, 0, 0, 0.2, true)
	addForeSphere(0.5, -0.7, 0.5, -0.8, 0, 0.2, 0, true)

	// Note: don't add infinite objects to a group if going to use a BVH
	world.AddObjects(floor, ceiling, westWall, eastWall, northWall, southWall)

	group.BuildBVH()

	world.AddObjects(group)

	camera := NewCamera(800, 400, 1.152)
	camera.SetTransform(EyeViewpoint(Point(-2.6, 1.5, -3.9), Point(-0.6, 1, -0.8), Vector(0, 1, 0)))

	scene := NewScene()
	scene.World = world
	scene.Camera = camera

	return scene
}

func TestReflectionsAndRefractions(t *testing.T) {
	TestWithImage(t)

	scene := createReflectionsAndRefractionsScene()

	scene.World.ErpCanvasToImage = ErpLinear
	scene.World.RenderToPNG(scene.Camera, "test_reflections_and_refractions.png")
}

func TestSceneWithBoundingBoxes(t *testing.T) {
	TestWithImage(t)

	world := NewWorld()

	light := NewPointLight(Point(-4.9, 4.9, -1), RGB(1, 1, 1))

	world.AddLights(light)

	wallMaterial := NewMaterial()
	wallMaterial.SetPattern(NewStripePattern(Gray(0.45), Gray(0.55)))
	wallMaterial.Pattern.SetTransform(Scaling(0.25, 0.25, 0.25), RotationY(1.5708))
	wallMaterial.SetAmbient(0).SetDiffuse(0.4).SetSpecular(0).SetReflective(0.3)

	floor := NewPlane()
	floor.SetTransform(RotationY(0.31415))
	floor.Material().SetPattern(NewCheckerPattern(Gray(0.65), Gray(0.35)))
	floor.Material().Pattern.SetTransform(Translation(0, 0.1, 0))
	floor.Material().SetSpecular(0).SetReflective(0.4)

	ceiling := NewPlane()
	ceiling.SetTransform(Translation(0, 5, 0))
	ceiling.Material().SetPattern(NewSolidColorPattern(Gray(0.8)))
	ceiling.Material().SetSpecular(0).SetAmbient(0.3)

	westWall := NewPlane()
	westWall.SetTransform(Translation(-5, 0, 0), RotationZ(1.5708), RotationY(1.5708))
	westWall.SetMaterial(wallMaterial)

	eastWall := NewPlane()
	eastWall.SetTransform(Translation(5, 0, 0), RotationZ(1.5708), RotationY(1.5708))
	eastWall.SetMaterial(wallMaterial)

	northWall := NewPlane()
	northWall.SetTransform(Translation(0, 0, 5), RotationX(1.5708))
	northWall.SetMaterial(wallMaterial)

	southWall := NewPlane()
	southWall.SetTransform(Translation(0, 0, -5), RotationX(1.5708))
	southWall.SetMaterial(wallMaterial)

	addBackSphere := func(group *Group, scale, tx, ty, tz, r, g, b float64) {
		s := NewSphere()
		s.Material().SetPattern(NewSolidColorPattern(RGB(r, g, b)))
		s.Material().SetShininess(50)
		s.SetTransform(Translation(tx, ty, tz), Scaling(scale, scale, scale))

		group.Add(s)
	}

	g1 := NewGroup()
	g2 := NewGroup()

	world.AddObjects(g1, g2)

	addBackSphere(g1, 0.4, 4.6, 0.4, 1, 0.8, 0.5, 0.3)
	addBackSphere(g1, 0.3, 4.7, 0.3, 0.4, 0.9, 0.4, 0.5)
	addBackSphere(g2, 0.5, -1, 0.5, 4.5, 0.4, 0.9, 0.6)
	addBackSphere(g2, 0.3, -1.7, 0.3, 4.7, 0.4, 0.6, 0.9)

	groupFore := NewGroup()

	world.AddObjects(groupFore)

	addForeSphere := func(scale, tx, ty, tz, r, g, b float64, glass bool) {
		s := NewSphere()
		s.Material().SetPattern(NewSolidColorPattern(RGB(r, g, b)))
		s.SetTransform(Translation(tx, ty, tz), Scaling(scale, scale, scale))

		if glass {
			s.Material().SetAmbient(0).SetDiffuse(0.4).SetSpecular(0.9).SetShininess(300).SetReflective(0.9).SetRefractive(0.9, 1.5)
		} else {
			s.Material().SetSpecular(0.4).SetShininess(5)
		}

		groupFore.Add(s)
	}

	addForeSphere(1, -0.6, 1, 0.6, 1, 0.3, 0.2, false)
	addForeSphere(0.7, 0.6, 0.7, -0.6, 0, 0, 0.2, true)
	addForeSphere(0.5, -0.7, 0.5, -0.8, 0, 0.2, 0, true)

	world.AddObjects(g1.BoundingBox(), g2.BoundingBox())
	world.AddObjects(groupFore.BoundingBox())

	world.AddObjects(floor, ceiling, westWall, eastWall, northWall, southWall)

	camera := NewCamera(800, 400, 1.152)
	camera.SetTransform(EyeViewpoint(Point(-2.6, 1.5, -3.9), Point(-0.6, 1, -0.8), Vector(0, 1, 0)))

	world.ErpCanvasToImage = ErpLinear
	world.RenderToPNG(camera, "test_reflections_w_bounding_boxes.png")
}
