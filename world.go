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

func (w *World) getPixelSampler(jitterId int) Sampler2d {
	if w.Options.Supersampling == 1 {
		return NewStratified2d(1, 1)
	}

	return NewJitteredStratified2d(w.Options.Supersampling, w.Options.Supersampling, NewRandomGenerator(13+int64(jitterId)*7))
}

func (w *World) GoDivisionRenderToCanvas(goers int, camera *Camera) Canvas {
	var wg sync.WaitGroup

	canvas := NewCanvas(camera.HSize, camera.VSize)

	samplesPerPixel := w.Options.Supersampling * w.Options.Supersampling

	renderer := func(m, r int) {
		defer wg.Done()

		rt := NewRaytracer(w)

		sampler := w.getPixelSampler(r)

		for y := 0; y < camera.VSize; y++ {
			for x := r; x < camera.HSize; x += m {
				sampler.Reset()
				for s := 0; s < samplesPerPixel; s++ {
					px, py := sampler.Next()
					px += float64(x)
					py += float64(y)
					ray := camera.RayForPixelF(px, py)
					col := rt.ColorAt(ray)
					canvas.AddPixelAt(px, py, col)
				}
			}
		}
	}

	wg.Add(goers)
	for i := 0; i < goers; i++ {
		go renderer(goers, i)
	}

	wg.Wait()

	canvas.Mul(1.0 / float64(samplesPerPixel))

	return canvas
}

func (w *World) GoPipelineRenderToCanvas(goers int, camera *Camera) Canvas {
	samplesPerPixel := w.Options.Supersampling * w.Options.Supersampling

	sampler := func(out chan XY) {
		js2d := w.getPixelSampler(0)

		for y := 0; y < camera.VSize; y++ {
			for x := 0; x < camera.HSize; x++ {
				js2d.Reset()
				for s := 0; s < samplesPerPixel; s++ {
					px, py := js2d.Next()
					out <- XY{float64(x) + px, float64(y) + py}
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
				defer wg.Done()

				rt := NewRaytracer(w)

				for xy := range in {
					ray := camera.RayForPixelF(xy.x, xy.y)
					col := rt.ColorAt(ray)
					out <- XYC{xy, col}
				}
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

		cvs.Mul(1.0 / float64(samplesPerPixel))

		out <- cvs
	}

	cxy := make(chan XY, 8)   // sampler -> (x,y) coordinates of point where pixel is sampled
	cpix := make(chan XYC, 8) // (x,y) -> renderers -> (x,y,color of sampled point)
	cimg := make(chan Canvas) // (x,y,color) -> image

	go sampler(cxy)

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
