// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/scanner"
)

type SbtParserOptions struct {
	FilenameBase string
}

// ParseSbtScene parses a scene description based on the .ray format, see:
// https://www.cs.utexas.edu/users/fussell/courses/cs384g/projects/raytracing/fileformat.html
// http://www.cs.cmu.edu/afs/cs.cmu.edu/academic/class/15864-s04/www/assignment4/format.html
//
// The orignal file format has been modified to match the available features,
// compatibility has never been thoroughly tested and is definitely not guaranteed!
func ParseSbtScene(reader io.Reader, options *SbtParserOptions) (scene *Scene, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			Debugln(err)
		}
	}()

	// Init scanner
	var s scanner.Scanner
	s.Init(reader)

	var token rune
	var tokenText string

	debug := false

	// Init object library
	objects := make(map[string]Groupable)

	// Init material library
	materials := make(map[string]*Material)
	defaultMaterial := NewMaterial()

	// Init scene
	scene = NewScene()
	scene.World.Ambient = Black

	// Get the next token from the scanner
	next := func() {
		token = s.Scan()
		tokenText = s.TokenText()

		if token == scanner.Int {
			token = scanner.Float
		}

		if debug {
			fmt.Printf("%+v %s: %s\n", token, s.Position, tokenText)
		}
	}

	// Raise an unrecoverable error
	raise := func() {
		msg := fmt.Sprintf("unexpected token %v (%s), pos=%s", token, tokenText, s.Position)

		panic(errors.New(msg))
	}

	// Check if current token matches expectation, gets the next token if successful
	check := func(text string) bool {
		f := tokenText == text
		if f {
			next()
		}
		return f
	}

	// Checks if current token matches expectation, raises an error if it doesn't
	match := func(id rune) {
		if token != id {
			if debug {
				Debugln("expected token, have", token, "want", id)
			}
			raise()
		}
		next()
	}

	// Matches a single float value
	var matchFloat func() float64

	matchFloat = func() (f float64) {
		check("+")
		negate := check("-")

		if check("pi") {
			f = Pi
			if check("/") {
				d := matchFloat()
				f /= d
			}
		} else {
			f, _ = strconv.ParseFloat(tokenText, 64)
			match(scanner.Float)
		}

		check(",")

		if negate {
			f = -f
		}

		return
	}

	parseTuple := func() (float64, float64, float64) {
		match('=')
		match('(')
		x := matchFloat()
		y := matchFloat()
		z := matchFloat()
		match(')')
		check(";")

		return x, y, z
	}

	parseFloat := func() float64 {
		match('=')
		f := matchFloat()
		check(";")

		return f
	}

	parseColor := func() (c Color) {
		if token == scanner.String {
			v, _ := strconv.Unquote(s.TokenText())
			match(scanner.String)
			c = CSS(v)
		} else {
			match('(')
			r := matchFloat()
			g := matchFloat()
			b := matchFloat()
			match(')')
			c = RGB(r, g, b)
		}
		check(",")

		return
	}

	parseString := func() string {
		check("=")
		v, _ := strconv.Unquote(s.TokenText())
		match(scanner.String)
		match(';')

		return v
	}

	parsePragma := func() {
		pragma := parseString()

		switch pragma {
		// For compatibility with scenes designed before gamma correction was added
		case "gamma=1.0":
			scene.World.ErpCanvasToImage = ErpLinear
		}
	}

	parseCamera := func() {
		pos := Point(0, 0, -4)
		dir := Vector(0, 0, 1)
		upd := Vector(0, 1, 0)
		fov := Pi / 2
		w := 0
		h := 0

		match('{')
		for !check("}") {
			switch {
			case check("position"):
				pos = Point(parseTuple())
			case check("viewdir"):
				dir = Vector(parseTuple())
			case check("target"): // Specify direction using a target point, must come after "position"
				dir = Vector(parseTuple()).Sub(pos)
			case check("aspectratio"):
				parseFloat() // Ignored
			case check("updir"):
				upd = Vector(parseTuple())
			case check("fov"):
				fov = (Pi * 2 * parseFloat() / 360)
			case check("fovrad"): // Specify FOV in radians
				fov = parseFloat()
			case check("viewsize"):
				w = int(parseFloat())
				h = int(matchFloat())
				match(';')

			default:
				raise()
			}
		}

		scene.Camera = NewCamera(w, h, fov)
		scene.Camera.SetTransform(EyeViewpoint(pos, pos.Add(dir), upd))
	}

	parseAmbientLight := func() {
		var col Color

		match('{')
		for !check("}") {
			switch {
			case check("colour"), check("color"):
				col = RGB(parseTuple())
			default:
				raise()
			}
		}

		scene.World.Ambient = scene.World.Ambient.Add(col)
	}

	parsePointLight := func() {
		var pos Tuple
		var col Color

		match('{')
		for !check("}") {
			switch {
			case check("position"):
				pos = Point(parseTuple())
			case check("colour"), check("color"):
				col = RGB(parseTuple())
			case check("constant_attenuation_coeff"), check("linear_attenuation_coeff"), check("quadratic_attenuation_coeff"):
				parseFloat() // Ignored
			default:
				raise()
			}
		}

		scene.World.AddLights(NewPointLight(pos, col))
	}

	parseDirectionalLight := func() {
		var dir Tuple
		var col Color

		match('{')
		for !check("}") {
			switch {
			case check("direction"):
				dir = Vector(parseTuple())
			case check("colour"), check("color"):
				col = RGB(parseTuple())
			default:
				raise()
			}
		}

		scene.World.AddLights(NewDirectionalLight(dir, col))
	}

	checkTransform := func(t Matrix) (Matrix, int) {
		n := 0

		trans := func(t, m Matrix) Matrix {
			return t.Mul(m)
		}

		for {
			switch {
			case check("rotate"):
				match('(')
				x := matchFloat() // Vector(x, y, z) is the axis of rotation
				y := matchFloat()
				z := matchFloat()
				a := matchFloat() // Rotation angle

				if x != 0 {
					t = trans(t, RotationX(a))
				}

				if y != 0 {
					t = trans(t, RotationY(a))
				}

				if z != 0 {
					t = trans(t, RotationZ(a))
				}
				n++
			case check("rotate_x"):
				match('(')
				t = trans(t, RotationX(matchFloat()))
				n++
			case check("rotate_y"):
				match('(')
				t = trans(t, RotationY(matchFloat()))
				n++
			case check("rotate_z"):
				match('(')
				t = trans(t, RotationZ(matchFloat()))
				n++
			case check("translate"):
				match('(')
				t = trans(t, Translation(matchFloat(), matchFloat(), matchFloat()))
				n++
			case check("scale"):
				match('(')
				x := matchFloat()
				y := x
				z := x
				if token == scanner.Float {
					y = matchFloat()
					z = matchFloat()
				}
				t = trans(t, Scaling(x, y, z))
				n++
			case check("transform"):
				match('(')
				m := Identity()
				for row := 0; row < 4; row++ {
					match('(')
					for col := 0; col < 4; col++ {
						m.SetAt(row, col, matchFloat())
					}
					match(')')
					match(',')
				}
				t = trans(t, m)
				n++
			default:
				return t, n
			}
		}
	}

	parsePattern := func() (p Pattern) {
		match('{')
		t, n := checkTransform(Identity())
		switch {
		case check("stripe"):
			p = NewStripePattern(parseColor(), parseColor())
		case check("checker"):
			p = NewCheckerPattern(parseColor(), parseColor())
		case check("gradient"):
			p = NewCheckerPattern(parseColor(), parseColor())
		default:
			p = NewSolidColorPattern(parseColor())
		}
		if n > 0 {
			p.SetTransform(t)
			for ; n > 0; n-- {
				match(')')
			}
		}
		match('}')

		return
	}

	parseMaterial := func() (m *Material, name string) {
		if check("{") {
			var p Pattern

			m = NewMaterial()

			m.SetAmbient(defaultMaterial.Ambient)
			m.SetShininess(defaultMaterial.Shininess)
			m.SetSpecular(defaultMaterial.Specular)
			m.SetReflective(defaultMaterial.ReflectLevel)
			m.SetRefractive(defaultMaterial.RefractLevel, defaultMaterial.Ior)

			for !check("}") {
				switch {
				case check("name"):
					name = parseString()
				case check("pattern"):
					p = parsePattern()
				case check("diffuse"):
					// Try to support both the original format and our modifications
					if p == nil {
						match('=')
						p = NewSolidColorPattern(parseColor())
						d := 1.0
						if check("*") {
							d = matchFloat()
							match(';')
						} else {
							check(";")
						}
						m.SetDiffuse(d)
					} else {
						d, _, _ := parseTuple() // Only one component supported
						m.SetDiffuse(d)
					}
				case check("specular"):
					s, _, _ := parseTuple() // Only one component supported
					m.SetSpecular(s)
				case check("shininess"):
					m.SetShininess(parseFloat())
				case check("ambient"):
					a, _, _ := parseTuple() // Only one component supported
					m.SetAmbient(a)
				case check("emissive"):
					parseTuple() // Unsupported
				case check("reflective"):
					m.SetReflect(1, RGB(parseTuple()))
				case check("transmissive"):
					m.SetRefract(1, RGB(parseTuple()))
				case check("index"):
					m.SetIor(parseFloat())
				default:
					raise()
				}
			}

			if p == nil {
				p = defaultMaterial.Pattern
			}
			m.SetPattern(p)

			check(";")
		} else {
			// Named material
			name = parseString()
			m = materials[name]
			if m == nil {
				Debugf("*** Warning: cannot find material '%s'\n", name)
				m = NewMaterial()
			}
		}

		return
	}

	groupStack := []*Group{}

	add := func(object Groupable) {
		objects[object.Name()] = object

		if len(groupStack) > 0 {
			groupStack[len(groupStack)-1].Add(object)
		} else {
			scene.World.AddObjects(object)
		}
	}

	checkStandardAttributes := func(object Groupable) bool {
		f := true

		switch {
		case check("name"):
			object.SetName(parseString())
		case check("material"):
			match('=')
			material, _ := parseMaterial()
			object.SetMaterial(material)
		default:
			f = false
		}

		return f
	}

	parsePolymesh := func(transform Matrix) {
		g := NewGroup()
		g.SetTransform(transform)

		autosmooth := false
		var info *ObjInfo

		match('{')
		for !check("}") {
			switch {
			case checkStandardAttributes(g):
				// Nothing to do
			case check("objfile"):
				filename := parseString()
				if _, err := os.Stat(filename); os.IsNotExist(err) && options != nil {
					filename = filepath.Join(options.FilenameBase, filename)
				}

				info = ParseWavefrontObjFromFile(filename)
				Debugf("%d triangles loaded from %q\n", len(info.F), filename)

				info.Normalize()

				if autosmooth {
					info.Autosmooth()
				}

				mesh := NewTrimesh(info, -1)
				mesh.AddToGroup(g)
			case check("gennormals"):
				match('=')
				if check("true") {
					autosmooth = true
				} else {
					check("false")
				}
				match(';')
			default:
				raise()
			}
		}

		g.BuildBVH() // For now, always build a BVH
		add(g)
	}

	var parseObject func()

	parseCsg := func(op CsgOp, transform Matrix) {
		g := NewGroup()
		groupStack = append(groupStack, g) // Push a temporary group into stack so new objects get added to it

		match('{')
		// Standard attributes must be specified before the objects, but they are ignored for now
		for checkStandardAttributes(g) {
		}
		parseObject()
		check(",")
		parseObject()
		match('}')

		groupStack = groupStack[:len(groupStack)-1] // Pop group, but don't add it

		csg := NewCsg(op, g.members[0], g.members[1])
		csg.SetTransform(transform)

		add(csg)
	}

	parseGroup := func(transform Matrix) {
		g := NewGroup()
		g.SetTransform(transform)

		groupStack = append(groupStack, g) // Push group into stack so new objects get added to it

		for !check("}") {
			switch {
			case checkStandardAttributes(g):
				// Nothing to do
			default:
				parseObject()
			}
		}

		groupStack = groupStack[:len(groupStack)-1] // Pop group

		add(g)
	}

	parseClone := func(transform Matrix) {
		name, _ := strconv.Unquote(s.TokenText())
		match(scanner.String)

		object := objects[name]
		if object == nil {
			Debugf("cannot find '%s' for cloning\n", name)
			raise()
		}

		s := object.Clone()

		match('{')
		for !check("}") {
			switch {
			case checkStandardAttributes(s):
				// Nothing to do
			default:
				raise()
			}
		}

		s.SetTransform(transform)

		add(s)
	}

	parseObject = func() {
		shape := func(object *Shape, transform Matrix) {
			match('{')
			for !check("}") {
				switch {
				case checkStandardAttributes(object):
					// Nothing to do
				case check("color"): // Shortcut to set a solid color material
					match('=')
					object.Material().SetPattern(NewSolidColorPattern(parseColor()))
					check(";")
				default:
					raise()
				}
			}
			object.SetTransform(transform)

			add(object)
		}

		t, n := checkTransform(Identity())
		switch {
		case check("{"):
			parseGroup(t)
		case check("box"):
			shape(NewCube(), t.Mul(Scaling(0.5, 0.5, 0.5))) // A box is a cube that goes from -0.5 to +0.5, so we need to add an initial transformation
		case check("sphere"):
			shape(NewSphere(), t)
		case check("cylinder"):
			shape(NewCylinder(0, 1, false), t.Mul(RotationZ(Pi/2))) // This cylinder is aligned on the Z axis rather than the Y axis
		// Objects not included in the original format or not fully compatible
		case check("group"):
			match('{')
			parseGroup(t)
		case check("clone"):
			parseClone(t)
		case check("cube"):
			shape(NewCube(), t)
		case check("cyl"):
			shape(NewCylinder(-1, +1, true), t)
		case check("cyl_infinite"):
			shape(NewInfiniteCylinder(), t)
		case check("cyl_uncapped"):
			shape(NewCylinder(-1, +1, false), t)
		case check("cone"):
			miny := matchFloat()
			maxy := matchFloat()
			shape(NewCone(miny, maxy, true), t)
		case check("cone_infinite"):
			shape(NewInfiniteCone(), t)
		case check("cone_uncapped"):
			miny := matchFloat()
			maxy := matchFloat()
			shape(NewCone(miny, maxy, false), t)
		case check("plane"):
			shape(NewPlane(), t)
		case check("polymesh"):
			parsePolymesh(t)
		case check("intersect"):
			parseCsg(CsgIntersection, t)
		case check("diff"):
			parseCsg(CsgDifference, t)
		case check("union"):
			parseCsg(CsgUnion, t)
		default:
			raise()
		}

		for ; n > 0; n-- {
			match(')')
		}
	}

	next()
	match(scanner.Ident) // FUN or SBT
	match('-')
	match(scanner.Ident) // raytracer
	match(scanner.Float) // 1.0

	for token != scanner.EOF {
		switch {
		case check("pragma"):
			parsePragma()
		case check("camera"):
			parseCamera()
		case check("ambient_light"):
			parseAmbientLight()
		case check("point_light"):
			parsePointLight()
		case check("directional_light"):
			parseDirectionalLight()
		case check("material"):
			material, name := parseMaterial()
			materials[name] = material
		case check(";"):
			// Just skip
		default:
			parseObject()
		}
	}

	return
}

func ParseSbtSceneFromString(src string) (*Scene, error) {
	return ParseSbtScene(strings.NewReader(src), nil)
}

func ParseSbtSceneFromFile(filename string) (*Scene, error) {
	f, err := os.Open(filename)

	if err == nil {
		defer f.Close()

		options := SbtParserOptions{
			FilenameBase: filepath.Dir(filename),
		}

		return ParseSbtScene(f, &options)
	}

	return nil, err
}
