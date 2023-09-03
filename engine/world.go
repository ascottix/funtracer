// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package engine

import (
	"image"
	"image/png"
	"os"
	"sync"

	. "ascottix/funtracer/maths"
	. "ascottix/funtracer/options"
	. "ascottix/funtracer/shapes"
	. "ascottix/funtracer/textures"
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
	w.ErpCanvasToImage = ErpLinearToGamma // Produce a gamma-corrected image by default

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
			ray := c.RayForPixel(float64(x)+0.5, float64(y)+0.5)
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

func (w *World) getPixelSampler(rand FloatGenerator) (s Sampler2d) {
	if w.Options.Supersampling == 1 {
		s = NewStratified2d(1, 1)
	} else {
		s = NewJitteredStratified2d(w.Options.Supersampling, w.Options.Supersampling, rand)
	}

	return s
}

func (w *World) GoDivisionRenderToCanvas(goers int, camera *Camera) Canvas {
	var wg sync.WaitGroup

	canvas := NewCanvas(camera.HSize, camera.VSize)

	samplesPerPixel := w.Options.Supersampling * w.Options.Supersampling

	renderer := func(m, r int) {
		defer wg.Done()

		rt := NewRaytracer(w)

		sampler := w.getPixelSampler(rt.rand) // NewRandomGenerator(13+int64(s)*7)

		for y := 0; y < camera.VSize; y++ {
			for x := r; x < camera.HSize; x += m {
				// Reset sampler to keep all values into the proper range
				sampler.Reset()

				for s := 0; s < samplesPerPixel; s++ {
					// Get the pixel coordinates
					px, py := sampler.Next()
					px += float64(x)
					py += float64(y)

					// Get ray from viewpoint to target pixel
					var ray Ray
					if w.Options.LensRadius > 0 {
						ray = camera.RayForPixelDepthOfField(px, py, w.Options.LensRadius, w.Options.FocalDistance, rt.rand)
					} else {
						ray = camera.RayForPixel(px, py)
					}

					// Render and store color
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

func (w *World) RenderToImage(c *Camera) image.Image {
	canvas := w.GoDivisionRenderToCanvas(w.Options.NumThreads, c)
	// Alternative renderers
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
