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
	s1.SetTransform(Translation(0, 1, 2))
	s1.SetMaterial(NewMaterial().SetDiffuseColor(CSS("orange")).SetDiffuse(1))

	s2 := NewSphere()
	s2.SetTransform(Translation(-2.5, 1, 6))
	s2.SetMaterial(NewMaterial().SetDiffuseColor(CSS("dodgerblue")).SetDiffuse(1))

	s3 := NewSphere()
	s3.SetTransform(Translation(1.5, 1, -1))
	s3.SetMaterial(NewMaterial().SetDiffuseColor(CSS("mediumspringgreen")).SetDiffuse(1))

	pole := NewCylinder(0, 2, true)
	pole.SetTransform(Scaling(0.1, 1, 0.1), Translation(-35, 0, 30))
	pole.SetMaterial(NewMaterial().SetDiffuseColor(CSS("crimson")).SetDiffuse(1))

	light := NewRectLight(RGB(1, 1, 1).Mul(0.9))
	light.SetSize(2, 2)
	light.SetDirection(Point(3, 5, -4), Point(0, 0, 0))

	world := NewWorld()
	world.SetAmbient(Gray(0.05))

	world.AddObjects(floor, s1, s2, s3, pole)

	world.AddLights(light)

	world.Options.AreaLightAdaptiveMinDepth = 6 // 5 is not enough to resolve the specular highlight on the green ball

	camera := NewCamera(800, 400, Pi/3)
	camera.SetTransform(EyeViewpoint(Point(0, 1.5, -5), Point(0, 1, 0), Vector(0, 1, 0)))
	// Good parameters for blur: LensRadius=0.2, FocalDistance=7, Supersampling=4

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

	light := NewPointLight(Point(0, 4, -4), RGB(1, 1, 1).Mul(0.9))

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

func applyTexture(s *Shape, filename string) *ImageTexture {
	txt := NewImageTexture()
	err := txt.LoadFromFile(filename)
	if err == nil {
		s.Material().SetPattern(txt).SetSpecular(0)
	} else {
		Debugln("Cannot load texture: ", filename)
	}

	return txt
}

func TestEarthFromSpace(t *testing.T) {
	t.SkipNow() // Test works but needs the texture files

	s0 := NewSphere()
	s0.SetTransform(Scaling(4.4), RotationY(-1.1), RotationZ(0.2), RotationX(0.4))
	s0.Material().SetAmbient(100).SetDiffuse(0)
	s0.SetShadow(false)
	applyTexture(s0, "../textures/2k_stars_milky_way.jpg")

	rot := Pi / 2.5

	s1 := NewSphere()
	s1.SetTransform(Scaling(1.20), RotationY(rot))
	s1.Material().SetSpecular(0).SetDiffuse(1).SetDiffuseColor(CSS("dodgerblue"))
	applyTexture(s1, "../textures/2k_earth_daymap.jpg")

	s2 := NewSphere()
	s2.SetTransform(Scaling(1.22), RotationY(rot))
	s2.Material().SetDiffuse(0).SetReflect(0, White).SetRefract(1, White).SetIor(1)
	s2.SetShadow(false)
	txt := applyTexture(s2, "../textures/2k_earth_clouds.jpg")

	// Add transparency to the texture, based on how bright it is
	for i := range txt.data {
		txt.data[i].a = txt.data[i].r*0.299 + txt.data[i].g*0.587 + txt.data[i].b*0.114
	}

	// Change the transparency of the sphere to match the texture at the hit
	txt.onApply = func(c ColorRGBA, ii *IntersectionInfo) {
		ii.Mat.DiffuseColor = c.RGB()
		ii.Mat.DiffuseLevel = 1.5 // Boost white a little
		ii.Mat.RefractLevel = float64(1 - c.a)
	}

	light := NewDirectionalLight(Vector(1, -1, 0.3), RGB(1, 1, 1).Mul(1))

	world := NewWorld()
	world.SetAmbient(Gray(0.01))

	world.AddObjects(s0, s1, s2)

	world.AddLights(light)
	world.Options.Supersampling = 4

	camera := NewCamera(800, 800, Pi/4)
	camera.SetTransform(EyeViewpoint(Point(0, 1.5, -4), Point(0, 0, 0), Vector(0, 1, 0)))

	world.RenderToPNG(camera, "test_earth_from_space.png")
}

