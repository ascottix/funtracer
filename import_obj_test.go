// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func TestGibberish(t *testing.T) {
	data := `
There was a lady named Bright
who could walk faster than light.
She set out one day
in a relative way,
and came back the previous night.    
	`
	ParseWavefrontObjFromString(data)
}

func TestVertex(t *testing.T) {
	data := `
v -1 1 0
v -1.0000 0.5000 0.0000 
v 1 0 0
v 1 1 0
	`
	info := ParseWavefrontObjFromString(data)

	v := info.Vertices

	if len(v) != 4 || !v[0].Equals(Point(-1, 1, 0)) || !v[1].Equals(Point(-1, 0.5, 0)) || !v[2].Equals(Point(1, 0, 0)) || !v[3].Equals(Point(1, 1, 0)) {
		t.Errorf("obj vertex failed")
	}
}

func TestTriangle(t *testing.T) {
	data := `
v -1 1 0 
v -1 0 0 
v 1 0 0 
v 1 1 0

f 1 2 3 
f 1 3 4
	`
	info := ParseWavefrontObjFromString(data)

	if len(info.Triangles) != 2 {
		t.Errorf("obj triangle failed")
		t.SkipNow()
	}

	t1 := info.Triangles[0]
	t2 := info.Triangles[1]

	v := info.Vertices

	if !t1.P1.Equals(v[0]) || !t1.P2.Equals(v[1]) || !t1.P3.Equals(v[2]) || !t2.P1.Equals(v[0]) || !t2.P2.Equals(v[2]) || !t2.P3.Equals(v[3]) {
		t.Errorf("obj triangle vertices mismatch")
	}
}

func TestPolygon(t *testing.T) {
	data := `
v -1 1 0 
v -1 0 0 
v 1 0 0 
v 1 1 0 
v 0 2 0
f 1 2 3 4 5
	`
	info := ParseWavefrontObjFromString(data)

	if len(info.Triangles) != 3 {
		t.Errorf("obj polygon failed")
		t.SkipNow()
	}

	t1 := info.Triangles[0]
	t2 := info.Triangles[1]
	t3 := info.Triangles[2]

	v := info.Vertices

	if !t1.P1.Equals(v[0]) || !t1.P2.Equals(v[1]) || !t1.P3.Equals(v[2]) ||
		!t2.P1.Equals(v[0]) || !t2.P2.Equals(v[2]) || !t2.P3.Equals(v[3]) ||
		!t3.P1.Equals(v[0]) || !t3.P2.Equals(v[3]) || !t3.P3.Equals(v[4]) {
		t.Errorf("obj polygon vertices mismatch")
	}
}

func TestGroup(t *testing.T) {
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
