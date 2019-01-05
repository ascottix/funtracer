// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

type TriangleSmooth struct {
	V  [3]int // Vertex indices
	VN [3]int // Vertex normals
	P1 Tuple  // First vertex
	E1 Tuple  // Edges
	E2 Tuple
	N  Tuple // Normal
}

type TrimeshVertexInfo struct {
	V  int // Index in vertex array
	VN int // Index in vertex normal array
}

type TrimeshTriangleInfo struct {
	F [3]TrimeshVertexInfo // Info for each of the three vertices
}

type Trimesh struct {
	Namer
	Grouper
	V        []Tuple // Vertices
	VN       []Tuple // Vertex normals
	F        []TriangleSmooth
	material *Material // TODO! Handling of material needs to be refactored
	bbox     Box
	culling  bool
}

func NewTrimesh() *Trimesh {
	mesh := Trimesh{}

	mesh.material = NewMaterial()

	mesh.SetNameForKind("mesh")
	mesh.SetTransform()

	return &mesh
}

func (s *Trimesh) Clone() Groupable {
	o := NewTrimesh()

	o.V = s.V
	o.VN = s.VN
	o.F = s.F
	o.material = s.material
	o.bbox = s.bbox
	o.culling = s.culling
	o.SetName("meshfrom_" + s.Name())

	return o
}

func (mesh *Trimesh) SetMesh(V []Tuple, VN []Tuple, T []TrimeshTriangleInfo) {
	mesh.V = V
	mesh.VN = VN
	mesh.F = mesh.F[:0]
	mesh.bbox = Box{PointAtInfinity(+1), PointAtInfinity(-1)}

	for _, t := range T {
		ts := TriangleSmooth{}

		var p [3]Tuple

		for i := 0; i < 3; i++ {
			ts.V[i] = t.F[i].V
			ts.VN[i] = t.F[i].VN

			p[i] = V[ts.V[i]] // Vertex
		}

		ts.P1 = p[0]
		ts.E1 = p[1].Sub(p[0])
		ts.E2 = p[2].Sub(p[0])
		ts.N = ts.E2.CrossProduct(ts.E1).Normalize()

		mesh.F = append(mesh.F, ts)

		mesh.bbox = mesh.bbox.Union(Box{
			Point(Min3(p[0].X, p[1].X, p[2].X), Min3(p[0].Y, p[1].Y, p[2].Y), Min3(p[0].Z, p[1].Z, p[2].Z)),
			Point(Max3(p[0].X, p[1].X, p[2].X), Max3(p[0].Y, p[1].Y, p[2].Y), Max3(p[0].Z, p[1].Z, p[2].Z)),
		})
	}
}

func (s *Trimesh) Material() *Material {
	return s.material
}

func (s *Trimesh) SetMaterial(m *Material) {
	s.material = m
}

func (s *Trimesh) SetCulling(f bool) {
	s.culling = f
}

func (s *Trimesh) AddIntersections(ray Ray, xs *Intersections) {
	ray = ray.Transform(s.Tinverse)

	// Iterate over triangles and check for an intersection with each
	for i, _ := range s.F {
		// Backface culling
		if s.culling {
			dd := ray.Direction.DotProduct(s.F[i].N)
			if dd <= 0 {
				continue
			}
		}

		e1 := s.F[i].E1 // For performance, don't use range to get the triangle as it will duffcopy a lot of useless data
		e2 := s.F[i].E2

		dirCrossE2 := ray.Direction.CrossProduct(e2)
		det := e1.DotProduct(dirCrossE2)

		// Check if ray is parallel to triangle plane
		if det > -Epsilon && det < +Epsilon {
			continue
		}

		p1 := s.F[i].P1
		f := 1.0 / det
		p1ToOrigin := ray.Origin.Sub(p1)
		u := f * p1ToOrigin.DotProduct(dirCrossE2)

		// Ray misses by the p1-p3 edge
		if u < 0 || u > 1 {
			continue
		}

		originCrossE1 := p1ToOrigin.CrossProduct(e1)
		v := f * ray.Direction.DotProduct(originCrossE1)

		// Ray misses by the p1-p2 or p2-p3 edge
		if v < 0 || (u+v) > 1 {
			continue
		}

		// Ray hits
		h := f * e2.DotProduct(originCrossE1)

		id := xs.AddWithData(s, h)

		id.tIdx = i
		id.tU = u
		id.tV = v
	}
}

func (s *Trimesh) NormalAtEx(point Tuple, xs *Intersections, i Intersection) Tuple {
	id := xs.Data(i)

	idx := id.tIdx

	if len(s.VN) == 0 {
		return s.NormalToWorld(s.F[idx].N) // Flat normal
	}
	// ...else interpolate

	N1 := s.VN[s.F[idx].VN[0]]
	N2 := s.VN[s.F[idx].VN[1]]
	N3 := s.VN[s.F[idx].VN[2]]

	n := N2.Mul(id.tU).Add(N3.Mul(id.tV)).Add(N1.Mul(1 - id.tU - id.tV))

	return s.NormalToWorld(n)
}

func (s *Trimesh) Bounds() Box {
	return s.bbox
}

// Normalize fits the entire shape into a (-1,-1,-1) to (+1,+1,+1) box
func (s *Trimesh) Normalize() {
	bbox := s.bbox

	sx := bbox.Max.X - bbox.Min.X
	sy := bbox.Max.Y - bbox.Min.Y
	sz := bbox.Max.Z - bbox.Min.Z

	scale := Max3(sx, sy, sz) / 2

	for i, v := range s.V {
		cx := bbox.Min.X + sx/2
		cy := bbox.Min.Y + sy/2
		cz := bbox.Min.Z + sz/2

		x := v.X - cx
		y := v.Y - cy
		z := v.Z - cz

		x /= scale
		y /= scale
		z /= scale

		s.V[i] = Point(x, y, z)
	}

	for i, _ := range s.F {
		p1 := s.V[s.F[i].V[0]]
		s.F[i].P1 = p1 // Need to update P1 with new vertex value
		s.F[i].E1 = s.V[s.F[i].V[1]].Sub(p1)
		s.F[i].E2 = s.V[s.F[i].V[2]].Sub(p1)
	}

	s.bbox = NewBox(Point(-1, -1, -1), Point(+1, +1, +1))
}

// Autosmooth recomputes all normals at vertices
func (s *Trimesh) Autosmooth() {
	s.VN = make([]Tuple, len(s.V))

	for i, _ := range s.V {
		n := Vector(0, 0, 0)
		c := 0.0

		for j, f := range s.F {
			if f.V[0] == i || f.V[1] == i || f.V[2] == i {
				n = n.Add(f.N)
				c++
			}

			s.F[j].VN[0] = f.V[0]
			s.F[j].VN[1] = f.V[1]
			s.F[j].VN[2] = f.V[2]
		}

		if c > 0 {
			n = n.Mul(1 / c)
		}

		s.VN[i] = n
	}
}
