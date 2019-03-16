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
	obj := ParseWavefrontObjFromString(data)

	v := obj.V

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
	obj := ParseWavefrontObjFromString(data)

	if len(obj.F) != 2 {
		t.Errorf("obj triangle failed")
		t.SkipNow()
	}

	t1 := obj.F[0]
	t2 := obj.F[1]

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
	obj := ParseWavefrontObjFromString(data)

	if len(obj.F) != 3 {
		t.Errorf("obj polygon failed")
		t.SkipNow()
	}

	t1 := obj.F[0]
	t2 := obj.F[1]
	t3 := obj.F[2]

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
	obj := ParseWavefrontObjFromString(data)

	if len(obj.Groups) != 3 {
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

	obj := ParseWavefrontObjFromString(data)
	obj.Normalize()

	world := NewWorld()

	world.AddLights(NewPointLight(Point(-2, 8, -9), White))

	mesh := NewTrimesh(obj, -1)

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

	obj := ParseWavefrontObjFromFile("scenes/teapot.obj")

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

	obj.Normalize()

	mesh := NewTrimesh(obj, -1)

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

	camera := NewCamera(800, 400, 1)
	camera.SetTransform(EyeViewpoint(Point(0, 2, -9), Point(0, 0.5, 0), Vector(0, 1, 0)))

	// Zoom on shadow (for testing the adaptive area light heuristics)
	if false {
		// At fov=0.3 there's a lot of shadow in the scene, even a 16x16 sampler cannot render
		// it correctly... use a minDepth of 7 for the adaptive sampler
		camera.SetFieldOfView(0.3)
		camera.SetTransform(EyeViewpoint(Point(0, 2, -9), Point(2.5, -0.75, 0), Vector(0, 1, 0)))
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

func TestObjSign(t *testing.T) {
	t.SkipNow()

	info := ParseWavefrontObjFromFile("sign/35 mph speed limit sign final.obj")

	info.Normalize()

	world := NewWorld()

	world.AddLights(NewPointLight(Point(-2, 8, 9), White))

	mesh := NewTrimesh(info, -1)

	txt := NewImageTexture()
	txt.LoadFromFile("sign/35 mph speed limit sign unwrap 4888.png")

	txt.onMapUv = func(u, v float64, ii *IntersectionInfo) (float64, float64) {
		return u, 1 - v
	}

	m := NewMaterial()
	m.SetPattern(txt)

	mesh.SetMaterial(m)

	group := NewGroup()
	group.SetTransform(Scaling(-1, 1, 1)) // Hack... texture is mirrored horizontally, so we flip the object to compensate

	mesh.AddToGroup(group)

	group.BuildBVH()

	world.AddObjects(group)

	camera := NewCamera(400*2, 400*2, Pi/5)
	camera.SetTransform(EyeViewpoint(Point(0, 1, 5), Point(0, 0, 0), Vector(0, 1, 0)))

	world.RenderToPNG(camera, "test_sign_35mph.png")
}

func TestObjCapsule(t *testing.T) {
	t.SkipNow()

	// info := ParseWavefrontObjFromFile("capsule/capsule.obj")
	info := ParseWavefrontObjFromFile("Castle/Castle OBJ.obj")

	info.Normalize()

	world := NewWorld()

	world.AddLights(NewPointLight(Point(-2, 8, 9), White))

	mesh := NewTrimesh(info, -1)

	group := NewGroup()
	// group.SetTransform(Scaling(1,1,1), RotationY(Pi/3), RotationX(-Pi/2))
	// group.SetTransform(RotationY(Pi/3), RotationX(+Pi/2))
	group.SetTransform(RotationY(Pi / 6))

	mesh.AddToGroup(group)

	group.BuildBVH()

	world.AddObjects(group)

	camera := NewCamera(400*2, 400*2, Pi/10)
	camera.SetTransform(EyeViewpoint(Point(0, 3, 5), Point(0, 0, 0), Vector(0, 1, 0)))

	world.RenderToPNG(camera, "test_capsule.png")
}

func TestObjMeshNormalMap(t *testing.T) {
	t.SkipNow()

	data := `
v 0 1 0
v 1 1 0
v 0 0 0
v 1 0 0

vt 0 1
vt 1 1
vt 0 0
vt 1 0

f 1/1 2/2 3/3 
f 2/2 4/4 3/3
	`
	info := ParseWavefrontObjFromString(data)

	info.Normalize()
	info.Autosmooth()

	world := NewWorld()

	world.AddLights(NewPointLight(Point(-2, 18, 9), White))

	mesh := NewTrimesh(info, -1)

	txt := NewImageTexture()
	txt.LoadFromFile("Wall_Stone_003_COLOR.jpg")

	// txt.onMapUv = func(u, v float64, ii *IntersectionInfo) (float64, float64) {
	// 	return u, 1 - v
	// }

	txt2 := NewImageTexture()
	txt2.linear = true
	txt2.LoadFromFile("Wall_Stone_003_NRM.jpg")

	m := NewMaterial()
	m.SetPattern(txt)
	m.NormalMap = txt2

	mesh.SetMaterial(m)

	group := NewGroup()
	// group.SetTransform(Scaling(-1, 1, 1)) // Hack... texture is mirrored horizontally, so we flip the object to compensate

	mesh.AddToGroup(group)

	group.BuildBVH()

	world.AddObjects(group)

	camera := NewCamera(400*2, 400*2, Pi/5)
	camera.SetTransform(EyeViewpoint(Point(0, 1, 5), Point(0, 0, 0), Vector(0, 1, 0)))

	world.RenderToPNG(camera, "test_mesh_normalmap.png")
}

// To use this scene, follows the [#] steps in order!
//
// [1] create a scenes/pan folder
// [2] visit the https://3dtextures.me/ site and download the Azulejos 003 texture from the Tiles section
// [3] unzip Azulejos_003_COLOR.jpg and Azulejos_003_NORM.jpg in the scenes/pan folder
// [4] download the model casserole_obj.zip from the original scene at http://www.oyonale.com/modeles.php?lang=en&page=41
// [5] unzip pan_obj.obj in the scenes/pan folder
// [6] comment line [6a] when ready
// [7] run go test -timeout 2h to render!
//
// Note: rendering takes about 20 minutes on my PC but it may be more or less on yours,
// disable supersampling or reduce the image size to go faster.
func TestObjPan(t *testing.T) {
	t.SkipNow() // [6a] comment this line to enable the test

	world := NewWorld()

	world.Ambient = Gray(0.2)

	// Note: a point light is much faster in case we need to just play with the geometry
	light := NewRectLight(RGB(1, 1, 1).Mul(0.9))
	light.SetSize(3, 1)
	light.SetDirection(Point(2, 0.75, -1), Point(0, 0, 0))
	world.AddLights(light)

	// Wall
	wall := NewPlane()
	wall.SetTransform(RotationX(Pi/2), Translation(0, 2, 0))

	// Texture for wall and floor
	txt_wall := NewImageTexture()
	txt_wall.LoadFromFile("scenes/pan/Azulejos_003_COLOR.jpg")
	txt_wall.onMapUv = func(u, v float64, ii *IntersectionInfo) (float64, float64) {
		return u, v + 0.26 // To align the tiles border with the planes crossing line
	}
	wall.Material().SetPattern(txt_wall)
	wall.Material().SetDiffuse(0.5)

	// Normal map gives just a slight depth to the tiles separation in this scene, but still worth it
	nrm_wall := NewImageTexture()
	nrm_wall.linear = true
	nrm_wall.LoadFromFile("scenes/pan/Azulejos_003_NORM.jpg")
	nrm_wall.onMapUv = txt_wall.onMapUv
	wall.Material().NormalMap = nrm_wall

	// Floor will use the same texture as the wall
	floor := NewPlane()
	floor.SetTransform(Translation(0, -0.5, 0))

	floor.Material().SetPattern(txt_wall)
	floor.Material().SetDiffuse(1.1)
	floor.Material().NormalMap = nrm_wall
	floor.Material().SetReflect(0.2, White)

	world.AddObjects(floor, wall)

	// Pan
	obj := ParseWavefrontObjFromFile("scenes/pan/pan_obj.obj")

	obj.Normalize()
	obj.Autosmooth()

	pan_body := NewTrimesh(obj, 1)
	pan_handle := NewTrimesh(obj, 2)
	pan_joint := NewTrimesh(obj, 3)

	// Copper material
	copper_color := CSS("#b85d33")
	copper := NewMaterial()
	copper.SetDiffuse(0.6)
	copper.SetDiffuseColor(copper_color)
	copper.SetReflect(0.4, copper_color.Mul(1.1))
	copper.Roughness = 0.1
	copper.SetSpecular(3)
	copper.SetShininess(10)

	pan_body.SetMaterial(copper)
	pan_joint.SetMaterial(copper)

	// Black plastic material
	black_plastic := NewMaterial()
	black_plastic.SetDiffuseColor(Gray(0.015))
	copper.Roughness = 0.5
	black_plastic.SetShininess(10)

	pan_handle.SetMaterial(black_plastic)

	// Create a group for the meshes so we can use BVH and optimize, or rendering would take forever
	group := NewGroup()
	group.SetTransform(Translation(0, -0.15, 0), RotationY(-Pi/2))

	pan_body.AddToGroup(group)
	pan_handle.AddToGroup(group)
	pan_joint.AddToGroup(group)

	group.BuildBVH()

	world.AddObjects(group)

	// Camera
	camera := NewCamera(800, 800, 0.15) // Note: 400x400 is usually good enough to test changes
	camera.SetTransform(EyeViewpoint(Point(-6, 4, -9), Point(0, 0, 0), Vector(0, 1, 0)))

	// Render!
	world.Options.Supersampling = 4 // Note: comment out to just "play" with the scene

	world.RenderToPNG(camera, "test_obj_pan.png")
}

func TestObjIss(t *testing.T) {
	t.SkipNow()

	obj := ParseWavefrontObjFromFile("scenes/ISS/iss.obj")

	obj.Normalize()
	obj.Autosmooth()

	world := NewWorld()

	world.AddLights(NewPointLight(Point(-3, 10, 15), Gray(0.8)))
	world.AddLights(NewDirectionalLight(Point(0, -3, 1), Gray(0.3)))

	mesh := NewTrimesh(obj, -1)

	txt_diff := NewImageTexture()
	txt_diff.LoadFromFile("scenes/ISS/textures/iss_diffuse.png")
	txt_diff.onMapUv = TextureMirrorV
	mesh.Material().SetPattern(txt_diff)

	txt_norm := NewImageTexture()
	txt_norm.linear = true
	txt_norm.LoadFromFile("scenes/ISS/textures/iss_normal.png")
	txt_norm.onMapUv = TextureMirrorV
	mesh.Material().NormalMap = txt_norm

	s := 10.

	iss := NewGroup()
	iss.SetTransform(RotationX(-Pi / 2))
	iss.SetTransform(Translation(0, s*0.927, 5.4), RotationX(Pi/2+Pi/10), Scaling(0.9))

	mesh.AddToGroup(iss)

	iss.BuildBVH()

	earth := NewSphere()
	earth.SetTransform(Scaling(s), RotationX(2.5))
	earth.Material().SetDiffuse(1.2)
	applyTexture(earth, "../textures/2k_earth_daymap.jpg")

	atmo := NewSphere()
	atmo.SetTransform(Scaling(s+0.04), RotationX(2.5))
	atmo.Material().SetDiffuse(0.1).SetDiffuseColor(CSS("dodgerblue")).SetReflect(0.01, White).SetRefract(0.9, CSS("white").Mul(1.2)).SetIor(1)
	atmo.SetShadow(false)

	clouds := NewSphere()
	clouds.SetTransform(Scaling(s+0.01), RotationX(2.5))
	clouds.Material().SetDiffuse(0).SetReflect(0, White).SetRefract(1, CSS("white")).SetIor(1)
	clouds.SetShadow(false)
	txt := applyTexture(clouds, "../textures/8k_earth_clouds.jpg") // Use 2k for test and 8k for final rendering

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

	world.AddObjects(iss, earth, atmo, clouds)

	camera := NewCamera(400*3, 300*3, Pi/7.2)
	camera.SetTransform(EyeViewpoint(Point(0, s-0.5, s), Point(0, s*0.9, 0), Vector(0, 1, 0)))

	world.Options.Supersampling = 4 // Note: comment out to just "play" with the scene

	world.RenderToPNG(camera, "test_iss.png")
}
