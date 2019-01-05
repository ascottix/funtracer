// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
)

type Camera struct {
	Transformer
	HSize      int
	VSize      int
	FOV        float64
	aspect     float64
	halfwidth  float64
	halfheight float64
	pixsize    float64
}

func NewCamera(hsize, vsize int, fov float64) *Camera {
	c := &Camera{HSize: hsize, VSize: vsize}

	c.SetTransform()
	c.SetFieldOfView(fov) // Triggers initialization of other parameters

	return c
}

func (c *Camera) SetFieldOfView(fov float64) {
	c.FOV = fov
	c.SetViewSize(c.HSize, c.VSize)
}

func (c *Camera) SetViewSize(hsize, vsize int) {
	c.HSize = hsize
	c.VSize = vsize
	c.aspect = float64(hsize) / float64(vsize)

	halfview := math.Tan(c.FOV / 2)

	if c.aspect >= 1 {
		c.halfwidth = halfview
		c.halfheight = halfview / c.aspect
	} else {
		c.halfwidth = halfview * c.aspect
		c.halfheight = halfview
	}

	c.pixsize = (c.halfwidth * 2) / float64(hsize) // Assume square pixel
}

func (c *Camera) RayForPixelI(x, y int) Ray {
	return c.RayForPixelF(float64(x)+0.5, float64(y)+0.5)
}

func (c *Camera) RayForPixelF(x, y float64) Ray {
	// Offset of the pixel center from edge of canvas
	xoffset := x * c.pixsize
	yoffset := y * c.pixsize

	// Note: camera always look toward -z, which places +x to the left

	// Get untransformed world coordinates
	worldx := c.halfwidth - xoffset
	worldy := c.halfheight - yoffset

	// Apply the transformation
	pixel := c.Tinverse.MulT(Point(worldx, worldy, -1))
	origin := c.Tinverse.MulT(Point(0, 0, 0))
	direction := pixel.Sub(origin).Normalize()

	return Ray{origin, direction}
}
