// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package shapes

import (
	. "ascottix/funtracer/maths"
	. "ascottix/funtracer/textures"
	. "ascottix/funtracer/traits"
)

type Container interface {
	Patternable
	NormalToWorld(Tuple) Tuple
}

type Groupable interface {
	Namable
	Transformable
	AddIntersections(Ray, *Intersections)
	SetMaterial(*Material)
	SetParent(Container)
	Bounds() Box
	Clone() Groupable
}

type Grouper struct {
	Transformer
	parent Container
}

type Group struct {
	Namer
	Grouper
	members  []Groupable
	bbox     Box // Bounding box
	bvhNodes []BvhLinearNode
}

func (g *Grouper) Parent() Container {
	return g.parent
}

func (g *Grouper) SetParent(p Container) {
	g.parent = p
}

func (g *Grouper) WorldToObject(point Tuple) Tuple {
	if g.parent != nil {
		point = g.parent.WorldToObject(point)
	}

	return g.Tinverse.MulT(point)
}

func (g *Grouper) NormalToWorld(normal Tuple) Tuple {
	normal = g.TinverseT.MulT(normal)
	normal.W = 0
	normal = normal.Normalize()

	if g.parent != nil {
		normal = g.parent.NormalToWorld(normal)
	}

	return normal
}

func NewGroup() *Group {
	g := &Group{}

	g.SetNameForKind("group")
	g.SetTransform()

	g.bbox = Box{PointAtInfinity(+1), PointAtInfinity(-1)}

	return g
}

func (g *Group) Clone() Groupable {
	o := NewGroup()

	o.SetName("groupfrom_" + g.Name())
	o.SetTransform(g.Transform())

	for _, s := range g.members {
		o.Add(s.Clone())
	}

	return o
}

func (g *Group) Add(elements ...Groupable) {
	for _, s := range elements {
		g.members = append(g.members, s)
		s.SetParent(g)
		g.bbox = g.bbox.Union(s.Bounds().Transform(s.Transform()))
	}
}

func (g *Group) SetMaterial(m *Material) {
	m = m.ProxifyPatterns(g)

	for _, s := range g.members {
		s.SetMaterial(m)
	}
}

func (g *Group) LocalIntersect(ray Ray) *Intersections { // Used only for testing
	xs := NewIntersections()

	for _, s := range g.members {
		s.AddIntersections(ray, xs)
	}

	xs.Sort()

	return xs
}

func (g *Group) AddIntersections(ray Ray, xs *Intersections) {
	if len(g.bvhNodes) > 0 {
		// Intersect using the BVH
		g.AddIntersectionsBvh(ray, xs)
	} else {
		// Standard intersection
		ray = ray.Transform(g.Tinverse)

		// Check hit against bounding box
		if !g.bbox.Intersects(ray) {
			return
		}

		// Check against all children
		for _, s := range g.members {
			s.AddIntersections(ray, xs)
		}
	}
}

func (g *Group) Bounds() Box {
	return g.bbox
}

func (g *Group) BoundingBox() *Shape {
	bb := g.Bounds().ToCube()

	bb.Material().SetDiffuseColor(RGB(0.5, 0.5, 0)).SetRefractive(1, 1)
	bb.SetShadow(false)
	bb.SetTransform(g.Transform(), bb.Transform())
	bb.SetLocked(true) // Protect the shape properties from overwrites

	return bb
}

func (g *Group) Members(index int) Groupable {
	return g.members[index]
}
