// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package maths

import (
	"testing" // https://golang.org/pkg/testing/
)

func TestNewMatrix(t *testing.T) {
	m := NewMatrix(4, 4,
		1, 2, 3, 4,
		5.5, 6.5, 7.5, 8.5,
		9, 10, 11, 12,
		13.5, 14.5, 15.5, 16.5)

	if m.At(0, 0) != 1 || m.At(0, 3) != 4 || m.At(1, 0) != 5.5 || m.At(1, 2) != 7.5 || m.At(2, 2) != 11 || m.At(3, 0) != 13.5 || m.At(3, 2) != 15.5 {
		t.Errorf("bad initialization %+v", m)
	}

	m2 := NewMatrix(2, 2,
		-3, 5,
		1, -2)

	if m2.At(0, 0) != -3 || m2.At(0, 1) != 5 || m2.At(1, 0) != 1 || m2.At(1, 1) != -2 {
		t.Errorf("bad initialization %+v", m2)
	}

	m3 := NewMatrix(3, 3,
		-3, 5, 0,
		1, -2, -7,
		0, 1, 1)

	if m3.At(0, 0) != -3 || m3.At(1, 1) != -2 || m3.At(2, 2) != 1 {
		t.Errorf("bad initialization %+v", m3)
	}
}

func TestEquals(t *testing.T) {
	a := NewMatrix(4, 4,
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16)
	b := NewMatrix(4, 4,
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16)
	c := NewMatrix(4, 4,
		2, 3, 4, 5,
		6, 7, 8, 9,
		8, 7, 6, 5,
		4, 3, 2, 1)

	if !a.Equals(b) {
		t.Errorf("expecting %+v = %+v", a, b)
	}

	if a.Equals(c) {
		t.Errorf("expecting %+v != %+v", a, c)
	}
}

func TestMul(t *testing.T) {
	a := NewMatrix(4, 4,
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 8, 7, 6,
		5, 4, 3, 2)
	b := NewMatrix(4, 4,
		-2, 1, 2, 3,
		3, 2, 1, -1,
		4, 3, 6, 5,
		1, 2, 7, 8)
	r := NewMatrix(4, 4,
		20, 22, 50, 48,
		44, 54, 114, 108,
		40, 58, 110, 102,
		16, 26, 46, 42)

	if !a.Mul(b).Equals(r) {
		t.Errorf("mul failed %+v", a.Mul(b))
	}
}

func TestMulT(t *testing.T) {
	a := NewMatrix(4, 4,
		1, 2, 3, 4,
		2, 4, 4, 2,
		8, 6, 4, 1,
		0, 0, 0, 1)
	b := Point(1, 2, 3)

	if !a.MulT(b).Equals(Point(18, 24, 33)) {
		t.Errorf("mul tuple failed %+v", a.MulT(b))
	}
}

func TestIdentity(t *testing.T) {
	a := NewMatrix(4, 4,
		-2, 1, 2, 3,
		3, 2, 1, -1,
		4, 3, 6, 5,
		1, 2, 7, 8)
	u := Tuple{1, 2, 3, 4}
	i := NewIdentityMatrix4x4()

	if !a.Mul(i).Equals(a) {
		t.Errorf("identity mul failed %+v", a.Mul(i))
	}

	if !i.MulT(u).Equals(u) {
		t.Errorf("identity mul tuple failed %+v", i.MulT(u))
	}
}

func TestTranspose(t *testing.T) {
	a := NewMatrix(4, 4,
		0, 9, 3, 0,
		9, 8, 0, 8,
		1, 8, 5, 3,
		0, 0, 5, 8)
	b := NewMatrix(4, 4,
		0, 9, 1, 0,
		9, 8, 8, 0,
		3, 0, 5, 5,
		0, 8, 3, 8)
	i := NewIdentityMatrix4x4()

	if !a.Transpose().Equals(b) {
		t.Errorf("transpose failed %+v", a.Transpose())
	}

	if !i.Transpose().Equals(i) {
		t.Errorf("identity transpose failed %+v", i.Transpose())
	}
}

func TestDeterminant(t *testing.T) {
	a := NewMatrix(2, 2,
		1, 5,
		-3, 2)

	if a.Determinant() != 17 {
		t.Errorf("2x2 determinant failed: %f", a.Determinant())
	}
}

func TestSubmatrix(t *testing.T) {
	a := NewMatrix(3, 3,
		1, 5, 0,
		-3, 2, 7,
		0, 6, -3)
	b := NewMatrix(2, 2,
		-3, 2,
		0, 6)
	c := NewMatrix(4, 4,
		-6, 1, 1, 6,
		-8, 5, 8, 6,
		-1, 0, 8, 2,
		-7, 1, -1, 1)
	d := NewMatrix(3, 3,
		-6, 1, 6,
		-8, 8, 6,
		-7, -1, 1)

	if !a.Submatrix(0, 2).Equals(b) {
		t.Errorf("submatrix failed %+v != %+v", a.Submatrix(0, 2), b)
	}

	if !c.Submatrix(2, 1).Equals(d) {
		t.Errorf("submatrix failed %+v != %+v", c.Submatrix(2, 1), d)
	}
}

