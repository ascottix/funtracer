// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"image/png"
	"os"
	"sync"
)

type World struct {
	Objects          []Groupable
	Lights           []Light
	Ambient          Color
	Options          *Options
	ErpCanvasToImage Interpolator
}

func NewWorld() *World {
	w := &World{}

	w.SetAmbient(White)
	w.SetOptions(NewOptions())
	w.ErpCanvasToImage = ErpGamma // Produce a gamma-corrected image by default

	return w
}

func (w *World) SetAmbient(color Color) {
	w.Ambient = color
}

func (w *World) SetOptions(options *Options) {
	w.Options = options
}

func (w *World) ObjectAsGroup(o interface{}) *Group {
	switch t := o.(type) {
	case *Group:
		return t
	}

	return nil
}

func (w *World) Find(name string) Groupable {
	for _, o := range w.Objects {
		if o.Name() == name {
			return o
		}
	}

	return nil
}

func (w *World) AddObjects(objects ...Groupable) {
	w.Objects = append(w.Objects, objects...)
}

func (w *World) AddLights(lights ...Light) {
	w.Lights = append(w.Lights, lights...)
}

func (w *World) Intersect(ray Ray) *Intersections {
	xs := NewIntersections()

	for _, o := range w.Objects {
		o.AddIntersections(ray, xs)
	}

	return xs
}

func (w *World) RenderToCanvas(c *Camera) Canvas {
	canvas := NewCanvas(c.HSize, c.VSize)

	for y := 0; y < c.VSize; y++ {
		for x := 0; x < c.HSize; x++ {
			ray := c.RayForPixelI(x, y)
			color := w.ColorAt(ray, w.Options.ReflectionDepth)
			canvas.FastSetPixelAt(x, y, color)
		}
	}

	return canvas
}

// Test methods
func (w *World) ShadeHit(ii *IntersectionInfo, depth int) (c Color) {
	return NewRaytracer(w).ShadeHit(ii, depth)
}

func (w *World) ColorAt(r Ray, depth int) Color {
	return NewRaytracer(w).ColorForRay(r, depth)
}

func (w *World) ReflectedColor(ii *IntersectionInfo, depth int) Color {
	return NewRaytracer(w).ReflectedColor(ii, depth)
}

func (w *World) RefractedColor(ii *IntersectionInfo, depth int) Color {
	return NewRaytracer(w).RefractedColor(ii, depth)
}

type XY struct {
	x, y float64
}

type XYC struct {
	XY
	c Color
}

func getPixelOffsetGenerator(random bool) (g FloatGenerator) {
	if random {
		g = NewRandomGenerator(1)
	} else {
		g = NewConstGenerator(0.5)
	}

	return g
}

func (w *World) GoDivisionRenderToCanvas(goers int, camera *Camera) Canvas {
	var wg sync.WaitGroup

	canvas := NewCanvas(camera.HSize, camera.VSize)

	renderer := func(m, r int) {
		rt := NewRaytracer(w)

		offset := getPixelOffsetGenerator(w.Options.SamplesPerPixel > 1)

		for y := 0; y < camera.VSize; y++ {
			for x := r; x < camera.HSize; x += m {
				for s := 0; s < w.Options.SamplesPerPixel; s++ {
					px := float64(x) + offset()
					py := float64(y) + offset()
					ray := camera.RayForPixelF(px, py)
					col := rt.ColorAt(ray)
					canvas.AddPixelAt(px, py, col)
				}
			}
		}

		wg.Done()
	}

	for i := 0; i < goers; i++ {
		wg.Add(1)
		go renderer(goers, i)
	}

	wg.Wait()

	canvas.Mul(1.0 / float64(w.Options.SamplesPerPixel))

	return canvas
}

func (w *World) GoPipelineRenderToCanvas(goers int, camera *Camera) Canvas {
	rasterizer := func(out chan XY) {
		offset := getPixelOffsetGenerator(w.Options.SamplesPerPixel > 1)

		for y := 0; y < camera.VSize; y++ {
			for x := 0; x < camera.HSize; x++ {
				for s := 0; s < w.Options.SamplesPerPixel; s++ {
					out <- XY{float64(x) + offset(), float64(y) + offset()}
				}
			}
		}

		close(out)
	}

	renderer := func(in chan XY, out chan XYC) {
		var wg sync.WaitGroup

		for i := 0; i < goers; i++ {
			wg.Add(1)
			go func() {
				rt := NewRaytracer(w)

				for xy := range in {
					ray := camera.RayForPixelF(xy.x, xy.y)
					col := rt.ColorAt(ray)
					out <- XYC{xy, col}
				}
				wg.Done()
			}()
		}

		wg.Wait()
		close(out)
	}

	imager := func(in chan XYC, out chan Canvas) {
		cvs := NewCanvas(camera.HSize, camera.VSize)

		for xyc := range in {
			cvs.AddPixelAt(xyc.x, xyc.y, xyc.c)
		}

		cvs.Mul(1.0 / float64(w.Options.SamplesPerPixel))

		out <- cvs
	}

	cxy := make(chan XY, 8)
	cpix := make(chan XYC, 8)
	cimg := make(chan Canvas)

	go rasterizer(cxy)

	go renderer(cxy, cpix)

	go imager(cpix, cimg)

	canvas := <-cimg

	return canvas
}

func (w *World) RenderToImage(c *Camera) image.Image {
	canvas := w.GoDivisionRenderToCanvas(w.Options.NumThreads, c)
	// Alternative renderers
	// canvas := w.GoPipelineRenderToCanvas(w.Options.NumThreads, c)
	// canvas := w.RenderToCanvas(c)

	return canvas.ToImage(w.ErpCanvasToImage)
}

func (w *World) RenderToPNG(c *Camera, filename string) error {
	img := w.RenderToImage(c)

	f, err := os.Create(filename)

	if err == nil {
		defer f.Close()

		err = png.Encode(f, img)
	}

	return err
}
