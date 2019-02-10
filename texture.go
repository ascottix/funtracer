// Copyright (c) 2019 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"os"
)

type TextureWrap int

const (
	TwPeriodic TextureWrap = iota
	TwClamp
	TwBlack
)

type ImageTexture struct {
	image image.Image
	wrap  TextureWrap
}

func NewImageTexture() *ImageTexture {
	return new(ImageTexture)
}

func (t *ImageTexture) Load(r io.Reader) error {
	image, _, err := image.Decode(r)

	if err == nil {
		t.image = image
	}

	return err
}

func (t *ImageTexture) LoadFromFile(filename string) error {
	f, err := os.Open(filename)

	if err == nil {
		err = t.Load(f)
		f.Close()
	}

	return err
}

func textureToImage(u float64, imgSize int) (i, j int, t float64) {
	w := float64(imgSize)
	h := 1 / (2 * w) // My own elucubration... 0 is fine

	switch {
	case u <= h:
		// All values are zero, so nothing to do
	case u >= (1 - h):
		i = imgSize - 1
		j = imgSize - 1
		t = 0
	default:
		s := w * (u - h)
		f := math.Floor(s)
		t = s - f
		i = int(f)
		j = i + 1

		if j >= imgSize {
			j--
		}
	}

	return
}

// ImageColorToTuple encodes an image color into a Tuple, using W to store alpha and X, Y, Z for R, G, B
func ImageColorToTuple(c color.Color) Tuple {
	const Max = 65535

	r, g, b, a := c.RGBA()

	t := Tuple{float64(r)/Max, float64(g)/Max, float64(b)/Max, float64(a)/Max}

	// TODO: need to consider actual color space
	t.X = ErpGammaToLinear(t.X)
	t.Y = ErpGammaToLinear(t.Y)
	t.Z = ErpGammaToLinear(t.Z)

	return t
}

func (t *ImageTexture) TextureAt(x, y int) Tuple {
	bounds := t.image.Bounds()

	return ImageColorToTuple( t.image.At(bounds.Min.X+x, bounds.Min.Y+y) )
}

func (t *ImageTexture) ApplyAtHit(ii *IntersectionInfo) {
	if t.image == nil {
		return
	}

	// Get surface coordinates
	u := ii.U
	v := ii.V

	// Handle wrapping
	if u < 0 || u > 1 || v < 0 || v > 1 {
		switch t.wrap {
		case TwBlack:
			ii.Mat.DiffuseColor = Black
			return
		case TwPeriodic:
			u -= math.Floor(u)
			v -= math.Floor(v)
		case TwClamp:
			u = Clamp(u)
			v = Clamp(v)
		}
	}

	// Need to consider that image may not start at (0,0)
	bounds := t.image.Bounds()

	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y

	// Bilinear filtering
	x0, x1, tu := textureToImage(u, w)
	y0, y1, tv := textureToImage(v, h)

	c00 := t.TextureAt(x0, y0).Mul((1 - tu) * (1 - tv))
	c01 := t.TextureAt(x0, y1).Mul((1 - tu) * tv)
	c10 := t.TextureAt(x1, y0).Mul(tu * (1 - tv))
	c11 := t.TextureAt(x1, y1).Mul(tu * tv)

	ct := c00.Add(c01).Add(c10).Add(c11)

	// Image color is alpha pre-multiplied, so we only need to merge in the target color
	ii.Mat.DiffuseColor = RGB(ct.X, ct.Y, ct.Z).Add(ii.Mat.DiffuseColor.Mul(1-ct.W))
}
