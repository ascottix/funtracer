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
	linear  bool // TODO: remove this!
}

type NormalMap struct {
}

func NewTextureOnMapUvHandler(addU, mulU, addV, mulV float64) TextureOnMapUv {
	return func(u, v float64, ii *IntersectionInfo) (float64, float64) {
		return u*mulU + addU, v*mulV + addV
	}
}

const MaxImageRgbaValue = 65535

var colorLinearLookup = make([]float32, MaxImageRgbaValue+1)
var colorGammaLookup = make([]float32, MaxImageRgbaValue+1)

func TextureMirrorV(u, v float64, ii *IntersectionInfo) (float64, float64) {
	return u, 1 - v
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

		colorLookup := colorGammaLookup

		if t.linear {
			colorLookup = colorLinearLookup
		}

		// Convert the texture to ColorRGBA format
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, b, a := image.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()

				o := x + y*w
				t.data[o].r = colorLookup[r]
				t.data[o].g = colorLookup[g]
				t.data[o].b = colorLookup[b]
				t.data[o].a = float32(a) / MaxImageRgbaValue
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
	} else {
		panic("Cannot load texture: " + filename)
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
	if t.w == 0 {
		return ColorRGBA{0, 0, 0, 1}
	}
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

// Converting from the image color space can be _very_ slow, a lookup table improves performance a lot
func init() {
	for v := 0; v <= MaxImageRgbaValue; v++ {
		f := float64(v) / MaxImageRgbaValue               // Value (a color component) is now a 64-bit float from 0 to 1
		colorLinearLookup[v] = float32(f)                 // No conversion
		colorGammaLookup[v] = float32(ErpsRGBToLinear(f)) // Convert from standard gamma
	}
}
