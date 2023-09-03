// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package engine

import (
	"math"

	. "ascottix/funtracer/maths"
	. "ascottix/funtracer/textures"
	. "ascottix/funtracer/traits"
)

type Camera struct {
	Transformer
	HSize      int     // Image width in pixel
	VSize      int     // Image height in pixel
	FOV        float64 // Field of view in radians
	aspect     float64 // Aspect ratio
	halfwidth  float64 // Half width of projected image
	halfheight float64 // Half height of projected image
	pixsize    float64 // Size of one pixel
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

func (c *Camera) RayForPixel(x, y float64) Ray {
	// Offset of the pixel center from edge of canvas
	xoffset := x * c.pixsize
	yoffset := y * c.pixsize

	// Note: camera is placed at (0,0,0) and looks toward -z, with +x to the left

	// Get untransformed world coordinates
	worldx := c.halfwidth - xoffset
	worldy := c.halfheight - yoffset

	// Apply the transformation
	pixel := c.Tinverse.MulT(Point(worldx, worldy, -1))
	origin := c.Tinverse.MulT(Point(0, 0, 0))
	direction := pixel.Sub(origin).Normalize()

	return Ray{Origin: origin, Direction: direction}
}

// ConcentricSampleDisk converts samples from [0,1)x[0,1) into
// a 2D unit disk centered at the origin (0,0)
func ConcentricSampleDisk(u, v float64) (float64, float64) {
	// Get two random samples in the [-1,+1] interval
	x := u*2 - 1
	y := v*2 - 1

	if x != 0 || y != 0 {
		// Apply concentric mapping to point
		var r, theta float64

		if math.Abs(x) > math.Abs(y) {
			r = x
			theta = (y / x) * Pi / 4
		} else {
			r = y
			theta = Pi/2 - (x/y)*Pi/4
		}

		x = r * math.Cos(theta)
		y = r * math.Sin(theta)
	}

	return x, y
}

func (c *Camera) RayForPixelDepthOfField(x, y, lensRadius, focalDistance float64, rand FloatGenerator) Ray {
	// Displace origin to a random point on the lens
	lensX, lensY := ConcentricSampleDisk(rand(), rand())
	origin := Point(lensX*lensRadius, lensY*lensRadius, 0)

	// Get the target pixel on the camera plane (in world coordinates)
	pixel := Point(c.halfwidth-x*c.pixsize, c.halfheight-y*c.pixsize, -1)

	// Adjust the target pixel to account for focal distance
	ft := focalDistance / -pixel.Normalize().Z // How far the pixel is from the focal plane
	pFocus := pixel.Mul(ft)                    // Target point on the focal plane

	pixel = pFocus.Sub(origin)  // Direction from adjusted origin to target plane
	pixel = pixel.Div(-pixel.Z) // Place pixel on the z=-1 plane
	pixel = pixel.Add(origin)   // Add origin back
	pixel.W = 1                 // Make sure we have a point

	// Build ray from transformed endpoints
	pixel = c.Tinverse.MulT(pixel)
	origin = c.Tinverse.MulT(origin)
	direction := pixel.Sub(origin).Normalize()

	return Ray{Origin: origin, Direction: direction}
}