func TestPlanets(t *testing.T) {
	t.SkipNow() // Test works but needs the texture files

	// These wonderful textures come from https://www.solarsystemscope.com/textures/

	// Universe
	s0 := NewSphere()
	s0.SetTransform(Scaling(6), RotationZ(Pi/3))

	applyTexture(s0, "../textures/2k_stars.jpg")

	// Mars
	s1 := NewSphere()
	s1.SetTransform(Translation(-1, -0.5, 0))
	s1.SetShadow(false)

	applyTexture(s1, "../textures/2k_mars.jpg")

	// Earth
	s2 := NewSphere()
	s2.SetTransform(Translation(0.2, -0.3, 1.5), RotationY(Pi/1.5))
	s2.SetShadow(false)

	applyTexture(s2, "../textures/2k_earth_daymap.jpg")

	// Jupiter
	s3 := NewSphere()
	s3.SetTransform(Translation(1.7, 0, 3), RotationY(Pi/6))
	s3.SetShadow(false)

	applyTexture(s3, "../textures/2k_jupiter.jpg")

	// Moon
	s4 := NewSphere()
	s4.SetTransform(Scaling(0.25), Translation(-5, 1.5, 12.5), RotationY(Pi-Pi/1.5), RotationX(-Pi/8))
	s4.SetShadow(false)

	applyTexture(s4, "../textures/2k_moon.jpg")

	light := NewPointLight(Point(0, 4, -4), RGB(1, 1, 1).Mul(0.9))

	world := NewWorld()
	world.SetAmbient(Gray(0.1))

	world.AddObjects(s0, s1, s2, s3, s4)

	world.AddLights(light)

	camera := NewCamera(300*4, 300*4, Pi/4)
	camera.SetTransform(EyeViewpoint(Point(0, 1.5, -4), Point(-0.15, 0, 0), Vector(0, 1, 0)))

	world.RenderToPNG(camera, "test_planets.png")
}

func TestNormalMaps(t *testing.T) {
	t.SkipNow()

	floor := NewPlane()
	floor.SetTransform(Scaling(4, -4, 4), RotationX(-Pi/2), Translation(0, -4, 0))

	txt := NewImageTexture()
	txt.linear = true
	// txt.LoadFromFile("head.jpg")
	// txt.LoadFromFile("wall.png")
	// txt.LoadFromFile("Well Preserved Chesterfield - (Normal Map_2).png")
	txt.LoadFromFile("marble_coloured_001_NRM.png")
	txt.wrap = TwPeriodic

	// txt.onMapUv = func(u, v float64, ii *IntersectionInfo) (float64, float64) {
	// 	return u / 3, v / 3
	// }

	txt3 := NewImageTexture()
	// txt.LoadFromFile("head.jpg")
	// txt.LoadFromFile("wall.png")
	// txt.LoadFromFile("Well Preserved Chesterfield - (Normal Map_2).png")
	txt3.LoadFromFile("marble_coloured_001_COLOR.png")
	txt3.wrap = TwPeriodic

	// txt3.onMapUv = func(u, v float64, ii *IntersectionInfo) (float64, float64) {
	// 	return u / 3, v / 3
	// }

	floor.Material().NormalMap = txt
	floor.Material().SetPattern(txt3)

	s1 := NewSphere()
	s1.SetTransform(Translation(0, 1, 0))
	s1.SetMaterial(NewMaterial().SetDiffuseColor(CSS("orange")).SetDiffuse(1))

	// txt2 := NewImageTexture()
	// txt2.linear = true
	// // txt2.LoadFromFile("NormalMap.jpg")
	// txt2.LoadFromFile("Well Preserved Chesterfield - (Normal Map_2).png")

	// txt2.onMapUv = func(u, v float64, ii *IntersectionInfo) (float64, float64) {
	// 	return u * 3, v * 3
	// }

	// txt2 := NewImageTexture()
	// txt2.linear = true
	// // txt2.LoadFromFile("NormalMap.jpg")
	// txt2.LoadFromFile("marble_coloured_001_NRM.png")

	// txt4 := NewImageTexture()
	// txt4.LoadFromFile("marble_coloured_001_COLOR.png")

	// txt2.onMapUv = func(u, v float64, ii *IntersectionInfo) (float64, float64) {
	// 	return u * 2, v * 2
	// }
	// txt4.onMapUv = func(u, v float64, ii *IntersectionInfo) (float64, float64) {
	// 	return u * 2, v * 2
	// }

	// s1.Material().NormalMap = txt2
	// s1.Material().SetPattern(txt2)

	// light := NewPointLight(Point(0,5,-5), RGB(1, 1, 1).Mul(0.9))
	// light := NewPointLight(Point(-5,5,-5), RGB(1, 1, 1).Mul(0.9))
	light := NewDirectionalLight(Vector(0, -2, 2), RGB(1, 1, 1).Mul(0.9))

	world := NewWorld()
	world.SetAmbient(Gray(0.1))

	world.AddObjects(floor /*, s1*/)

	world.AddLights(light)

	camera := NewCamera(800, 800, Pi/4)
	camera.SetTransform(EyeViewpoint(Point(0, 1, -5), Point(0, 1, 0), Vector(0, 1, 0)))

	world.RenderToPNG(camera, "test_normal_maps.png")
}
