// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package engine

import (
	"math"
	"testing"

	. "ascottix/funtracer/maths"
)

const Pi2 = Pi / 2

func TestCameraNew(t *testing.T) {
	c := NewCamera(160, 120, Pi2)

	if c.HSize != 160 || c.VSize != 120 || c.FOV != Pi/2 || !c.Transform().Equals(Identity()) {
		t.Errorf("camera new failed")
	}
}

func TestCameraPixsize(t *testing.T) {
	c1 := NewCamera(200, 125, Pi2)

	if !FloatEqual(c1.pixsize, 0.01) {
		t.Errorf("camera pixel size h>v failed")
	}

	c2 := NewCamera(125, 200, Pi2)

	if !FloatEqual(c2.pixsize, 0.01) {
		t.Errorf("camera pixel size h<v failed %+v", c2)
	}
}

func TestCameraRays(t *testing.T) {
	c := NewCamera(201, 101, Pi2)

	r := c.RayForPixel(100.5, 50.5)
	if !r.Origin.Equals(Point(0, 0, 0)) || !r.Direction.Equals(Vector(0, 0, -1)) {
		t.Error("ray for pixel 1 failed")
	}

	r = c.RayForPixel(0.5, 0.5)
	if !r.Origin.Equals(Point(0, 0, 0)) || !r.Direction.Equals(Vector(0.6651864, 0.33259321, -0.6685124)) {
		t.Error("ray for pixel 2 failed")
	}

	c.SetTransform(RotationY(Pi / 4).Mul(Translation(0, -2, 5)))
	r = c.RayForPixel(100.5, 50.5)
	if !r.Origin.Equals(Point(0, 2, -5)) || !r.Direction.Equals(Vector(math.Sqrt(2)/2, 0, -math.Sqrt(2)/2)) {
		t.Error("ray for pixel 2 failed")
	}
}
