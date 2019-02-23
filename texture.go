// Copyright (c) 2019 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"os"
)

type TextureWrap int

// Extension point to allow custom processing
type TextureOnMapUv func(u, v float64, ii *IntersectionInfo) (float64, float64) // Triggers after u and v are fetched, but before they are used
type TextureOnImage func(data []ColorRGBA, w, h int)
type TextureOnApply func(c ColorRGBA, ii *IntersectionInfo) // Triggers before color is applied to hit, replaces standard processing

// ColorRGBA is a RGB color with alpha, in linear space
type ColorRGBA struct {
	r, g, b, a float32
}

func (c ColorRGBA) Add(d ColorRGBA) ColorRGBA {
	return ColorRGBA{r: c.r + d.r, g: c.g + d.g, b: c.b + d.b, a: c.a + d.a}
}

func (c ColorRGBA) Mul(f float32) ColorRGBA {
	return ColorRGBA{r: c.r * f, g: c.g * f, b: c.b * f, a: c.a * f}
}

func (c ColorRGBA) RGB() Color {
	return RGB(float64(c.r), float64(c.g), float64(c.b))
}

const (
	TwPeriodic TextureWrap = iota
	TwFlip
	TwClamp
	TwBlack
)

type ImageTexture struct {
	data    []ColorRGBA
	w, h    int
	wrap    TextureWrap
	onImage TextureOnImage
	onMapUv TextureOnMapUv
	onApply TextureOnApply
	linear  bool
}

type NormalMap struct {
}

func NewImageTexture() *ImageTexture {
	return new(ImageTexture)
}

func (t *ImageTexture) Load(r io.Reader) error {
	image, _, err := image.Decode(r)

	if err == nil {
		bounds := image.Bounds()
		
		w := bounds.Max.X - bounds.Min.X
		h := bounds.Max.Y - bounds.Min.Y

		t.data = make([]ColorRGBA, w*h)
		t.w = w
		t.h = h

		// Convert the texture to ColorRGBA format
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				const Max = 65535

				var c ColorRGBA

				r, g, b, a := image.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()

				if t.linear {
					// Linear space is used for normal maps
					c = ColorRGBA{
						float32(r) / Max,
						float32(g) / Max,
						float32(b) / Max,
						float32(a) / Max,
					}
				} else {
					// sRGB is the color space for most JPEG and PNG images
					c = ColorRGBA{
						float32(ErpsRGBToLinear(float64(r) / Max)),
						float32(ErpsRGBToLinear(float64(g) / Max)),
						float32(ErpsRGBToLinear(float64(b) / Max)),
						float32(a) / Max,
					}
				}

				t.data[y*w+x] = c
			}
		}

		if t.onImage != nil {
			t.onImage(t.data, t.w, t.h)
		}
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

func textureToImage(u float64, imgSize int) (i, j int, t float32) {
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
		t = float32(s - f)
		i = int(f)
		j = i + 1

		if j >= imgSize {
			j--
		}
	}

	return
}

func (t *ImageTexture) TextureAt(x, y int) ColorRGBA {
	return t.data[y*t.w+x]
}

func (t *ImageTexture) TexelAtHit(ii *IntersectionInfo) ColorRGBA {
	// Get surface coordinates
	u := ii.U
	v := ii.V

	if t.onMapUv != nil {
		u, v = t.onMapUv(u, v, ii)
	}

	// Handle wrapping
	if u < 0 || u > 1 || v < 0 || v > 1 {
		switch t.wrap {
		case TwPeriodic:
			u -= math.Floor(u)
			v -= math.Floor(v)
		case TwFlip:
			fu := math.Floor(u)
			fv := math.Floor(v)
			u -= fu
			v -= fv
			if (int(math.Abs(fu+fv)) % 2) == 1 {
				u = 1 - u
				v = 1 - v
			}
		case TwClamp:
			u = Clamp(u)
			v = Clamp(v)
		case TwBlack:
			return ColorRGBA{0, 0, 0, 1}
		}
	}

	// Bilinear filtering
	x0, x1, tu := textureToImage(u, t.w)
	y0, y1, tv := textureToImage(v, t.h)

	c00 := t.TextureAt(x0, y0).Mul((1 - tu) * (1 - tv))
	c01 := t.TextureAt(x0, y1).Mul((1 - tu) * tv)
	c10 := t.TextureAt(x1, y0).Mul(tu * (1 - tv))
	c11 := t.TextureAt(x1, y1).Mul(tu * tv)

	return c00.Add(c01).Add(c10).Add(c11)
}

func (t *ImageTexture) ApplyAtHit(ii *IntersectionInfo) {
	ct := t.TexelAtHit(ii)

	if t.onApply != nil {
		// Use a custom handler if defined
		t.onApply(ct, ii)
	} else {
		// Image color is alpha pre-multiplied, so we only need to merge in the target color
		ii.Mat.DiffuseColor = ct.RGB().Add(ii.Mat.DiffuseColor.Mul(float64(1 - ct.a)))
	}
}

func (t *ImageTexture) NormalAtHit(ii *IntersectionInfo) Tuple {
	ct := t.TexelAtHit(ii)

	// Convert color into vector
	ct = ct.Mul(2).Add(ColorRGBA{-1, -1, -1, 0})

	return Vector(float64(ct.r), float64(ct.g), float64(ct.b))
}
