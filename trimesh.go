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
	VT       []Tuple // Vertex texture coordinates
	T        []MeshTriangle
	material *Material // TODO! Handling of material needs to be refactored
}

type MeshTriangle struct {
	mesh *Trimesh
	V    [3]int // Vertex indices
	VN   [3]int // Vertex normals
	VT   [3]int // Texture vertices
	E1   Tuple  // Edges
	E2   Tuple
	N    Tuple // Normal
	T    Tuple // Tangent vector
	B    Tuple // Bitangent vector
	Mat  *Material
}

func NewTrimesh(info *ObjInfo, group int) *Trimesh {
	mesh := Trimesh{}

	mesh.material = NewMaterial()

	mesh.SetNameForKind("mesh")
	mesh.SetTransform()

	mesh.V = info.V
	mesh.VN = info.VN
	mesh.VT = info.VT

	for _, f := range info.F {
		if group == -1 || f.G == group {
			var p [3]Tuple

			mt := MeshTriangle{}
			mt.mesh = &mesh

			mt.V = f.V
			mt.VN = f.VN
			mt.VT = f.VT

			for i := 0; i < 3; i++ {
				p[i] = mesh.V[mt.V[i]] // Vertex
			}

			// Compute normal
			mt.E1 = p[1].Sub(p[0])
			mt.E2 = p[2].Sub(p[0])
			mt.N = mt.E2.CrossProduct(mt.E1).Normalize()

			// Compute tangent and bitangent vectors
			if len(mesh.VT) > 0 {
				vt0 := mesh.VT[mt.VT[0]]
				vt1 := mesh.VT[mt.VT[1]]
				vt2 := mesh.VT[mt.VT[2]]

				u1 := vt1.X - vt0.X
				v1 := vt1.Y - vt0.Y
				u2 := vt2.X - vt0.X
				v2 := vt2.Y - vt0.Y

				d := 1 / (u1*v2 - v1*u2)

				q1 := mt.E1
				q2 := mt.E2

				T := Vector(q1.X*v2-q2.X*v1, q1.Y*v2-q2.Y*v1, q1.Z*v2-q2.Z*v1).Mul(d)
				B := Vector(q2.X*u1-q1.X*u2, q2.Y*u1-q1.Y*u2, q2.Z*u1-q1.Z*u2).Mul(d)

				mt.T = T.Normalize()
				mt.B = B.Normalize()
			} else {
				mt.T = Vector(0xBAD, 0, 0)
				mt.B = Vector(0xBAD, 0, 0)
			}

			mt.Mat = f.M

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

	for i := range s.T {
		s.T[i].SetMaterial(nil)
	}
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
	if t.Mat != nil {
		return t.Mat
	}

	return t.mesh.Material()
}

func (t *MeshTriangle) SetMaterial(m *Material) {
	t.Mat = m
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
	id := xs.Data(&ii.Intersection)

	// Texture
	if len(t.mesh.VT) > 0 {
		VT1 := t.mesh.VT[t.VT[0]]
		VT2 := t.mesh.VT[t.VT[1]]
		VT3 := t.mesh.VT[t.VT[2]]

		vt := VT2.Mul(id.tU).Add(VT3.Mul(id.tV)).Add(VT1.Mul(1 - id.tU - id.tV))
		ii.U = vt.X
		ii.V = vt.Y
	}

	// Normal
	N := t.N

	// This is a trick to mark the triangle edges (wireframe),
	// but the line width depends on the triangle area and therefore is not constant,
	// which should be fixed
	/*
		const EdgeWidth = 0.02
		if id.tU < EdgeWidth || id.tV < EdgeWidth || (1-id.tU-id.tV) < EdgeWidth {
			return Vector(0,0,0)
		}
	*/

	if len(t.mesh.VN) > 0 {
		N1 := t.mesh.VN[t.VN[0]]
		N2 := t.mesh.VN[t.VN[1]]
		N3 := t.mesh.VN[t.VN[2]]

		N = N2.Mul(id.tU).Add(N3.Mul(id.tV)).Add(N1.Mul(1 - id.tU - id.tV))
	}

	// Apply normal map if present
	if nmap := ii.GetNormalMap(); nmap != nil {
		n := nmap.NormalAtHit(ii)

		// TODO: I think T and B should be interpolated like the normal!
		T := t.T
		B := t.B

		ii.SurfNormalv = (T.Mul(n.X).Add(B.Mul(n.Y)).Add(N.Mul(n.Z))).Normalize()
		ii.HasSurfNormalv = true
	}

	return t.mesh.NormalToWorld(N)
}

func (t *MeshTriangle) WorldToObject(point Tuple) Tuple {
	return t.mesh.WorldToObject(point)
}
