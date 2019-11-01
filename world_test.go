// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
	"testing"
)

const (
	SkipAllTestsWithImages = false
)

func TestWithImage(t *testing.T) {
	if SkipAllTestsWithImages {
		t.SkipNow()
	}
}

func TestWorldNew(t *testing.T) {
	w := NewWorld()

	if len(w.Objects) != 0 || len(w.Lights) != 0 {
		t.Errorf("world is not empty!")
	}
}

func createDefaultWorld() *World {
	w := NewWorld()

	s1 := NewSphere()
	s1.Material().SetDiffuse(0.7).SetSpecular(0.2).SetDiffuseColor(RGB(0.8, 1.0, 0.6))
	s2 := NewSphere()
	s2.SetTransform(Scaling(0.5, 0.5, 0.5))

	w.AddObjects(s1, s2)

	l1 := NewPointLight(Point(-10, 10, -10), RGB(1, 1, 1))

	w.AddLights(l1)

	return w
}

func TestWorldDefault(t *testing.T) {
	w := createDefaultWorld()

	if len(w.Lights) != 1 {
		t.Errorf("bad default world light")
	}

	if len(w.Objects) != 2 {
		t.Errorf("bad default world objects")
	}
}

func TestWorldIntersect(t *testing.T) {
	r := NewRay(Point(0, 0, -5), Vector(0, 0, 1))

	w := createDefaultWorld()
	xs := w.Intersect(r)

	if xs != nil {
		xs.Sort() // Sort may be skipped for performance
	}

	if xs == nil || xs.Len() != 4 || xs.At(0).T != 4 || xs.At(1).T != 4.5 || xs.At(2).T != 5.5 || xs.At(3).T != 6 {
		t.Errorf("bad world intersects")
	}

	i1 := NewIntersection(4, NewSphere())
	ii := NewIntersectionInfo(i1, r, nil)

	if ii.T != i1.T || ii.O != i1.O || !ii.Point.Equals(Point(0, 0, -1)) || !ii.Eyev.Equals(Vector(0, 0, -1)) || !ii.Normalv.Equals(Vector(0, 0, -1)) {
		t.Errorf("bad intersection info 1: %+v", ii)
	}

	if ii.Inside {
		t.Errorf("bad outside intersection info")
	}

	ii = NewIntersectionInfo(NewIntersection(1, NewSphere()), NewRay(Point(0, 0, 0), Vector(0, 0, 1)), nil)

	if !ii.Point.Equals(Point(0, 0, 1)) || !ii.Eyev.Equals(Vector(0, 0, -1)) || !ii.Inside || !ii.Normalv.Equals(Vector(0, 0, -1)) {
		t.Errorf("bad inside intersection info: %+v", ii)
	}

	s1 := NewSphere()
	s1.SetTransform(Translation(0, 0, 1))
	ii = NewIntersectionInfo(NewIntersection(5, s1), NewRay(Point(0, 0, -5), Vector(0, 0, 1)), nil)

	if ii.OverPoint.Z >= -Epsilon/2 || ii.Point.Z <= ii.OverPoint.Z {
		t.Errorf("over point failed: p=%+v, o=%+v", ii.Point, ii.OverPoint)
	}

	if ii.UnderPoint.Z <= -Epsilon/2 || ii.Point.Z > ii.UnderPoint.Z {
		t.Errorf("under point failed: p=%+v, u=%+v", ii.Point, ii.UnderPoint)
	}
}

func getHittableAt(world *World, index int) Hittable {
	object := world.Objects[index]

	switch t := object.(type) {
	case Hittable:
		return t
	}

	return nil
}

func TestWorldShadeHit(t *testing.T) {
	w := createDefaultWorld()
	r := NewRay(Point(0, 0, -5), Vector(0, 0, 1))
	s := getHittableAt(w, 0)
	i := NewIntersection(4, s)
	ii := NewIntersectionInfo(i, r, nil)
	c := w.ShadeHit(ii, 0)

	if !c.Equals(RGB(0.38066, 0.47583, 0.2855)) {
		t.Errorf("drats! shade hit failed! %+v", c)
	}

	w.Lights[0] = NewPointLight(Point(0, 0.25, 0), RGB(1, 1, 1))
	r = NewRay(Point(0, 0, 0), Vector(0, 0, 1))
	s = getHittableAt(w, 1)
	i = NewIntersection(0.5, s)
	ii = NewIntersectionInfo(i, r, nil)
	c = w.ShadeHit(ii, 0)

	if !c.Equals(RGB(0.90498, 0.90498, 0.90498)) {
		t.Errorf("oh no! inside shade hit failed, ii=%+v, c=%+v", ii, c)
	}
}

