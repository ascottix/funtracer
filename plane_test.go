// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func createPlaneScene() (*World, *Camera) {
	newMaterial := func(c Color) *Material {
		return NewMaterial().SetDiffuseColor(c)
	}

	floor := NewPlane()
	floor.SetMaterial(newMaterial(RGB(1, 0.9, 0.9)).SetSpecular(0))

	middle := NewSphere()
	middle.SetTransform(Translation(-0.5, 1, 0.5))
	middle.SetMaterial(newMaterial(RGB(0.1, 1, 0.5)).SetDiffuse(0.7).SetSpecular(0.3))

	right := NewSphere()
	right.SetTransform(Translation(1.5, 0.5, -0.5), Scaling(0.5, 0.5, 0.5))
	right.SetMaterial(newMaterial(RGB(0.5, 1, 0.1)).SetDiffuse(0.7).SetSpecular(0.3))

	left := NewSphere()
	left.SetTransform(Translation(-1.5, 0.33, -0.75), Scaling(0.33, 0.33, 0.33))
	left.SetMaterial(newMaterial(RGB(1, 0.8, 0.1)).SetDiffuse(0.7).SetSpecular(0.3))

	light := NewPointLight(Point(-10, 10, -10), RGB(0.8, 0.8, 0.8))

	light2 := NewDirectionalLight(Vector(0, -1, 0), RGB(0.2, 0.2, 0.2))

	world := NewWorld()

	world.AddObjects(floor, middle, left, right)

	world.AddLights(light, light2)

	world.ErpCanvasToImage = ErpLinear // Test scenes from book are not gamma corrected

	camera := NewCamera(640, 320, Pi/3)
	camera.SetTransform(EyeViewpoint(Point(0, 1.5, -5), Point(0, 1, 0), Vector(0, 1, 0)))

	return world, camera
}

func TestPlaneScene(t *testing.T) {
	TestWithImage(t)

	world, camera := createPlaneScene()

	world.RenderToPNG(camera, "test_simple_scene_w_plane.png")
}

func createPlaneSceneWithPatterns() (*World, *Camera) {
	world, camera := createPlaneScene()

	floor := getHittableAt(world, 0)
	cp := NewCheckerPattern(White, RGB(1, 0.3, 0.3))
	cp.SetTransform(Translation(0, 0.1, 0)) // Bump pattern a little bit to avoid "pattern acne"
	floor.Material().SetPattern(cp)

	middle := getHittableAt(world, 1)
	sp := NewStripePattern(RGB(0, 0.6, 0.2), RGB(0, 0.8, 0.4))
	sp.SetTransform(Scaling(0.15, 0.15, 0.15).Mul(RotationZ(Pi / 2)).Mul(RotationY(+Pi / 4)))
	middle.Material().SetPattern(sp)

	right := getHittableAt(world, 3)
	gp := NewGradientPattern(RGB(1, 0.4, 0.0), RGB(1, 1, 0.6))
	gp.SetTransform(Translation(-1.5, 0, 0), Scaling(2, 1, 1), RotationY(-Pi/8))
	right.Material().SetPattern(gp)

	return world, camera
}

func TestPlaneSceneWithPatterns(t *testing.T) {
	TestWithImage(t)

	world, camera := createPlaneSceneWithPatterns()

	world.RenderToPNG(camera, "test_simple_scene_w_patterns.png")
}

func TestPlaneSceneWithReflection(t *testing.T) {
	TestWithImage(t)

	world, camera := createPlaneSceneWithPatterns()

	floor := getHittableAt(world, 0)
	floor.Material().SetDiffuse(1)
	floor.Material().SetReflective(0.4)

	left := getHittableAt(world, 2)
	left.Material().SetPattern(NewSolidColorPattern(White))
	left.Material().SetDiffuse(0.7)
	left.Material().SetReflective(1)

	world.ErpCanvasToImage = ErpLinear
	world.RenderToPNG(camera, "test_simple_scene_w_reflection.png")
}

func TestPlaneSceneWithRefraction(t *testing.T) {
	TestWithImage(t)

	world, camera := createPlaneSceneWithPatterns()

	floor := getHittableAt(world, 0)
	floor.Material().SetDiffuse(1)
	floor.Material().SetReflective(0.4)

	left := getHittableAt(world, 2)
	left.Material().SetPattern(NewSolidColorPattern(White))
	left.Material().SetDiffuse(0.7)
	left.Material().SetReflective(1)

	right := getHittableAt(world, 3)
	right.Material().SetPattern(NewSolidColorPattern(White))
	right.Material().SetDiffuse(0.1)
	right.Material().SetShininess(300)
	right.Material().SetReflective(1)
	right.Material().SetRefractive(1, 1.52)

	wall := NewPlane()
	wall.SetTransform(RotationX(Pi/2), Translation(0, 5, -2), RotationZ(Pi/4))
	wall.Material().SetReflective(0.5)
	world.AddObjects(wall)

	camera.SetFieldOfView(Pi / 2)

	world.ErpCanvasToImage = ErpLinear
	world.RenderToPNG(camera, "test_simple_scene_w_refraction.png")
}
