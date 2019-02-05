// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func TestObjGibberish(t *testing.T) {
	data := `
There was a lady named Bright
who could walk faster than light.
She set out one day
in a relative way,
and came back the previous night.    
	`
	ParseWavefrontObjFromString(data)
}

func TestObjVertex(t *testing.T) {
	data := `
v -1 1 0
v -1.0000 0.5000 0.0000 
v 1 0 0
v 1 1 0
	`
	info := ParseWavefrontObjFromString(data)

	v := info.V

	if len(v) != 4 || !v[0].Equals(Point(-1, 1, 0)) || !v[1].Equals(Point(-1, 0.5, 0)) || !v[2].Equals(Point(1, 0, 0)) || !v[3].Equals(Point(1, 1, 0)) {
		t.Errorf("obj vertex failed")
	}
}

func TestObjTriangle(t *testing.T) {
	data := `
v -1 1 0 
v -1 0 0 
v 1 0 0 
v 1 1 0

f 1 2 3 
f 1 3 4
	`
	info := ParseWavefrontObjFromString(data)

	if len(info.F) != 2 {
		t.Errorf("obj triangle failed")
		t.SkipNow()
	}

	t1 := info.F[0]
	t2 := info.F[1]

	if t1.V[0] != 0 || t1.V[1] != 1 || t1.V[2] != 2 || t2.V[0] != 0 || t2.V[1] != 2 || t2.V[2] != 3 {
		t.Errorf("obj triangle vertices mismatch")
	}
}

func TestObjPolygon(t *testing.T) {
	data := `
v -1 1 0 
v -1 0 0 
v 1 0 0 
v 1 1 0 
v 0 2 0
f 1 2 3 4 5
	`
	info := ParseWavefrontObjFromString(data)

	if len(info.F) != 3 {
		t.Errorf("obj polygon failed")
		t.SkipNow()
	}

	t1 := info.F[0]
	t2 := info.F[1]
	t3 := info.F[2]

	if t1.V[0] != 0 || t1.V[1] != 1 || t1.V[2] != 2 || t2.V[0] != 0 || t2.V[1] != 2 || t2.V[2] != 3 || t3.V[0] != 0 || t3.V[1] != 3 || t3.V[2] != 4 {
		t.Errorf("obj polygon vertices mismatch")
	}
}

func TestObjGroup(t *testing.T) {
	data := `
v -1 1 0 
v -1 0 0 
v 1 0 0 
v 1 1 0 
v 0 2 0
g FirstGroup
f 1 2 3 
g SecondGroup
f 1 4 5
	`
	info := ParseWavefrontObjFromString(data)

	if len(info.Groups) != 3 {
		t.Errorf("obj group count failed")
	}
}

func TestObjDodecahedron(t *testing.T) {
	TestWithImage(t)

	data := `
v  -0.57735  -0.57735  0.57735
v  0.934172  0.356822  0
v  0.934172  -0.356822  0
v  -0.934172  0.356822  0
v  -0.934172  -0.356822  0
v  0  0.934172  0.356822
v  0  0.934172  -0.356822
v  0.356822  0  -0.934172
v  -0.356822  0  -0.934172
v  0  -0.934172  -0.356822
v  0  -0.934172  0.356822
v  0.356822  0  0.934172
v  -0.356822  0  0.934172
v  0.57735  0.57735  -0.57735
v  0.57735  0.57735  0.57735
v  -0.57735  0.57735  -0.57735
v  -0.57735  0.57735  0.57735
v  0.57735  -0.57735  -0.57735
v  0.57735  -0.57735  0.57735
v  -0.57735  -0.57735  -0.57735

f  19  3  2
f  12  19  2
f  15  12  2
f  8  14  2
f  18  8  2
f  3  18  2
f  20  5  4
f  9  20  4
f  16  9  4
f  13  17  4
f  1  13  4
f  5  1  4
f  7  16  4
f  6  7  4
f  17  6  4
f  6  15  2
f  7  6  2
f  14  7  2
f  10  18  3
f  11  10  3
f  19  11  3
f  11  1  5
f  10  11  5
f  20  10  5
f  20  9  8
f  10  20  8
f  18  10  8
f  9  16  7
f  8  9  7
f  14  8  7
f  12  15  6
f  13  12  6
f  17  13  6
f  13  1  11
f  12  13  11
f  19  12  11
	`

	info := ParseWavefrontObjFromString(data)
	info.Normalize()

	world := NewWorld()

	world.AddLights(NewPointLight(Point(-2, 8, -9), White))

	mesh := NewTrimesh(info, -1)

	group := NewGroup()

	mesh.AddToGroup(group)

	group.BuildBVH()

	world.AddObjects(group)

	camera := NewCamera(400*2, 400*2, Pi/5)
	camera.SetTransform(EyeViewpoint(Point(0, 0, -5), Point(0, 0, 0), Vector(0, 1, 0)))

	world.RenderToPNG(camera, "test_obj_dodecahedron.png")
}

