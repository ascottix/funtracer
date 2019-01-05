// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"image"
	"io"
)

type Canvas struct {
	Width  int
	Height int
	Pix    []Color
}

func NewCanvas(width, height int) Canvas {
	canvas := Canvas{
		width,
		height,
		make([]Color, width*height, width*height),
	}

	return canvas
}

func xyFloatToInt(fx, fy float64, w, h int) (x, y int, ok bool) {
	x = int(fx)
	y = int(fy)

	ok = (x >= 0) && (x < w) && (y >= 0) && (y < h)

	return
}

func (canvas Canvas) PixelAt(fx, fy float64) (c Color) {
	if x, y, ok := xyFloatToInt(fx, fy, canvas.Width, canvas.Height); ok {
		c = canvas.Pix[x+y*canvas.Width]
	}

	return
}

func (canvas Canvas) SetPixelAt(fx, fy float64, c Color) {
	if x, y, ok := xyFloatToInt(fx, fy, canvas.Width, canvas.Height); ok {
		canvas.Pix[x+y*canvas.Width] = c
	}
}

func (canvas Canvas) AddPixelAt(fx, fy float64, c Color) {
	// Warning... no boundary checking here!
	x := int(fx)
	y := int(fy)
	o := x + y*canvas.Width

	canvas.Pix[o] = canvas.Pix[o].Add(c)
}

func (canvas Canvas) FastPixelAt(x, y int) Color {
	return canvas.Pix[x+y*canvas.Width]
}

func (canvas Canvas) FastSetPixelAt(x, y int, c Color) {
	canvas.Pix[x+y*canvas.Width] = c
}

func (canvas Canvas) Fill(c Color) {
	for i := range canvas.Pix {
		canvas.Pix[i] = c
	}
}

func (canvas Canvas) Mul(f float64) {
	for i := range canvas.Pix {
		canvas.Pix[i] = canvas.Pix[i].Mul(f)
	}
}

func (canvas Canvas) ToImage(erp Interpolator) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, canvas.Width, canvas.Height))

	for y := 0; y < canvas.Height; y++ {
		for x := 0; x < canvas.Width; x++ {
			col := canvas.FastPixelAt(x, y)
			img.Set(x, y, col.Erp(erp))
		}
	}

	return img
}

// Exports the canvas in PPM format
func (canvas Canvas) WriteAsPPM(w io.Writer) {
	fmt.Fprintf(w, "P3\n") // Magic
	fmt.Fprintf(w, "%d %d\n", canvas.Width, canvas.Height)
	fmt.Fprintf(w, "255\n") // Maximum value of a color component

	for y := 0; y < canvas.Height; y++ {
		for x := 0; x < canvas.Width; x++ {
			r, g, b, _ := canvas.Pix[x+y*canvas.Width].RGBA()
			fmt.Fprintf(w, "%d %d %d ", r, g, b) // In theory, each line should not exceed 70 characters
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w)
}
