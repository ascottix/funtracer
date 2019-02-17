// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

type Plane struct {
}

// NewPlane returns a Shape based on a plane on the x and z axis (i.e. y=0 always)
func NewPlane() *Shape {
	return NewShape("plane", &Plane{})
}

func (p *Plane) Bounds() Box {
	return Box{
		PointAtInfinity(-1),
		PointAtInfinity(+1),
	}
}

func (p *Plane) LocalIntersect(ray Ray) []float64 {
	if ray.Direction.Y <= -Epsilon || ray.Direction.Y >= Epsilon {
		t := -ray.Origin.Y / ray.Direction.Y

		return []float64{t}
	}
	// ...else the ray is parallel to the plane

	return nil
}

func (p *Plane) LocalNormalAt(point Tuple) Tuple {
	return Vector(0, 1, 0) // Plane is always on the x and z axis, so normal does not depend on point
}

func (p *Plane) NormalAtHit(point Tuple, ii *IntersectionInfo) Tuple {
	ii.U = point.X
	ii.V = point.Z

	if nmap := ii.GetNormalMap(); nmap != nil {
		n := nmap.NormalAtHit(ii)

		ii.HasSurfNormalv = true
		ii.SurfNormalv = Vector(n.X, n.Z, n.Y).Normalize()
	}

	return p.LocalNormalAt(point)
}