func TestObjTeapot(t *testing.T) {
	TestWithImage(t)

	info := ParseWavefrontObjFromFile("scenes/teapot.obj")

	world := NewWorld()

	world.Ambient = Gray(0.2)

	if false {
		// Use a spot-light
		world.AddLights(NewSpotLight(Point(-10, 20, -10), Point(0, 0, 0), 0.1, 0.2, Gray(0.7)))
	} else {
		// Use an area light
		light := NewRectLight(RGB(1, 1, 1).Mul(0.8))
		light.SetSize(3, 3)
		light.SetDirection(Point(-10, 20, -10), Point(0, 0, 0))
		world.AddLights(light)
	}

	info.Normalize()

	mesh := NewTrimesh(info, -1)

	group := NewGroup()
	group.SetTransform(RotationX(-Pi/2), Scaling(3, 3, 4), Translation(0, 0, 0.12))
	m := NewMaterial()
	m.SetDiffuseColor(RGB(0.87, 0.87, 0.9))
	m.SetReflect(0.3, White)
	m.SetDiffuse(1)
	m.SetShininess(256)
	mesh.SetMaterial(m)

	mesh.AddToGroup(group)

	group.BuildBVH()

	world.AddObjects(group)

	p := NewCheckerPattern(Gray(0.5), Gray(0.7))
	p.SetTransform(Translation(0, 0.1, 0), Scaling(0.7))
	x := NewMaterial()
	x.SetPattern(p)
	x.SetSpecular(0)

	wall := NewPlane()
	wall.SetTransform(RotationX(Pi/2), Translation(0, 5, 0))
	wall.SetMaterial(x)

	floor := NewPlane()
	floor.SetTransform(Translation(0, -1, 0))
	floor.SetMaterial(x)

	world.AddObjects(floor, wall)

	camera := NewCamera(400*2, 200*2, 1)
	camera.SetTransform(EyeViewpoint(Point(0, 2, -9), Point(0, 0.5, 0), Vector(0, 1, 0)))

	// Zoom on shadow
	if false {
		camera.SetFieldOfView(0.5)
		camera.SetTransform(EyeViewpoint(Point(0, 2, -9), Point(1.5, 0, 0), Vector(0, 1, 0)))
	}

	world.RenderToPNG(camera, "test_obj_teapot.png")
}

func TestObjDragon(t *testing.T) {
	TestWithImage(t)

	info := ParseWavefrontObjFromFile("../obj/dragon.obj")

	world := NewWorld()

	world.Ambient = Gray(0.1)

	world.AddLights(NewSpotLight(Point(-10, 20, -10), Point(0, 0, 0), 0.2, 0.25, Gray(0.9)))

	info.Normalize()

	mesh := NewTrimesh(info, -1)

	group := NewGroup()

	/* Interesting parameters for the Stanford Dragon */
	// m := NewMaterial()
	// m.SetPattern(JadePattern())
	// m.SetAmbient(0.2)
	// m.Pattern.SetTransform(Scaling(0.3, 0.3, 0.3))
	// m.SetReflect(0.1, White)
	// m.SetDiffuse(1)
	// m.SetShininess(10)

	m := NewMaterial()
	m.SetAmbient(0)
	m.SetDiffuse(0)
	m.SetReflect(0.05, White)
	m.SetRefract(0.95, CSS("CadetBlue"))
	m.SetIor(1.1)

	group.SetTransform(Scaling(2), Translation(0, 0.2, 0))

	camera := NewCamera(1920/4, 1080/4, 0.7)
	camera.SetTransform(EyeViewpoint(Point(0, 2, -9), Point(0, 0.3, 0), Vector(0, 1, 0)))

	mesh.SetMaterial(m)

	mesh.AddToGroup(group)

	group.BuildBVH()

	world.AddObjects(group)

	p := NewCheckerPattern(Gray(0.5), Gray(0.7))
	p.SetTransform(Translation(0, 0.1, 0), Scaling(0.7))
	x := NewMaterial()
	x.SetPattern(p)
	x.SetSpecular(0)

	wall := NewPlane()
	wall.SetTransform(RotationX(Pi/2), Translation(0, 5, 0))
	wall.SetMaterial(x)

	floor := NewPlane()
	floor.SetTransform(Translation(0, -1, 0))
	floor.SetMaterial(x)

	world.AddObjects(floor, wall)

	world.RenderToPNG(camera, "test_obj_dragon.png")
}