func TestWorldColorAt(t *testing.T) {
	w := createDefaultWorld()
	r := NewRay(Point(0, 0, -5), Vector(0, 1, 0))
	c := w.ColorAt(r, 0)

	if !c.Equals(RGB(0, 0, 0)) {
		t.Errorf("color at (miss) should be black")
	}

	r = NewRay(Point(0, 0, -5), Vector(0, 0, 1))
	c = w.ColorAt(r, 0)

	if !c.Equals(RGB(0.38066, 0.47583, 0.2855)) {
		t.Errorf("color at (hit) failed")
	}

	outer := getHittableAt(w, 0)
	outer.Material().SetAmbient(1)
	inner := getHittableAt(w, 1)
	inner.Material().SetAmbient(1)
	r = NewRay(Point(0, 0, 0.75), Vector(0, 0, -1))
	c = w.ColorAt(r, 0)

	if !c.Equals(White.Blend(w.Ambient) /* should be inner.Material().Color but we cannot inspect that */) {
		t.Errorf("color at (intersection behind ray) failed: %+v", c)
	}
}

func TestWorldCameraRender(t *testing.T) {
	w := createDefaultWorld()
	c := NewCamera(11, 11, Pi/2)
	from := Point(0, 0, -5)
	to := Point(0, 0, 0)
	up := Vector(0, 1, 0)
	c.SetTransform(EyeViewpoint(from, to, up))
	canvas := w.RenderToCanvas(c)

	if !canvas.FastPixelAt(5, 5).Equals(RGB(0.38066, 0.47583, 0.2855)) {
		t.Errorf("render failed!")
	}
}

func TestWorldScene(t *testing.T) {
	TestWithImage(t)

	newMaterial := func(c Color) *Material {
		return NewMaterial().SetDiffuseColor(c)
	}

	floor := NewSphere()
	floor.SetTransform(Scaling(10, 0.01, 10))
	floor.SetMaterial(newMaterial(RGB(1, 0.9, 0.9)).SetSpecular(0))

	leftWall := NewSphere()
	leftWall.SetTransform(Translation(0, 0, 5).Mul(RotationY(-Pi / 4)).Mul(RotationX(Pi / 2)).Mul(Scaling(10, 0.01, 10)))
	leftWall.SetMaterial(floor.Material())

	rightWall := NewSphere()
	rightWall.SetTransform(Translation(0, 0, 5).Mul(RotationY(Pi / 4)).Mul(RotationX(Pi / 2)).Mul(Scaling(10, 0.01, 10)))
	rightWall.SetMaterial(floor.Material())

	middle := NewSphere()
	middle.SetTransform(Translation(-0.5, 1, 0.5))
	middle.SetMaterial(newMaterial(RGB(0.1, 1, 0.5)).SetDiffuse(0.7).SetSpecular(0.3))

	right := NewSphere()
	right.SetTransform(Translation(1.5, 0.5, -0.5).Mul(Scaling(0.5, 0.5, 0.5)))
	right.SetMaterial(newMaterial(RGB(0.5, 1, 0.1)).SetDiffuse(0.7).SetSpecular(0.3))

	left := NewSphere()
	left.SetTransform(Translation(-1.5, 0.33, -0.75).Mul(Scaling(0.33, 0.33, 0.33)))
	left.SetMaterial(newMaterial(RGB(1, 0.8, 0.1)).SetDiffuse(0.7).SetSpecular(0.3))

	light := NewPointLight(Point(-10, 10, -10), RGB(1, 1, 1))

	world := NewWorld()

	world.AddObjects(floor, leftWall, rightWall, middle, left, right)

	world.AddLights(light)

	world.ErpCanvasToImage = ErpLinear

	camera := NewCamera(640, 320, Pi/3)
	camera.SetTransform(EyeViewpoint(Point(0, 1.5, -5), Point(0, 1, 0), Vector(0, 1, 0)))

	world.RenderToPNG(camera, "test_simple_scene.png")
}

