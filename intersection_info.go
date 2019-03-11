// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

type IntersectionInfo struct {
	Intersection
	U, V           float64        // Surface coordinates of intersection point
	Point          Tuple          // Intersection point
	OverPoint      Tuple          // Intersection point adjusted a bit in the normal direction (over the surface), used for shadows
	UnderPoint     Tuple          // Intersection point adjusted a bit in the opposide normal direction (under the surface), used for refractions
	Eyev           Tuple          // Eye vector
	Normalv        Tuple          // Normal vector at intersection point
	SurfNormalv    Tuple          // Surface normal used by light (it's the normal vector possibly perturbed e.g. by a normal map)
	Reflectv       Tuple          // Reflected ray
	N1             float64        // Refractive index of "outside" (which ray is leaving) material
	N2             float64        // Refractive index of "inside" (which ray is entering) material
	Mat            MaterialParams // Material info
	Inside         bool           // True if the ray originates inside the intersected object
	HasSurfNormalv bool           // True if the surface normal may be different from the geometric normal
	// The following is for performance optimization only and does not contain actual information
	_containers []Hittable // To avoid allocating a new slice at every hit
}

func NewIntersectionInfo(i Intersection, r Ray, xs *Intersections) *IntersectionInfo {
	ii := IntersectionInfo{}

	ii.Update(i, r, xs)

	return &ii
}

func (ii *IntersectionInfo) Update(i Intersection, r Ray, xs *Intersections) {
	// Fill in the basic information
	ii.Intersection = i
	ii.Point = r.Position(i.T)
	ii.Eyev = r.Direction.Neg()
	ii.HasSurfNormalv = false

	n := i.O.NormalAtHit(ii, xs) // Get the normal at the intersection, necessary for all code that follows

	ii.Inside = n.DotProduct(ii.Eyev) < 0 // If the normal points away from the eye direction, it means the eye is inside the object

	if ii.Inside {
		n = n.Neg()
		ii.SurfNormalv = ii.SurfNormalv.Neg()
	}

	ii.Normalv = n
	ii.OverPoint = ii.Point.Add(n.Mul(Epsilon)) // Used to avoid objects casting shadows on themselves
	ii.UnderPoint = ii.Point.Sub(n.Mul(Epsilon))
	ii.Reflectv = r.Direction.Reflect(n)
	ii.N1 = 1
	ii.N2 = 1

	if !ii.HasSurfNormalv {
		ii.SurfNormalv = n
	}

	ii.O.Material().GetParamsAt(ii) // Get the material parameters

	// Handle refraction
	if xs != nil && ii.Mat.RefractLevel > 0 {
		// The purpose of this code is to get the refractive index of the material the ray is leaving
		// and of the material the ray is entering, it does so by tracking the ray thru all intersections
		// as it enters and leaves objects
		xs.Sort()

		hit := i

		containers := ii._containers[:0] // Reusing the same slice keeps memory clean and brings a very good performance boost

		removeContainer := func(o Hittable) bool {
			for i, v := range containers {
				if v == o {
					containers = append(containers[:i], containers[i+1:]...) // Remove object
					return true
				}
			}
			return false
		}

		for j := 0; j < xs.Len(); j++ {
			x := xs.At(j)

			if x.T == hit.T {
				// Get refraction index of entering ray
				if len(containers) == 0 {
					ii.N1 = 1.0 // Entering from "vacuum"
				} else {
					ii.N1 = containers[len(containers)-1].Material().Ior // Refraction index of last container
				}
			}

			// Add current object if not in the list, otherwise remove it
			if !removeContainer(x.O) {
				containers = append(containers, x.O)
			}

			if x.T == hit.T {
				// Get refraction index of leaving ray
				if len(containers) == 0 {
					ii.N2 = 1.0 // Leaving into "vacuum"
				} else {
					ii.N2 = containers[len(containers)-1].Material().Ior // Refraction index of last container
				}

				break
			}
		}
	}
}

func (ii *IntersectionInfo) GetNormalMap() (nmap *ImageTexture) { // TODO: wrong type
	if ii.O != nil {
		nmap = ii.O.Material().NormalMap
	}

	return nmap
}
