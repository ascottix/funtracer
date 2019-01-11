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

type ObjInfo struct {
	Vertices  []Tuple
	Triangles []*Triangle
	VNormals  []Tuple
	Faces     []TrimeshTriangleInfo
	Groups    []*Group
	Bounds    Box
}

func (o *ObjInfo) Dump() {
	Debugf("v=%d, vn=%d, f=%d, bbox=%+v\n", len(o.Vertices), len(o.VNormals), len(o.Faces), o.Bounds)
}

var (
	ReVertex       = regexp.MustCompile(`\s*v\s+(.+)\s+(.+)\s+([^\s]+)\s*`)
	ReVertexNormal = regexp.MustCompile(`\s*vn\s+(.+)\s+(.+)\s+([^\s]+)\s*`)
	RePolygon      = regexp.MustCompile(`\s*f(\s+\d+(/\d*/\d+)?){3,}\s*`)
	ReGroup        = regexp.MustCompile(`g\s+(\w.*)`)
	ReBlank        = regexp.MustCompile(`\s+`)
)

func ParseWavefrontObj(rd io.Reader) *ObjInfo {
	var info ObjInfo

	info.Bounds = Box{PointAtInfinity(+1), PointAtInfinity(-1)}
	info.Groups = []*Group{NewGroup()}

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

		// Vertex
		s = ReVertex.FindStringSubmatch(line)
		if s != nil {
			info.Vertices = append(info.Vertices, Point(s2f(s[1]), s2f(s[2]), s2f(s[3])))
		}

		// Vertex normal
		s = ReVertexNormal.FindStringSubmatch(line)
		if s != nil {
			info.VNormals = append(info.VNormals, Point(s2f(s[1]), s2f(s[2]), s2f(s[3])))
		}

		// Polygon
		s = RePolygon.FindStringSubmatch(line)
		if s != nil {
			vs := info.Vertices

			v := ReBlank.Split(strings.TrimSpace(line), -1)
			i1, _, n1 := f2cs(v[1])
			for i := 3; i < len(v); i++ {
				i2, _, n2 := f2cs(v[i-1])
				i3, _, n3 := f2cs(v[i-0])

				t := NewTriangle(vs[i1-1], vs[i2-1], vs[i3-1])

				info.Triangles = append(info.Triangles, t)

				info.Groups[len(info.Groups)-1].Add(NewShape("triangle", t))

				face := TrimeshTriangleInfo{}
				face.F[0].V = i1 - 1
				face.F[1].V = i2 - 1
				face.F[2].V = i3 - 1
				face.F[0].VN = n1 - 1
				face.F[1].VN = n2 - 1
				face.F[2].VN = n3 - 1

				info.Faces = append(info.Faces, face)

				info.Bounds = info.Bounds.Union(t.Bounds())
			}
		}

		// Group
		s = ReGroup.FindStringSubmatch(line)
		if s != nil {
			g := NewGroup()
			g.SetName(strings.TrimSpace(s[1]))

			info.Groups[0].Add(g)

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
