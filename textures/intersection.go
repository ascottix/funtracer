// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package textures

import (
	"sort"

	. "ascottix/funtracer/maths"
)

type Patternable interface {
	WorldToObject(Tuple) Tuple
}

type Intersectable interface {
	NormalAtHit(ii *IntersectionInfo, xs *Intersections) Tuple
}

type Hittable interface {
	Patternable
	Intersectable
	Material() *Material
}

type Intersection struct {
	T float64
	O Hittable
	D int // 1-based index of data associated to this intersection
}

type IntersectionData struct {
	// Used by CSG
	Lhit bool
	// Used by Trimesh
	TU float64
	TV float64
}

type Intersections struct {
	L    []Intersection
	hit  Intersection
	data []IntersectionData
	// The following attributes are used to pass information from the renderer to the objects
	Shadows bool // If true, only shadowing objects should add to the intersections
}

func NewIntersection(t float64, o Hittable) Intersection {
	return Intersection{T: t, O: o, D: 0}
}

func (i Intersection) Valid() bool {
	return i.O != nil
}

func NewIntersections() *Intersections {
	return &Intersections{}
}

func (x *Intersections) Len() int {
	return len(x.L)
}

func (x *Intersections) Reset() {
	x.L = x.L[:0]
	x.data = x.data[:0]
	x.hit.O = nil
	x.hit.T = 0
	x.hit.D = 0
}

// Computes the hit of the current list
func (x *Intersections) UpdateHit() {
	x.hit.O = nil
	x.hit.T = 0
	x.hit.D = 0

	for _, i := range x.L {
		if i.T >= 0 && (x.hit.O == nil || x.hit.T > i.T) {
			x.hit = i
		}
	}
}

// Remove removes an intersection from the list, note this may invalidate the hit
// so it's often a good idea to call UpdateHit() when done with the list manipulations
func (x *Intersections) Remove(i int) {
	x.L = append(x.L[:i], x.L[i+1:]...) // Remove object
}

func (x *Intersections) SortRange(start int) *Intersections {
	less := func(a, b int) bool {
		return x.L[start+a].T < x.L[start+b].T
	}

	sort.Slice(x.L[start:], less)

	return x
}

func (x *Intersections) Sort() *Intersections {
	less := func(a, b int) bool {
		return x.L[a].T < x.L[b].T
	}

	sort.Slice(x.L, less)

	return x
}

func add(x *Intersections, i Intersection) {
	x.L = append(x.L, i)

	if i.T >= 0 && (x.hit.O == nil || x.hit.T > i.T) {
		x.hit = i
	}
}

func (x *Intersections) Add(o Hittable, t ...float64) *Intersections {
	for _, v := range t {
		i := Intersection{T: v, O: o, D: 0}
		add(x, i)
	}

	return x
}

func (x *Intersections) AddWithData(o Hittable, t float64) *IntersectionData {
	x.data = append(x.data, IntersectionData{})
	i := Intersection{T: t, O: o, D: len(x.data)}
	add(x, i)

	return x.Data(&i)
}

// Data returns a pointer to the data associated with an intersection, it may return nil
func (x *Intersections) Data(i *Intersection) (d *IntersectionData) {
	if i.D > 0 {
		d = &x.data[i.D-1]
	}

	return
}

// DataAt returns a pointer to the data associated with an intersection,
// if there is no such data then it is automaticallt allocated
// so this function never returns nil
func (x *Intersections) DataAt(i int) *IntersectionData {
	if x.L[i].D == 0 {
		x.data = append(x.data, IntersectionData{}) // Data is kept into an internal array, which can be reused (to help memory management costs)
		x.L[i].D = len(x.data)
	}

	return &x.data[x.L[i].D-1]
}

func (x *Intersections) At(i int) Intersection {
	return x.L[i]
}

func (x *Intersections) Hit() Intersection {
	return x.hit
}
