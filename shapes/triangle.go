// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package shapes

import (
	. "ascottix/funtracer/maths"
	. "ascottix/funtracer/textures"
)

// Note: creating a shape from this Shapable is not very efficient, a Trimesh is almost always a better choice

// Triangle is defined by three points (not on the same line) in space
type Triangle struct {
	P1 Tuple // Vertices
	P2 Tuple
	P3 Tuple
	E1 Tuple // Edges
	E2 Tuple
	N  Tuple // Normal
}

func NewTriangle(p1, p2, p3 Tuple) *Triangle {
	e1 := p2.Sub(p1)
	e2 := p3.Sub(p1)
	n := e2.CrossProduct(e1).Normalize()

	return &Triangle{p1, p2, p3, e1, e2, n}
}

func (t *Triangle) Bounds() Box {
	return Box{
		Point(Min3(t.P1.X, t.P2.X, t.P3.X), Min3(t.P1.Y, t.P2.Y, t.P3.Y), Min3(t.P1.Z, t.P2.Z, t.P3.Z)),
		Point(Max3(t.P1.X, t.P2.X, t.P3.X), Max3(t.P1.Y, t.P2.Y, t.P3.Y), Max3(t.P1.Z, t.P2.Z, t.P3.Z)),
	}
}

// LocalIntersect uses the Möller–Trumbore algorithm to find the intersection with a ray
func (t *Triangle) LocalIntersect(ray Ray) []float64 {
	dirCrossE2 := ray.Direction.CrossProduct(t.E2)
	det := t.E1.DotProduct(dirCrossE2)

	// Check if ray is parallel to triangle plane
	if det > -Epsilon && det < +Epsilon {
		return nil
	}

	f := 1.0 / det
	p1ToOrigin := ray.Origin.Sub(t.P1)
	u := f * p1ToOrigin.DotProduct(dirCrossE2)

	// Ray misses by the p1-p3 edge
	if u < 0 || u > 1 {
		return nil
	}

	originCrossE1 := p1ToOrigin.CrossProduct(t.E1)
	v := f * ray.Direction.DotProduct(originCrossE1)

	// Ray misses by the p1-p2 or p2-p3 edge
	if v < 0 || (u+v) > 1 {
		return nil
	}

	// Ray hits
	h := f * t.E2.DotProduct(originCrossE1)

	return []float64{h}
}

func (t *Triangle) LocalNormalAt(point Tuple) Tuple {
	return t.N
}
