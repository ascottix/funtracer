// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

type Trimesh struct {
	Namer
	Grouper
	V        []Tuple // Vertices
	VN       []Tuple // Vertex normals
	T        []MeshTriangle
	material *Material // TODO! Handling of material needs to be refactored
}

type MeshTriangle struct {
	mesh *Trimesh
	V    [3]int // Vertex indices
	VN   [3]int // Vertex normals
	E1   Tuple  // Edges
	E2   Tuple
	N    Tuple // Normal
}

func NewTrimesh(info *ObjInfo, group int) *Trimesh {
	mesh := Trimesh{}

	mesh.material = NewMaterial()

	mesh.SetNameForKind("mesh")
	mesh.SetTransform()

	mesh.V = info.V
	mesh.VN = info.VN

	for _, f := range info.F {
		if group == -1 || f.G == group {
			var p [3]Tuple

			mt := MeshTriangle{}
			mt.mesh = &mesh

			mt.V = f.V
			mt.VN = f.VN

			for i := 0; i < 3; i++ {
				p[i] = mesh.V[mt.V[i]] // Vertex
			}

			mt.E1 = p[1].Sub(p[0])
			mt.E2 = p[2].Sub(p[0])
			mt.N = mt.E2.CrossProduct(mt.E1).Normalize()

			mesh.T = append(mesh.T, mt)
		}
	}

	return &mesh
}

func (s *Trimesh) Material() *Material {
	return s.material
}

func (s *Trimesh) SetMaterial(m *Material) {
	s.material = m
}

func (s *Trimesh) AddToGroup(group *Group) {
	for i := range s.T {
		group.Add(&(s.T[i]))
	}
}

func (t *MeshTriangle) Transform() Matrix {
	return t.mesh.Transform()
}

func (t *MeshTriangle) InverseTransform() Matrix {
	return t.mesh.InverseTransform()
}

// SetTransform should never be called on a single mesh triangle
func (t *MeshTriangle) SetTransform(transforms ...Matrix) {
	t.mesh.SetTransform(transforms...)
}

func (t *MeshTriangle) Clone() Groupable {
	o := *t

	return &o
}

func (t *MeshTriangle) Bounds() Box {
	p1 := t.mesh.V[t.V[0]]
	p2 := t.mesh.V[t.V[1]]
	p3 := t.mesh.V[t.V[2]]
	return Box{
		Point(Min3(p1.X, p2.X, p3.X), Min3(p1.Y, p2.Y, p3.Y), Min3(p1.Z, p2.Z, p3.Z)),
		Point(Max3(p1.X, p2.X, p3.X), Max3(p1.Y, p2.Y, p3.Y), Max3(p1.Z, p2.Z, p3.Z)),
	}
}

func (t *MeshTriangle) Material() *Material {
	return t.mesh.Material()
}

// SetMaterial should never be called on a single mesh triangle
func (t *MeshTriangle) SetMaterial(m *Material) {
	t.mesh.SetMaterial(m)
}

func (t *MeshTriangle) Name() string {
	// Let's make a name for debugging
	return fmt.Sprintf("%s_t_%d_%d_%d", t.mesh.Name(), t.V[0], t.V[1], t.V[2])
}

// SetName should never be called on a single mesh triangle
func (t *MeshTriangle) SetName(name string) {
}

func (t *MeshTriangle) Parent() Container {
	return t.mesh.Parent()
}

func (t *MeshTriangle) SetParent(p Container) {
	t.mesh.SetParent(p)
}

func (t *MeshTriangle) AddIntersections(ray Ray, xs *Intersections) {
	dirCrossE2 := ray.Direction.CrossProduct(t.E2)
	det := t.E1.DotProduct(dirCrossE2)

	const E = 1e-9 // We need higher precision than usual here for models with many small triangles

	// Check if ray is parallel to triangle plane
	if det > -E && det < +E {
		return
	}

	f := 1.0 / det
	p1 := t.mesh.V[t.V[0]]
	p1ToOrigin := ray.Origin.Sub(p1)
	u := f * p1ToOrigin.DotProduct(dirCrossE2)

	// Ray misses by the p1-p3 edge
	if u < 0 || u > 1 {
		return
	}

	originCrossE1 := p1ToOrigin.CrossProduct(t.E1)
	v := f * ray.Direction.DotProduct(originCrossE1)

	// Ray misses by the p1-p2 or p2-p3 edge
	if v < 0 || (u+v) > 1 {
		return
	}

	// Ray hits
	h := f * t.E2.DotProduct(originCrossE1)

	id := xs.AddWithData(t, h)

	id.tU = u
	id.tV = v
}

func (t *MeshTriangle) NormalAtHit(ii *IntersectionInfo, xs *Intersections) Tuple {
	if len(t.mesh.VN) == 0 {
		return t.mesh.NormalToWorld(t.N) // Flat normal
	}
	// ...else interpolate

	id := xs.Data(&ii.Intersection)

	N1 := t.mesh.VN[t.VN[0]]
	N2 := t.mesh.VN[t.VN[1]]
	N3 := t.mesh.VN[t.VN[2]]

	n := N2.Mul(id.tU).Add(N3.Mul(id.tV)).Add(N1.Mul(1 - id.tU - id.tV))

	return t.mesh.NormalToWorld(n)
}

func (t *MeshTriangle) WorldToObject(point Tuple) Tuple {
	return t.mesh.WorldToObject(point)
}
