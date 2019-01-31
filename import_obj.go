// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// See:
// https://en.wikipedia.org/wiki/Wavefront_.obj_file
// https://www.mathworks.com/matlabcentral/mlc-downloads/downloads/submissions/27982/versions/5/previews/help%20file%20format/MTL_format.html
type ObjMaterial struct {
	Name string
	Ka   Color   // Ambient
	Kd   Color   // Diffuse
	Ke   Color   // Emissive
	Ks   Color   // Specular
	Ns   float64 // Specular exponent (typically from 0 to 1000)
	Tf   Color   // Transmission filter
	Tr   float64 // Transparency (can be also specified with 'd' i.e. "dissolve", then Tr=1-d)
	Ni   float64 // Index of refraction
}

type ObjInfoFace struct {
	V  [3]int // Indices in vertex array
	VN [3]int // Indices in vertex normals array
	G  int    // Group
}

type ObjInfoGroup struct {
	Name string
}

type ObjInfo struct {
	V      []Tuple
	VN     []Tuple
	F      []ObjInfoFace
	Groups []ObjInfoGroup
}

func (o *ObjInfo) Bounds() Box {
	bbox := NewBox(PointAtInfinity(+1), PointAtInfinity(-1))

	for _, f := range o.F {
		p1 := o.V[f.V[0]]
		p2 := o.V[f.V[1]]
		p3 := o.V[f.V[2]]

		bbox = bbox.Union(Box{
			Point(Min3(p1.X, p2.X, p3.X), Min3(p1.Y, p2.Y, p3.Y), Min3(p1.Z, p2.Z, p3.Z)),
			Point(Max3(p1.X, p2.X, p3.X), Max3(p1.Y, p2.Y, p3.Y), Max3(p1.Z, p2.Z, p3.Z)),
		})
	}

	return bbox
}

// Normalize fits the entire mesh into a (-1,-1,-1) to (+1,+1,+1) box
func (o *ObjInfo) Normalize() {
	bbox := o.Bounds()

	sx := bbox.Max.X - bbox.Min.X
	sy := bbox.Max.Y - bbox.Min.Y
	sz := bbox.Max.Z - bbox.Min.Z

	scale := Max3(sx, sy, sz) / 2

	for i, v := range o.V {
		cx := bbox.Min.X + sx/2
		cy := bbox.Min.Y + sy/2
		cz := bbox.Min.Z + sz/2

		x := v.X - cx
		y := v.Y - cy
		z := v.Z - cz

		x /= scale
		y /= scale
		z /= scale

		o.V[i] = Point(x, y, z)
	}
}

// Autosmooth sets all vertex normals to the average of the normals of the adjacent triangles
func (o *ObjInfo) Autosmooth() {
	o.VN = make([]Tuple, len(o.V))
	c := make([]int, len(o.V))

	for j, _ := range o.F {
		v0 := o.F[j].V[0]
		v1 := o.F[j].V[1]
		v2 := o.F[j].V[2]

		// Compute face normal
		p0 := o.V[v0]
		p1 := o.V[v1]
		p2 := o.V[v2]
		e1 := p1.Sub(p0)
		e2 := p2.Sub(p0)
		fn := e2.CrossProduct(e1).Normalize()

		// Add to vertex normals
		o.VN[v0] = o.VN[v0].Add(fn)
		o.VN[v1] = o.VN[v1].Add(fn)
		o.VN[v2] = o.VN[v2].Add(fn)

		c[v0]++
		c[v1]++
		c[v2]++

		// Update face information
		o.F[j].VN[0] = v0
		o.F[j].VN[1] = v1
		o.F[j].VN[2] = v2
	}

	for i, _ := range o.V {
		n := c[i]
		if n > 0 {
			o.VN[i] = o.VN[i].Mul(1 / float64(n))
		}
	}
}

func (o *ObjInfo) Dump() {
	Debugf("v=%d, vn=%d, f=%d\n", len(o.V), len(o.VN), len(o.F))
}

var (
	ReBlank = regexp.MustCompile(`\s+`)
)

func ParseWavefrontObj(rd io.Reader) *ObjInfo {
	var info ObjInfo

	info.Groups = []ObjInfoGroup{ObjInfoGroup{Name: "default"}}

	s2f := func(s string) float64 {
		f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)

		if err != nil {
			panic(err)
		}

		return f
	}

	s2i := func(s string) int {
		s = strings.TrimSpace(s)

		if s == "" {
			return 0
		}

		i, err := strconv.Atoi(s)

		if err != nil {
			panic(err)
		}

		return i
	}

	f2cs := func(s string) (v, tn, vn int) {
		cs := strings.Split(s, "/")
		v = s2i(cs[0])
		if len(cs) > 1 {
			tn = s2i(cs[1])
		}
		if len(cs) > 2 {
			vn = s2i(cs[2])
		}
		return
	}

	reader := bufio.NewReader(rd)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		var s []string

		s = ReBlank.Split(strings.TrimSpace(line), -1)

		switch line[0] {
		case 'v':
			// Vertex
			if line[1] == 'n' {
				// Vertex normal
				info.VN = append(info.VN, Point(s2f(s[1]), s2f(s[2]), s2f(s[3])))
			} else if line[1] == 't' {
				// Texture
			} else {
				// Standard vertex
				info.V = append(info.V, Point(s2f(s[1]), s2f(s[2]), s2f(s[3])))
			}
		case 'f':
			// Polygon
			i1, _, n1 := f2cs(s[1])
			for i := 3; i < len(s); i++ {
				i2, _, n2 := f2cs(s[i-1])
				i3, _, n3 := f2cs(s[i-0])

				f := ObjInfoFace{
					[3]int{i1 - 1, i2 - 1, i3 - 1},
					[3]int{n1 - 1, n2 - 1, n3 - 1},
					len(info.Groups) - 1,
				}

				info.F = append(info.F, f)
			}
		case 'g':
			// Group
			g := ObjInfoGroup{}

			if len(s) > 1 {
				g.Name = s[1]
			}

			info.Groups = append(info.Groups, g)
		}
	}

	return &info
}

func ParseWavefrontObjFromString(src string) *ObjInfo {
	return ParseWavefrontObj(strings.NewReader(src))
}

func ParseWavefrontObjFromFile(filename string) *ObjInfo {
	f, err := os.Open(filename)

	if err == nil {
		defer f.Close()

		return ParseWavefrontObj(f)
	}

	panic(err)
}
