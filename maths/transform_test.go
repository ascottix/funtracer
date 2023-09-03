// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package maths

import (
	"math"
	"testing"
)

func TestTranslation(t *testing.T) {
	a := Translation(5, -3, 2)

	if !a.MulT(Point(-3, 4, 5)).Equals(Point(2, 1, 7)) {
		t.Errorf("point translation failed")
	}

	if !a.Inverse().MulT(Point(-3, 4, 5)).Equals(Point(-8, 7, 3)) {
		t.Errorf("point translation failed")
	}

	v := Vector(-3, 4, 5)

	if !a.MulT(v).Equals(v) {
		t.Errorf("vector translation failed")
	}
}

func TestScaling(t *testing.T) {
	a := Scaling(2, 3, 4)

	if !a.MulT(Point(-4, 6, 8)).Equals(Point(-8, 18, 32)) {
		t.Errorf("point scaling failed")
	}

	if !a.MulT(Vector(-4, 6, 8)).Equals(Vector(-8, 18, 32)) {
		t.Errorf("vector scaling failed")
	}

	if !a.Inverse().MulT(Vector(-4, 6, 8)).Equals(Vector(-2, 2, 2)) {
		t.Errorf("vector scaling failed")
	}

	if !Scaling(-1, 1, 1).MulT(Point(2, 3, 4)).Equals(Point(-2, 3, 4)) {
		t.Errorf("point x reflection failed")
	}
}

func TestRotation(t *testing.T) {
	px := Point(1, 0, 0)
	py := Point(0, 1, 0)
	pz := Point(0, 0, 1)
	h2 := math.Sqrt(2) / 2

	if !RotationX(Pi / 4).MulT(py).Equals(Point(0, h2, h2)) {
		t.Errorf("x rotation(pi/4) failed")
	}

	if !RotationX(Pi / 2).MulT(py).Equals(pz) {
		t.Errorf("x rotation(pi/2) failed")
	}

	if !RotationX(Pi / 4).Inverse().MulT(py).Equals(Point(0, h2, -h2)) {
		t.Errorf("x inverse rotation(pi/4) failed")
	}

	if !RotationY(Pi / 4).MulT(pz).Equals(Point(h2, 0, h2)) {
		t.Errorf("y rotation(pi/4) failed")
	}

	if !RotationY(Pi / 2).MulT(pz).Equals(px) {
		t.Errorf("y rotation(pi/2) failed")
	}

	if !RotationZ(Pi / 4).MulT(py).Equals(Point(-h2, h2, 0)) {
		t.Errorf("z rotation(pi/4) failed")
	}

	if !RotationZ(Pi / 2).MulT(py).Equals(Point(-1, 0, 0)) {
		t.Errorf("z rotation(pi/2) failed")
	}
}

func TestShearing(t *testing.T) {
	p := Point(2, 3, 4)

	if !Shearing(1, 0, 0, 0, 0, 0).MulT(p).Equals(Point(5, 3, 4)) {
		t.Errorf("shearing xy failed")
	}

	if !Shearing(0, 1, 0, 0, 0, 0).MulT(p).Equals(Point(6, 3, 4)) {
		t.Errorf("shearing xz failed")
	}

	if !Shearing(0, 0, 1, 0, 0, 0).MulT(p).Equals(Point(2, 5, 4)) {
		t.Errorf("shearing yx failed")
	}

	if !Shearing(0, 0, 0, 1, 0, 0).MulT(p).Equals(Point(2, 7, 4)) {
		t.Errorf("shearing yz failed")
	}

	if !Shearing(0, 0, 0, 0, 1, 0).MulT(p).Equals(Point(2, 3, 6)) {
		t.Errorf("shearing zx failed")
	}

	if !Shearing(0, 0, 0, 0, 0, 1).MulT(p).Equals(Point(2, 3, 7)) {
		t.Errorf("shearing zy failed")
	}
}

func TestComposition(t *testing.T) {
	p1 := Point(1, 0, 1)
	a := RotationX(Pi / 2)
	b := Scaling(5, 5, 5)
	c := Translation(10, 5, 7)

	p2 := a.MulT(p1)
	if !p2.Equals(Point(1, -1, 0)) {
		t.Errorf("transform failed")
	}

	p3 := b.MulT(p2)
	if !p3.Equals(Point(5, -5, 0)) {
		t.Errorf("transform failed")
	}

	p4 := c.MulT(p3)
	if !p4.Equals(Point(15, 0, 7)) {
		t.Errorf("transform failed")
	}

	if !c.Mul(b).Mul(a).MulT(p1).Equals(p4) {
		t.Errorf("chained transform failed")
	}
}

// TODO!!!
// func TestAnalogClock(t *testing.T) {
// 	t.SkipNow()

// 	// Canvas
// 	c := NewCanvas(200, 200)

// 	w := RGB(1, 1, 1)

// 	// Plot clock hours
// 	p := Point(0, 1, 0) // 12th hour
// 	s := Scaling(80, 80, 0)
// 	x := Translation(100, 100, 0)
// 	for hour := 0; hour < 12; hour++ {
// 		r := RotationZ(2 * Pi * float64(hour) / 12)
// 		h := x.MulT(r.Mul(s).MulT(p))

// 		c.SetPixelAt(h.X, h.Y, w)
// 	}

// 	c.WriteAsPPM(os.Stdout)
// }

func TestEyeViewpoint(t *testing.T) {
	if !EyeViewpoint(Point(0, 0, 0), Point(0, 0, -1), Vector(0, 1, 0)).Equals(Identity()) {
		t.Errorf("view transform failed")
	}

	if !EyeViewpoint(Point(0, 0, 0), Point(0, 0, 1), Vector(0, 1, 0)).Equals(Scaling(-1, 1, -1)) {
		t.Errorf("view transform (positive z) failed")
	}

	if !EyeViewpoint(Point(0, 0, 8), Point(0, 0, 0), Vector(0, 1, 0)).Equals(Translation(0, 0, -8)) {
		t.Errorf("view transform (world moves) failed")
	}

	m := NewMatrix(4, 4,
		-0.50709255, 0.50709255, 0.6761234, -2.36643191,
		0.7677159, 0.60609153, 0.1212183, -2.8284271,
		-0.3585686, 0.5976143, -0.717137166, 0.00000,
		0, 0, 0, 1)

	e := EyeViewpoint(Point(1, 3, 2), Point(4, -2, 8), Vector(1, 1, 0))

	if !e.Equals(m) {
		t.Errorf("view transform (all) failed: %+v", e)
	}
}
