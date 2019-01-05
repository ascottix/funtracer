// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

// ObjCustomDice builds a dice centered in the origin, the size of an edge is 2*(1+R)
func ObjCustomDice(matBody, matDots *Material, R, D float64) *Csg {
	objects := []Groupable{}

	add := func(object *Shape) {
		objects = append(objects, object)
	}

	dice_edge := func(x, y, z, rotx, rotz float64) {
		cylinder := NewCylinder(-1, +1, true) // Cylinder must be capped or visual artifacts occur!
		cylinder.SetTransform(Translation(x, y, z), RotationX(rotx), RotationZ(rotz), Scaling(R, 1, R))
		add(cylinder)
	}

	dice_corner := func(x, y, z float64) {
		sphere := NewSphere()
		sphere.SetTransform(Translation(x, y, z), Scaling(R))
		add(sphere)
	}

	dice_face := func(x, y, z float64) {
		cube := NewCube()
		cube.SetTransform(Scaling(1+x, 1+y, 1+z))
		add(cube)
	}

	dice_number := func(x, y, z float64) {
		sphere := NewSphere()
		sphere.SetTransform(Translation(x, y, z), Scaling(R))
		add(sphere)

		// Debugf("translate(%.2f, %.2f, %.2f, scale(%.2f, sphere {} ))\n", x, y, z, R)
	}

	dice_edge(-1, 0, -1, 0, 0)
	dice_edge(-1, 0, +1, 0, 0)
	dice_edge(+1, 0, -1, 0, 0)
	dice_edge(+1, 0, +1, 0, 0)
	dice_edge(-1, -1, 0, Pi/2, 0)
	dice_edge(-1, +1, 0, Pi/2, 0)
	dice_edge(+1, -1, 0, Pi/2, 0)
	dice_edge(+1, +1, 0, Pi/2, 0)
	dice_edge(0, -1, -1, 0, Pi/2)
	dice_edge(0, +1, -1, 0, Pi/2)
	dice_edge(0, -1, +1, 0, Pi/2)
	dice_edge(0, +1, +1, 0, Pi/2)

	dice_corner(-1, -1, -1)
	dice_corner(-1, +1, -1)
	dice_corner(+1, -1, -1)
	dice_corner(+1, +1, -1)
	dice_corner(-1, -1, +1)
	dice_corner(-1, +1, +1)
	dice_corner(+1, -1, +1)
	dice_corner(+1, +1, +1)

	dice_face(R, 0, 0)
	dice_face(0, R, 0)
	dice_face(0, 0, R)

	csgBody := NewCsgUnion(objects...)
	csgBody.SetMaterial(matBody)

	objects = objects[:0]

	dice_number(0, 0, -1-R)  // 1
	dice_number(1+R, -D, -D) // 2
	dice_number(1+R, D, D)
	dice_number(0, -1-R, 0) // 3
	dice_number(-D, -1-R, -D)
	dice_number(D, -1-R, D)
	dice_number(-D, 1+R, -D) // 4
	dice_number(-D, 1+R, +D)
	dice_number(+D, 1+R, -D)
	dice_number(+D, 1+R, +D)
	dice_number(-1-R, D, D) // 5
	dice_number(-1-R, -D, D)
	dice_number(-1-R, -D, -D)
	dice_number(-1-R, D, -D)
	dice_number(-1-R, 0, 0)
	dice_number(-D, -D, 1+R) // 6
	dice_number(-D, 0, 1+R)
	dice_number(-D, D, 1+R)
	dice_number(+D, -D, 1+R)
	dice_number(+D, 0, 1+R)
	dice_number(+D, D, 1+R)

	csgDots := NewCsgUnion(objects...)
	csgDots.SetMaterial(matDots)

	csg := NewCsg(CsgDifference, csgBody, csgDots)

	return csg
}

func ObjDice(body, dots Color) *Csg {
	return ObjCustomDice(MatMatte(body), MatMatte(dots), 0.18, 0.60)
}

func init() {
	ObjDice(Black, White)
}