func TestMinor(t *testing.T) {
	a := NewMatrix(3, 3,
		3, 5, 0,
		2, -1, -7,
		6, -1, 5)
	b := a.Submatrix(1, 0)

	if b.Determinant() != 25 || a.Minor(1, 0) != 25 {
		t.Errorf("minor failed")
	}
}

func TestCofactor(t *testing.T) {
	a := NewMatrix(3, 3,
		3, 5, 0,
		2, -1, -7,
		6, -1, 5)

	if a.Cofactor(0, 0) != -12 || a.Cofactor(0, 0) != a.Minor(0, 0) {
		t.Errorf("b.cofactor(0,0) failed")
	}

	if a.Cofactor(1, 0) != -25 || a.Cofactor(1, 0) != -a.Minor(1, 0) {
		t.Errorf("b.cofactor(1,0) failed")
	}

	b := NewMatrix(3, 3,
		1, 2, 6,
		-5, 8, -4,
		2, 6, 4)

	if b.Cofactor(0, 0) != 56 || b.Cofactor(0, 1) != 12 || b.Cofactor(0, 2) != -46 || b.Determinant() != -196 {
		t.Errorf("b.cofactor/determinant 3x3 failed")
	}

	c := NewMatrix(4, 4,
		-2, -8, 3, 5,
		-3, 1, 7, 3,
		1, 2, -9, 6,
		-6, 7, 7, -9)

	if c.Cofactor(0, 0) != 690 || c.Cofactor(0, 1) != 447 || c.Cofactor(0, 2) != 210 || c.Cofactor(0, 3) != 51 || c.Determinant() != -4071 {
		t.Errorf("b.cofactor/determinant 4x4 failed")
	}
}

func TestIsInvertible(t *testing.T) {
	a := NewMatrix(4, 4,
		6, 4, 4, 4,
		5, 5, 7, 6,
		4, -9, 3, -7,
		9, 1, 7, -6)

	if a.Determinant() != -2120 || !a.IsInvertible() {
		t.Errorf("a invertible failed")
	}

	b := NewMatrix(4, 4,
		-4, 2, -2, -3,
		9, 6, 2, 6,
		0, -5, 1, -5,
		0, 0, 0, 0)

	if b.Determinant() != 0 || b.IsInvertible() {
		t.Errorf("b invertible failed")
	}
}

func TestInverse(t *testing.T) {
	a := NewMatrix(4, 4,
		-5, 2, 6, -8,
		1, -5, 1, 8,
		7, 7, -6, -7,
		1, -3, 7, 4)

	b := a.Inverse()

	c := NewMatrix(4, 4,
		0.21804511, 0.4511278, 0.2406015, -0.04511278,
		-0.80827068, -1.4567669, -0.443609, 0.52067669,
		-0.0789474, -0.2236842, -0.0526316, 0.19736842,
		-0.5225564, -0.81390978, -0.30075188, 0.306390978,
	)

	if a.Determinant() != 532 || a.Cofactor(2, 3) != -160 || a.Cofactor(3, 2) != 105 {
		t.Errorf("a.cofactor/determinant failed")
	}

	if !b.Equals(c) {
		t.Errorf("inverse failed: %+v != %+v", b, c)
	}

	d := NewMatrix(4, 4,
		8, -5, 9, 2,
		7, 5, 6, 1,
		-6, 0, 9, 6,
		-3, 0, -9, -4)

	e := NewMatrix(4, 4,
		-0.15384615, -0.15384615, -0.28205128, -0.53846154,
		-0.076923077, 0.12307692, 0.025641026, 0.03076923,
		0.35897436, 0.35897436, 0.435897436, 0.92307692,
		-0.692307692, -0.692307692, -0.76923077, -1.92307692)

	if !d.Inverse().Equals(e) {
		t.Errorf("inverse failed: %+v != %+v", d.Inverse(), e)
	}

	f := NewMatrix(4, 4,
		9, 3, 0, 9,
		-5, -2, -6, -3,
		-4, 9, 6, 4,
		-7, 6, 6, 2)

	g := NewMatrix(4, 4,
		-0.04074074, -0.07777778, 0.14444444, -0.22222222,
		-0.07777778, 0.03333333, 0.36666667, -0.33333333,
		-0.02901234, -0.1462963, -0.10925926, 0.12962963,
		0.17777778, 0.06666667, -0.26666667, 0.33333333)

	if !f.Inverse().Equals(g) {
		t.Errorf("inverse failed: %+v", f.Inverse())
	}

	h := NewMatrix(4, 4,
		3, -9, 7, 3,
		3, -8, 2, -9,
		-4, 4, 4, 1,
		-6, 5, -1, 1)
	j := NewMatrix(4, 4,
		8, 2, 2, 2,
		3, -1, 7, 0,
		7, 0, 5, 4,
		6, -2, 0, 5)

	k := h.Mul(j)

	if !k.Mul(j.Inverse()).Equals(h) {
		t.Errorf("inverse verification failed: %+v", k)
	}

	i := NewIdentityMatrix4x4()

	if !i.Inverse().Equals(i) {
		t.Errorf("error inverting the identity matrix")
	}

	if !a.Mul(a.Inverse()).Equals(i) {
		t.Errorf("expecting identity, found %+v", a.Mul(a.Inverse()))
	}
}