func TestWorldShadow(t *testing.T) {
	w := createDefaultWorld()
	rt := NewRaytracer(w)

	light := w.Lights[0].(*PointLight)

	if IsShadowed(light.Pos, rt, Point(0, 10, 0)) {
		t.Errorf("no shadow expected")
	}

	if !IsShadowed(light.Pos, rt, Point(10, -10, 10)) {
		t.Errorf("shadow expected")
	}

	if IsShadowed(light.Pos, rt, Point(-20, 20, 20)) {
		t.Errorf("no shadow expected (2)")
	}

	if IsShadowed(light.Pos, rt, Point(-2, 2, -2)) {
		t.Errorf("no shadow expected (3)")
	}
}

func TestWorldReflect(t *testing.T) {
	w := createDefaultWorld()

	r := NewRay(Point(0, 0, 0), Vector(0, 0, 1))
	s := getHittableAt(w, 1)
	s.Material().SetAmbient(1)
	i := NewIntersection(1, s)
	ii := NewIntersectionInfo(i, r, nil)
	c := w.ReflectedColor(ii, 0)

	if !c.Equals(Black) {
		t.Errorf("world reflect (non-reflective) failed")
	}

	w = createDefaultWorld()

	p := NewPlane()
	p.SetTransform(Translation(0, -1, 0))
	p.Material().SetReflective(0.5)

	w.AddObjects(p)

	r = NewRay(Point(0, 0, -3), Vector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	i = NewIntersection(math.Sqrt(2), p)
	ii = NewIntersectionInfo(i, r, nil)
	c = w.ReflectedColor(ii, 1)

	if !c.Equals(RGB(0.19033, 0.23791, 0.14274)) {
		t.Errorf("world reflect failed: %+v", c)
	}

	c = w.ShadeHit(ii, 1)

	if !c.Equals(RGB(0.87676, 0.92434, 0.82917)) {
		t.Errorf("world shade hit failed: %+v", c)
	}

	c = w.ReflectedColor(ii, 0)

	if !c.Equals(Black) {
		t.Errorf("reflection at 0-depth failed: %+v", c)
	}
}

func TestWorldReflectInfinite(t *testing.T) {
	w := NewWorld()

	w.AddLights(NewPointLight(Point(0, 0, 0), RGB(1, 1, 1)))

	lower := NewPlane()
	lower.Material().SetReflective(1)
	lower.SetTransform(Translation(0, -1, 0))

	upper := NewPlane()
	upper.Material().SetReflective(1)
	upper.SetTransform(Translation(0, +1, 0))

	w.AddObjects(lower, upper)

	r := NewRay(Point(0, 0, 0), Vector(0, 1, 0))

	w.ColorAt(r, 4) // Will crash if recursion is not handled correctly
}

func TestWorldRefract(t *testing.T) {
	w := createDefaultWorld()

	r := NewRay(Point(0, 0, -5), Vector(0, 0, 1))
	s := getHittableAt(w, 1)

	xs := NewIntersections()
	xs.Add(s, 4, 6)

	ii := NewIntersectionInfo(xs.At(0), r, xs)
	c := w.RefractedColor(ii, 5)

	if !c.Equals(Black) {
		t.Errorf("refracted color (opaque) failed")
	}

	// Now add some transparency but test with depth=0
	s.Material().SetRefractive(1, 1.5)

	ii = NewIntersectionInfo(xs.At(0), r, xs)

	c = w.RefractedColor(ii, 0)

	if !c.Equals(Black) {
		t.Errorf("refracted color (maximum depth reached) failed")
	}

	// Expecting black because of total internal reflection
	r = NewRay(Point(0, 0, math.Sqrt(2)/2), Vector(0, 1, 0))
	xs = NewIntersections()
	xs.Add(s, -math.Sqrt(2)/2, math.Sqrt(2)/2)
	ii = NewIntersectionInfo(xs.At(1), r, xs)

	c = w.RefractedColor(ii, 5)

	if !c.Equals(Black) {
		t.Errorf("refracted color (total internal reflection) failed")
	}

	// Actual refraction
	w = createDefaultWorld()
	w.Ambient = White
	a := getHittableAt(w, 0)
	a.Material().SetAmbient(1)
	a.Material().SetPattern(NewPointPattern())

	b := getHittableAt(w, 1)
	b.Material().SetRefractive(1, 1.5)

	r = NewRay(Point(0, 0, 0.1), Vector(0, 1, 0))
	xs = NewIntersections()
	xs.Add(a, -0.9899, 0.9899)
	xs.Add(b, -0.4899, 0.4899)
	xs.Sort()

	ii = NewIntersectionInfo(xs.At(2), r, xs)
	c = w.RefractedColor(ii, 5)

	if !c.Equals(RGB(0, 0.99888, 0.04722)) {
		t.Errorf("refracted color failed: %+v", c)
	}
}

func TestWorldShadeHitWithRefract(t *testing.T) {
	w := createDefaultWorld()

	floor := NewPlane()
	floor.SetTransform(Translation(0, -1, 0))
	floor.Material().SetRefractive(0.5, 1.5)

	ball := NewSphere()
	ball.SetTransform(Translation(0, -3.5, -0.5))
	ball.Material().SetDiffuseColor(CSS("red")).SetAmbient(0.5)

	w.AddObjects(floor, ball)

	r := NewRay(Point(0, 0, -3), Vector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := NewIntersections().Add(floor, math.Sqrt(2))
	ii := NewIntersectionInfo(xs.At(0), r, xs)

	c := w.ShadeHit(ii, 5)

	if !c.Equals(RGB(0.93642, 0.68642, 0.68642)) {
		t.Errorf("shade hit with refract failed: %+v", c)
	}
}

func TestWorldSimpleRefractScene(t *testing.T) {
	TestWithImage(t)

	world := NewWorld()

	floor := NewPlane()
	pattern := NewCheckerPattern(Black, White)
	pattern.SetTransform(Translation(0, 0.1, 0))
	floor.Material().SetPattern(pattern)
	floor.SetTransform(Translation(0, -10.1, 0))

	sphere := NewSphere()
	sphere.Material().SetDiffuse(0.1)
	sphere.Material().SetShininess(300)
	sphere.Material().SetReflective(1)
	sphere.Material().SetRefractive(1, 1.52)

	sphere2 := NewSphere()
	sphere2.Material().SetDiffuse(0.1)
	sphere2.Material().SetShininess(300)
	sphere2.Material().SetReflective(1)
	sphere2.Material().SetRefractive(1, 1) // Air
	sphere2.SetTransform(Scaling(0.5, 0.5, 0.5))

	world.AddObjects(floor, sphere, sphere2)

	light := NewPointLight(Point(20, 10, 0), RGB(0.7, 0.7, 0.7))

	world.AddLights(light)

	camera := NewCamera(480, 480, Pi/3)
	camera.SetTransform(EyeViewpoint(Point(0, 2.5, 0), Point(0, 0, 0), Vector(1, 0, 0)))

	world.ErpCanvasToImage = ErpLinear
	world.RenderToPNG(camera, "test_simple_refract_ball.png")
}

func TestWorldSchlick(t *testing.T) {
	s := NewSphere()
	s.Material().SetRefractive(1, 1.5)
	r := NewRay(Point(0, 0, math.Sqrt(2)/2), Vector(0, 1, 0))
	xs := NewIntersections().Add(s, -math.Sqrt(2)/2, math.Sqrt(2)/2)
	ii := NewIntersectionInfo(xs.At(1), r, xs)

	reflectance := SchlickReflectance(ii)

	if !FloatEqual(reflectance, 1) {
		t.Errorf("reflectance with diagonal ray failed")
	}

	r = NewRay(Point(0, 0, 0), Vector(0, 1, 0))
	xs = NewIntersections().Add(s, -1, +1)
	ii = NewIntersectionInfo(xs.At(1), r, xs)

	reflectance = SchlickReflectance(ii)

	if !FloatEqual(reflectance, 0.04) {
		t.Errorf("reflectance with perpendicular ray failed")
	}

	r = NewRay(Point(0, 0.99, -2), Vector(0, 0, 1))
	xs = NewIntersections().Add(s, 1.8589)
	ii = NewIntersectionInfo(xs.At(0), r, xs)

	reflectance = SchlickReflectance(ii)

	if !FloatEqual(reflectance, 0.48873) {
		t.Errorf("reflectance with small angle and n2>n1 failed")
	}
}
