// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package shapes

import (
	. "ascottix/funtracer/textures"
	. "ascottix/funtracer/traits"
)

type CsgOp int

const (
	CsgUnion CsgOp = iota
	CsgDifference
	CsgIntersection
)

type Csg struct {
	Namer
	Grouper
	L      Groupable
	R      Groupable
	Op     CsgOp
	OpTest func(bool, bool, bool) bool
}

// CsgUnion builds a CSG shape that is the union of a set of objects
func NewCsgUnion(o ...Groupable) (csg *Csg) {
	// Try to get a balanced tree
	for len(o) >= 2 {
		// Even indices become CSG
		for i := 1; i < len(o); i += 2 {
			o[i-1] = NewCsg(CsgUnion, o[i-1], o[i])
		}

		// Odd indices are removed
		n := (len(o) + 1) / 2
		for i := 1; i < n; i++ {
			o[i] = o[i*2]
		}

		o = o[:n]
	}

	if len(o) == 1 {
		switch t := o[0].(type) {
		case *Csg:
			csg = t
		}
	}

	return
}

func NewCsg(op CsgOp, a, b Groupable) *Csg {
	g := &Csg{}

	g.SetNameForKind("csg")
	g.SetTransform()
	g.L = a
	g.R = b
	g.Op = op

	g.L.SetParent(g)
	g.R.SetParent(g)

	switch op {
	case CsgUnion:
		g.OpTest = opTestUnion
	case CsgDifference:
		g.OpTest = opTestDifference
	case CsgIntersection:
		g.OpTest = opTestIntersection
	}

	return g
}

func (g *Csg) Clone() Groupable {
	o := NewCsg(g.Op, g.L.Clone(), g.R.Clone())

	o.SetName("csgfrom_" + g.Name())
	o.SetTransform(g.Transform())

	return o
}

// Hit belongs to union if:
// - it's in the left shape, and the ray is not in the right shape
// - it's in the right shape, and the ray is not in the left shape
func opTestUnion(lhit, inL, inR bool) bool {
	return (lhit && !inR) || (!lhit && !inL)
}

// Hit belongs to difference if:
// - it's in the left shape, and the ray is not in the right shape
// - it's in the right shape, and the ray is in the left shape
func opTestDifference(lhit, inL, inR bool) bool {
	return (lhit && !inR) || (!lhit && inL)
}

// Hit belongs to intersection if:
// - it's in the left shape, and the ray is in the right shape
// - it's in the right shape, and the ray is in the left shape
func opTestIntersection(lhit, inL, inR bool) bool {
	return (lhit && inR) || (!lhit && inL)
}

func (g *Csg) AddIntersections(ray Ray, xs *Intersections) {
	ray = ray.Transform(g.Tinverse)

	sIdx := xs.Len()

	// Add intersections for the left shape
	g.L.AddIntersections(ray, xs)

	for i := sIdx; i < xs.Len(); i++ {
		xs.DataAt(i).Lhit = true
	}

	// Add intersections for the right shape
	rIdx := xs.Len()
	g.R.AddIntersections(ray, xs)

	for i := rIdx; i < xs.Len(); i++ {
		xs.DataAt(i).Lhit = false
	}

	// Sort our intersections
	xs.SortRange(sIdx)

	// Now scan the intersection list keeping track of objects as the ray enters and leaves them
	inL := false
	inR := false

	for i := sIdx; i < xs.Len(); {
		lhit := xs.DataAt(i).Lhit

		ok := g.OpTest(lhit, inL, inR)

		// Update object trackers
		if lhit {
			inL = !inL
		} else {
			inR = !inR
		}

		if !ok {
			// Remove intersection
			xs.Remove(i)
		} else {
			// Keep intersection and move to the next!
			i++
		}
	}

	// Recompute the hit as we may have disrupted the hit list
	xs.UpdateHit()
}

func (g *Csg) SetMaterial(m *Material) {
	m = m.ProxifyPatterns(g)

	g.L.SetMaterial(m)
	g.R.SetMaterial(m)
}

func (g *Csg) Bounds() Box {
	b := g.L.Bounds().Transform(g.L.Transform())

	return b.Union(g.R.Bounds().Transform(g.R.Transform()))
}
