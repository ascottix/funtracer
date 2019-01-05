// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"testing"
)

// Color
func TestColorAdd(t *testing.T) {
	c1 := RGB(0.9, 0.6, 0.75)
	c2 := RGB(0.7, 0.1, 0.25)

	if !c1.Add(c2).Equals(RGB(1.6, 0.7, 1.0)) {
		t.Errorf("color add failed")
	}
}

func TestColorSub(t *testing.T) {
	c1 := RGB(0.9, 0.6, 0.75)
	c2 := RGB(0.1, 0.1, 0.25)

	if !c1.Sub(c2).Equals(RGB(0.8, 0.5, 0.5)) {
		t.Errorf("color sub failed: %+v - %+v = %+v", c1, c2, c1.Sub(c2))
	}
}

func TestColorMul(t *testing.T) {
	c1 := RGB(0.2, 0.3, 0.4)

	if !c1.Mul(2).Equals(RGB(0.4, 0.6, 0.8)) {
		t.Errorf("color mul failed")
	}
}

func TestColorBlend(t *testing.T) {
	c1 := RGB(1, 0.2, 0.3)
	c2 := RGB(0.9, 1, 0.1)

	if !c1.Blend(c2).Equals(RGB(0.9, 0.2, 0.03)) {
		t.Errorf("color blend failed: %+v - %+v = %+v", c1, c2, c1.Blend(c2))
	}
}

// Canvas
func TestNewCanvas(t *testing.T) {
	c := NewCanvas(10, 20)

	if c.Width != 10 {
		t.Errorf("bad new canvas width")
	}

	if c.Height != 20 {
		t.Errorf("bad new canvas height")
	}

	if !c.PixelAt(5, 3).Equals(RGB(0, 0, 0)) {
		t.Errorf("new canvas pixels should be black")
	}
}

func TestCanvasWrite(t *testing.T) {
	c := NewCanvas(10, 20)
	r := RGB(1, 0, 0)

	if !c.PixelAt(2, 3).Equals(RGB(0, 0, 0)) {
		t.Errorf("new canvas pixels should be black")
	}

	c.SetPixelAt(2, 3, r)
	if !c.PixelAt(2, 3).Equals(r) {
		t.Errorf("cannot set pixel")
	}
}

func TestCanvasToPPM(t *testing.T) {
	t.SkipNow()

	c := NewCanvas(10, 2)

	c.Fill(RGB(1, 0.8, 0.6))

	c.WriteAsPPM(os.Stdout)
}

func TestCanvasDraw(t *testing.T) {
	t.SkipNow()

	// Projectile
	p_position := Point(0, 1, 0)
	p_velocity := Vector(1, 1.8, 0).Normalize().Mul(11.25)

	// Environment
	e_gravity := Vector(0, -0.1, 0)
	e_wind := Vector(-0.01, 0, 0)

	// Canvas
	c := NewCanvas(900, 500)

	r := RGB(1, 0.5, 0.5)

	// Plot graph
	tick := func() {
		c.SetPixelAt(p_position.X, float64(c.Height)-p_position.Y, r)
		p_position = p_position.Add(p_velocity)
		p_velocity = p_velocity.Add(e_gravity).Add(e_wind)
	}

	for p_position.Y > 0 {
		tick()
	}

	c.WriteAsPPM(os.Stdout)
}
