// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
)

type Raytracer struct {
	world *World
	xs    *Intersections
	ii    *IntersectionInfo
	rand  FloatGenerator
}

func NewRaytracer(world *World) *Raytracer {
	rt := Raytracer{
		world: world,
		xs:    NewIntersections(),
		ii:    &IntersectionInfo{},
		rand:  NewRandomGenerator(1),
	}

	return &rt
}

func (rt *Raytracer) HitForShadow(ray Ray) Intersection {
	xs := rt.xs
	xs.shadows = true // We look only for shadows now
	xs.Reset()
	for _, o := range rt.world.Objects {
		o.AddIntersections(ray, xs)
	}
	xs.shadows = false // Reset the shadow flag, as this list will be reused

	return xs.hit
}

func (rt *Raytracer) ShadeHit(ii *IntersectionInfo, depth int) (c Color) {
	m := ii.Mat

	c = m.DiffuseColor.Blend(rt.world.Ambient.Mul(ii.O.Material().Ambient))

	for _, light := range rt.world.Lights {
		c = c.Add(light.LightenHit(ii, rt))
	}

	if depth > 0 {
		if m.ReflectLevel > 0 {
			if m.RefractLevel > 0 {
				info := *ii // Need to copy the info locally because there are two recursive calls now and ii will be overwritten
				reflectance := SchlickReflectance(&info)
				reflected := rt.ReflectedColor(&info, depth)
				refracted := rt.RefractedColor(&info, depth)

				return c.
					Add(reflected.Mul(reflectance)).
					Add(refracted.Mul(1 - reflectance))
			} else {
				return c.Add(rt.ReflectedColor(ii, depth))
			}
		} else if m.RefractLevel > 0 {
			return c.Add(rt.RefractedColor(ii, depth))
		}
	}

	return c
}

func (rt *Raytracer) ColorForRay(r Ray, depth int) Color {
	xs := rt.xs // Reusing the intersection list greatly reduces memory usage and provides a very significant performance boost
	xs.Reset()

	// Intersect ray with all objects
	for _, o := range rt.world.Objects {
		o.AddIntersections(r, xs)
	}

	// Did we hit something?
	hit := xs.Hit()

	if hit.Valid() {
		ii := rt.ii
		ii.Update(hit, r, xs)

		return rt.ShadeHit(ii, depth)
	} else {
		return Black
	}
}

func (rt *Raytracer) ColorAt(ray Ray) Color {
	return rt.ColorForRay(ray, rt.world.Options.ReflectionDepth)
}

func (rt *Raytracer) ReflectedColor(ii *IntersectionInfo, depth int) (c Color) {
	if depth > 0 {
		reflectedRay := NewRay(ii.OverPoint, ii.Reflectv)

		c = ii.O.Material().Reflect.Blend(rt.ColorForRay(reflectedRay, depth-1))
	}

	return c
}

func (rt *Raytracer) RefractedColor(ii *IntersectionInfo, depth int) (c Color) {
	if depth > 0 {
		// Check for total internal reflection
		nRatio := ii.N1 / ii.N2
		cosThetai := ii.Eyev.DotProduct(ii.Normalv)               // θi is the angle of incidence
		sin2Thetat := nRatio * nRatio * (1 - cosThetai*cosThetai) // sin(θt)^2, where θt is the angle of refraction

		if sin2Thetat <= 1 {
			cosThetat := math.Sqrt(1 - sin2Thetat)
			direction := ii.Normalv.Mul(nRatio*cosThetai - cosThetat).Sub(ii.Eyev.Mul(nRatio))
			refractedRay := NewRay(ii.UnderPoint, direction)

			c = ii.O.Material().Refract.Blend(rt.ColorForRay(refractedRay, depth-1))
		}
		// ...else we got total internal reflection
	}

	return c
}
