// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package shapes

import (
	"math"

	. "ascottix/funtracer/maths"
	. "ascottix/funtracer/textures"
)

type Cube struct {
}

func NewCube() *Shape {
	return NewShape("cube", &Cube{})
}

func (p *Cube) Bounds() Box {
	return Box{Point(-1, -1, -1), Point(+1, +1, +1)}
}

func (p *Cube) LocalIntersect(ray Ray) []float64 {

	checkAxis := func(origin, direction float64) (tmin, tmax float64) {
		tminNumerator := -1.0 - origin
		tmaxNumerator := +1.0 - origin

		if math.Abs(direction) >= Epsilon {
			tmin = tminNumerator / direction
			tmax = tmaxNumerator / direction
		} else {
			tmin = tminNumerator * math.Inf(+1)
			tmax = tmaxNumerator * math.Inf(+1)
		}

		if tmin > tmax {
			tmin, tmax = tmax, tmin
		}

		return tmin, tmax
	}

	xtmin, xtmax := checkAxis(ray.Origin.X, ray.Direction.X)
	ytmin, ytmax := checkAxis(ray.Origin.Y, ray.Direction.Y)
	ztmin, ztmax := checkAxis(ray.Origin.Z, ray.Direction.Z)

	tmin := math.Inf(-1)
	tmax := math.Inf(+1)

	if xtmin > tmin {
		tmin = xtmin
	}
	if ytmin > tmin {
		tmin = ytmin
	}
	if ztmin > tmin {
		tmin = ztmin
	}

	if xtmax < tmax {
		tmax = xtmax
	}
	if ytmax < tmax {
		tmax = ytmax
	}
	if ztmax < tmax {
		tmax = ztmax
	}

	if tmin <= tmax {
		return []float64{tmin, tmax}
	}

	return nil
}

func (p *Cube) LocalNormalAt(point Tuple) Tuple {
	ax := math.Abs(point.X)
	ay := math.Abs(point.Y)
	az := math.Abs(point.Z)

	maxc := math.Max(ax, math.Max(ay, az))

	switch maxc {
	case ax:
		return Vector(point.X, 0, 0)
	case ay:
		return Vector(0, point.Y, 0)
	default:
		return Vector(0, 0, point.Z)
	}
}

func (p *Cube) NormalAtHit(point Tuple, ii *IntersectionInfo) Tuple {
	n := p.LocalNormalAt(point)

	return n
}
