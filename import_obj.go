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

type ObjInfoFace struct {
	V	[3]int	// Indices in vertex array
	VN	[3]int	// Indices in vertex normals array
	G 	int		// Group
}

type ObjInfoGroup struct {
	Name string
}

type ObjInfo struct {
	Vertices  []Tuple
	VNormals  []Tuple
	F 		  []ObjInfoFace
	Groups    []ObjInfoGroup
}

func (o *ObjInfo) Dump() {
	Debugf("v=%d, vn=%d, f=%d\n", len(o.Vertices), len(o.VNormals), len(o.F))
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

	info.Groups = []ObjInfoGroup{ObjInfoGroup{Name:"default"}}

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
			v := ReBlank.Split(strings.TrimSpace(line), -1)
			i1, _, n1 := f2cs(v[1])
			for i := 3; i < len(v); i++ {
				i2, _, n2 := f2cs(v[i-1])
				i3, _, n3 := f2cs(v[i-0])

				f := ObjInfoFace{
					[3]int{i1-1, i2-1, i3-1},
					[3]int{n1-1, n2-1, n3-1},
					len(info.Groups)-1,
				}

				info.F = append(info.F, f)
			}
		}

		// Group
		s = ReGroup.FindStringSubmatch(line)
		if s != nil {
			g := ObjInfoGroup{
				Name: strings.TrimSpace(s[1]),
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
