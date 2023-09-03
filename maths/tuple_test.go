// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package maths

import (
	"math"
	"testing" // https://golang.org/pkg/testing/
)

func claim(t *testing.T, cond bool) {
	if !cond {
		t.Fail()
	}
}

func TestTuplePoint(t *testing.T) {
	a := Tuple{4.3, -4.2, 3.1, 1.0}

	claim(t, a.X == 4.3)
	claim(t, a.Y == -4.2)
	claim(t, a.Z == 3.1)
	claim(t, a.W == 1.0)

	if !a.IsPoint() {
		t.Errorf("a should be a point")
	}

	if a.IsVector() {
		t.Errorf("a should not be a vector")
	}

	x, y, z := 4.0, -4.0, 3.0

	if !Point(x, y, z).Equals(Tuple{x, y, z, 1}) {
		t.Errorf("w should be 1 for a point")
	}
}

func TestTupleVector(t *testing.T) {
	a := Tuple{4.3, -4.2, 3.1, 0.0}

	claim(t, a.X == 4.3)
	claim(t, a.Y == -4.2)
	claim(t, a.Z == 3.1)
	claim(t, a.W == 0.0)

	if a.IsPoint() {
		t.Errorf("a should not be a point")
	}

	if !a.IsVector() {
		t.Errorf("a should be a vector")
	}

	x, y, z := 4.0, -4.0, 3.0

	if !Vector(x, y, z).Equals(Tuple{x, y, z, 0}) {
		t.Errorf("w should be 0 for a vector")
	}
}

func TestTupleAdd(t *testing.T) {
	a1 := Point(3, -2, 5)
	a2 := Vector(-2, 3, 1)
	r := Tuple{1, 1, 6, 1}

	if !a1.Add(a2).Equals(r) {
		t.Errorf("%+v + %+v should be %+v", a1, a2, r)
	}
}

func TestTupleSub(t *testing.T) {
	p1 := Point(3, 2, 1)
	p2 := Point(5, 6, 7)
	r1 := Vector(-2, -4, -6)

	if !p1.Sub(p2).Equals(r1) {
		t.Errorf("%+v - %+v should be %+v", p1, p2, r1)
	}

	v1 := Vector(3, 2, 1)
	v2 := Vector(5, 6, 7)
	r2 := Point(-2, -4, -6)

	if !p1.Sub(v2).Equals(r2) {
		t.Errorf("%+v - %+v should be %+v", p1, v2, r2)
	}

	if !v1.Sub(v2).Equals(r1) {
		t.Errorf("%+v - %+v should be %+v", v1, v2, r1)
	}
}

func TestTupleNeg(t *testing.T) {
	v := Vector(1, -2, 3)
	z := Vector(0, 0, 0)

	if !z.Sub(v).Equals(Vector(-1, 2, -3)) {
		t.Errorf("%+v - %+v failed", z, v)
	}

	a := Tuple{1, -2, 3, -4}
	r := Tuple{-1, 2, -3, 4}

	if !a.Neg().Equals(r) {
		t.Errorf("negating %+v should yield %+v", a, r)
	}
}

func TestTupleMul(t *testing.T) {
	a := Tuple{1, -2, 3, -4}
	f := 3.5

	if !a.Mul(f).Equals(Tuple{3.5, -7, 10.5, -14}) {
		t.Errorf("%+v mul by %f failed", a, f)
	}

	h := 0.5
	r := Tuple{0.5, -1, 1.5, -2}

	if !a.Mul(h).Equals(r) {
		t.Errorf("%+v mul by %f failed", a, h)
	}

	d := 2.0

	if !a.Div(d).Equals(r) {
		t.Errorf("%+v div by %f failed", a, d)
	}
}

func TestTupleMagnitude(t *testing.T) {
	if Vector(1, 0, 0).Length() != 1 {
		t.Errorf("bad length for (1,0,0)")
	}

	if Vector(0, 1, 0).Length() != 1 {
		t.Errorf("bad length for (0,1,0)")
	}

	if Vector(0, 0, 1).Length() != 1 {
		t.Errorf("bad length for (0,0,1)")
	}

	if Vector(1, 2, 3).Length() != math.Sqrt(14) {
		t.Errorf("bad length for (1,2,3)")
	}

	if Vector(-1, -2, -3).Length() != math.Sqrt(14) {
		t.Errorf("bad length for (-1,-2,-3)")
	}
}

func TestTupleNormalization(t *testing.T) {
	if !Vector(4, 0, 0).Normalize().Equals(Vector(1, 0, 0)) {
		t.Errorf("bad normalization for (4,0,0)")
	}

	v := Vector(1, 2, 3)
	u := math.Sqrt(14)
	n := Vector(1/u, 2/u, 3/u)

	if !v.Normalize().ApproxEquals(n) {
		t.Errorf("bad normalization for %+v: %+v should be %+v", v, v.Normalize(), n)
	}

	if !FloatEqual(v.Normalize().Length(), 1) {
		t.Errorf("bad normalization length for %+v: %f", v, v.Normalize().Length())
	}
}

func TestTupleDotProduct(t *testing.T) {
	if Vector(1, 2, 3).DotProduct(Vector(2, 3, 4)) != 20 {
		t.Errorf("dot product failed")
	}
}

func TestTupleCrossProduct(t *testing.T) {
	a := Vector(1, 2, 3)
	b := Vector(2, 3, 4)

	if !a.CrossProduct(b).Equals(Vector(-1, 2, -1)) {
		t.Errorf("cross product (a,b) failed: %+v", a.CrossProduct(b))
	}

	if !b.CrossProduct(a).Equals(Vector(1, -2, 1)) {
		t.Errorf("cross product (b,a) failed: %+v", b.CrossProduct(a))
	}
}

func TestTupleReflect(t *testing.T) {
	if !Vector(1, -1, 0).Reflect(Vector(0, 1, 0)).Equals(Vector(1, 1, 0)) {
		t.Errorf("reflect 1 failed")
	}

	if !Vector(0, -1, 0).Reflect(Vector(math.Sqrt(2)/2, math.Sqrt(2)/2, 0)).Equals(Vector(1, 0, 0)) {
		t.Errorf("reflect 2 failed")
	}
}
