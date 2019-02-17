// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

// Shapable provides the basic geometry of a primitive shape
type Shapable interface {
	Bounds() Box
	LocalNormalAt(Tuple) Tuple
	LocalIntersect(Ray) []float64
	NormalAtHit(Tuple, *IntersectionInfo) Tuple
}

// Shape is an higher level object that represents a geometric primitive,
// it uses the Shapable interface for anything related to the shape geometry
type Shape struct {
	Namer
	Grouper
	material *Material
	shapable Shapable
	shadow   bool
	Locked   bool
}

func NewShape(kind string, shapable Shapable) *Shape {
	s := &Shape{}

	s.material = NewMaterial()
	s.shapable = shapable
	s.shadow = true
	s.SetNameForKind(kind)
	s.SetTransform()

	return s
}

func (s *Shape) Clone() Groupable {
	o := NewShape("", s.shapable)

	o.material = s.material
	o.shadow = s.shadow

	o.SetName("shape_from_" + s.Name())
	o.SetTransform(s.Transform())

	return o
}

func (s *Shape) Material() *Material {
	return s.material
}

func (s *Shape) SetMaterial(m *Material) {
	if !s.Locked {
		s.material = m
	}
}

func (s *Shape) SetLocked(f bool) {
	s.Locked = f // Used for testing and for special objects like bounding boxes: disallows setting some shape properties
}

func (s *Shape) SetShadow(f bool) {
	s.shadow = f
}

func (s *Shape) Parent() Container {
	return s.parent
}

func (s *Shape) Intersect(ray Ray) *Intersections { // Self-contained but slow: currently used only for testing
	xs := NewIntersections()

	s.AddIntersections(ray, xs)

	return xs
}

func (s *Shape) AddIntersections(ray Ray, xs *Intersections) {
	// If looking for shadows and this shape does not cast one, exit now
	if xs.shadows && !s.shadow {
		return
	}

	ray = ray.Transform(s.Tinverse)

	localxs := s.shapable.LocalIntersect(ray)

	xs.Add(s, localxs...)
}

// NormalAt returns the normal at the specified point on this shape,
// it has been replaced by NormalAtHit() and now used only for tests.
func (s *Shape) NormalAt(point Tuple) Tuple {
	ii := IntersectionInfo{Point: point}

	return s.NormalAtHit(&ii, nil)
}

// NormalAtHit returns the normal at an intersection point on this shape,
// it may also fill other information in the IntersectionInfo (e.g. the u,v surface coords)
func (s *Shape) NormalAtHit(ii *IntersectionInfo, xs *Intersections) Tuple {
	// Convert the point into object space, so it can be handled by the simple primitive
	point := s.WorldToObject(ii.Point)

	// Get the normal in object space, may fill other info as well
	normal := s.shapable.NormalAtHit(point, ii)

	if ii.HasSurfNormalv {
		ii.SurfNormalv = s.NormalToWorld(ii.SurfNormalv)
	}

	// Return the normal in world space
	return s.NormalToWorld(normal)
}

// Make a shape a Shapable object itself
func (s *Shape) Bounds() Box {
	return s.shapable.Bounds()
}

func (s *Shape) LocalNormalAt(point Tuple) Tuple {
	return s.NormalAt(point)
}

func (s *Shape) LocalIntersect(ray Ray) []float64 {
	ray = ray.Transform(s.Tinverse)

	return s.shapable.LocalIntersect(ray)
}
